package v4_test

import (
	"context"
	"errors"
	"testing"
	"time"

	auditlog "github.com/larsartmann/samber-do-auditlog"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func TestAuditMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("records before and after events", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)
		mw := v4.AuditMiddleware[testCLIConfig](plugin)

		err := mw(t.Context(), &testCLIConfig{}, v4.CommandInfo{Name: "run"}, func() error {
			return nil
		})
		if err != nil {
			t.Fatalf("middleware returned error: %v", err)
		}

		events := plugin.Events()
		if len(events) != 2 {
			t.Fatalf("Events() length = %d, want 2", len(events))
		}

		before, after := events[0], events[1]

		if !before.IsCommand() || !after.IsCommand() {
			t.Fatalf("expected command events, got types %q/%q", before.EventType, after.EventType)
		}

		if before.Phase != auditlog.PhaseBefore {
			t.Errorf("first event phase = %q, want %q", before.Phase, auditlog.PhaseBefore)
		}

		if after.Phase != auditlog.PhaseAfter {
			t.Errorf("second event phase = %q, want %q", after.Phase, auditlog.PhaseAfter)
		}

		if before.ServiceName != "run" || after.ServiceName != "run" {
			t.Errorf("ServiceName = %q/%q, want %q on both", before.ServiceName, after.ServiceName, "run")
		}

		if before.DurationMs != nil {
			t.Errorf("before event DurationMs = %v, want nil", *before.DurationMs)
		}

		if after.DurationMs == nil || *after.DurationMs < 0 {
			t.Errorf("after event DurationMs = %v, want non-nil and >= 0", after.DurationMs)
		}

		if after.Error != nil {
			t.Errorf("after event Error = %v, want nil on success", *after.Error)
		}
	})

	t.Run("propagates error and records it", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)
		mw := v4.AuditMiddleware[testCLIConfig](plugin)

		sentinel := errors.New("boom")

		err := mw(t.Context(), &testCLIConfig{}, v4.CommandInfo{Name: "fail"}, func() error {
			return sentinel
		})
		if !errors.Is(err, sentinel) {
			t.Fatalf("middleware error = %v, want the sentinel from next", err)
		}

		events := plugin.Events()
		if len(events) != 2 {
			t.Fatalf("Events() length = %d, want 2", len(events))
		}

		after := events[1]

		if after.Error == nil {
			t.Fatal("after event Error = nil, want the recorded error message")
		}

		if *after.Error != sentinel.Error() {
			t.Errorf("after event Error = %q, want %q", *after.Error, sentinel.Error())
		}
	})

	t.Run("nil plugin is a passthrough", func(t *testing.T) {
		t.Parallel()

		mw := v4.AuditMiddleware[testCLIConfig](nil)

		called := false
		err := mw(t.Context(), &testCLIConfig{}, v4.CommandInfo{Name: "run"}, func() error {
			called = true

			return nil
		})
		if err != nil {
			t.Fatalf("middleware returned error: %v", err)
		}

		if !called {
			t.Error("next was not called for nil plugin")
		}
	})

	t.Run("prefers FullPath over Name", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)
		mw := v4.AuditMiddleware[testCLIConfig](plugin)

		info := v4.CommandInfo{Name: "cmd", FullPath: "root sub cmd"}

		_ = mw(t.Context(), &testCLIConfig{}, info, func() error { return nil })

		events := plugin.Events()
		if len(events) == 0 {
			t.Fatal("Events() is empty")
		}

		if got := string(events[0].ServiceName); got != "root sub cmd" {
			t.Errorf("ServiceName = %q, want FullPath %q", got, "root sub cmd")
		}
	})
}

func TestWithAuditMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("records command events through full CLI execution", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)

		cli, err := v4.NewCLI(
			"test", "Test", testCLIConfig{},
			v4.WithAuditLog(plugin),
			v4.WithAuditMiddleware[testCLIConfig](plugin),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		addCommand(t, cli, newTestCLICommand[testCLIConfig](t, "auditme"))

		if err := cli.ExecuteWithArgs(t.Context(), []string{"auditme"}); err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		var commandEvents []auditlog.Event

		for _, evt := range plugin.Events() {
			if evt.IsCommand() {
				commandEvents = append(commandEvents, evt)
			}
		}

		if len(commandEvents) != 2 {
			t.Fatalf("command events = %d, want 2 (before + after)", len(commandEvents))
		}

		before, after := commandEvents[0], commandEvents[1]

		if got := string(before.ServiceName); got != "test auditme" {
			t.Errorf("before ServiceName = %q, want %q", got, "test auditme")
		}

		if got := string(after.ServiceName); got != "test auditme" {
			t.Errorf("after ServiceName = %q, want %q", got, "test auditme")
		}

		if after.DurationMs == nil {
			t.Error("after event DurationMs = nil, want a recorded duration")
		}

		if after.Error != nil {
			t.Errorf("after event Error = %q, want nil on success", *after.Error)
		}
	})

	t.Run("records the handler error on failure", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)

		cli, err := v4.NewCLI(
			"test", "Test", testCLIConfig{},
			v4.WithAuditLog(plugin),
			v4.WithAuditMiddleware[testCLIConfig](plugin),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		failing := errors.New("handler failed")
		cmd, err := v4.NewCommand(
			"fails", v4.NoFlags{},
			func(_ context.Context, _ *testCLIConfig, _ v4.NoFlags) error { return failing },
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		addCommand(t, cli, cmd)

		_ = cli.ExecuteWithArgs(t.Context(), []string{"fails"})

		var after *auditlog.Event

		for i, evt := range plugin.Events() {
			if evt.IsCommand() && evt.Phase == auditlog.PhaseAfter {
				evtCopy := plugin.Events()[i]
				after = &evtCopy
			}
		}

		if after == nil {
			t.Fatal("no PhaseAfter command event recorded")
		}

		if after.Error == nil || *after.Error != failing.Error() {
			t.Errorf("after event Error = %v, want %q", after.Error, failing.Error())
		}
	})

	t.Run("middleware context reaches the handler", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)

		cli, err := v4.NewCLI(
			"test", "Test", testCLIConfig{},
			v4.WithAuditMiddleware[testCLIConfig](plugin),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type ctxKey struct{}

		var gotDeadline bool
		cmd, err := v4.NewCommand(
			"deadline", v4.NoFlags{},
			func(ctx context.Context, _ *testCLIConfig, _ v4.NoFlags) error {
				_, gotDeadline = ctx.Deadline()

				return nil
			},
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		addCommand(t, cli, cmd)

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		if err := cli.ExecuteWithArgs(ctx, []string{"deadline"}); err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if !gotDeadline {
			t.Error("handler did not observe the caller's context")
		}
	})
}
