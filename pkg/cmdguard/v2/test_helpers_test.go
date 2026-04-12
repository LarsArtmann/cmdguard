package v2

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func assertNoError(t *testing.T, err error, msg ...string) {
	t.Helper()

	if err != nil {
		if len(msg) > 0 {
			t.Fatalf("%s: %v", msg[0], err)
		}

		t.Fatalf("unexpected error: %v", err)
	}
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()

	didPanic := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()

		fn()
	}()

	if !didPanic {
		t.Error("expected panic, got none")
	}
}

func assertPanicsWithMessage(t *testing.T, fn func(), msg string) {
	t.Helper()

	didPanic := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()

		fn()
	}()

	if !didPanic {
		t.Errorf("expected panic with message containing %q, got none", msg)
	}
}

func assertDurationField(t *testing.T, d Duration, expected time.Duration) {
	t.Helper()

	if d.Duration() != expected {
		t.Errorf("expected %v, got %v", expected, d.Duration())
	}
}

func assertErrorContains(t *testing.T, err error, substrings ...string) {
	t.Helper()

	errMsg := err.Error()
	for _, s := range substrings {
		if !strings.Contains(errMsg, s) {
			t.Errorf("error should contain %q, got %q", s, errMsg)
		}
	}
}

func assertStderrContains(t *testing.T, stderr string, substrings ...string) {
	t.Helper()

	for _, s := range substrings {
		if !strings.Contains(strings.ToLower(stderr), strings.ToLower(s)) {
			t.Errorf("stderr should contain %q, got %q", s, stderr)
		}
	}
}

type testAppConfig struct {
	Verbose bool   `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
	Output  string `default:"-"     flag:"output"  help:"Output file"           short:"o"`
}

func newTestCmd(use string, err ...error) Command[testAppConfig, NoFlags] {
	var runErr error
	if len(err) > 0 {
		runErr = err[0]
	}

	return Command[testAppConfig, NoFlags]{
		Use: use,
		RunE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
			return runErr
		},
	}
}

func registerAndSetFlag[T any](
	t *testing.T,
	registry *FlagRegistry,
	cmd *cobra.Command,
	cfg *T,
	flagName, flagValue string,
) {
	t.Helper()

	registerAndParseFlags(t, registry, cmd, flagName, flagValue)

	err := registry.ParseFlags(cmd, cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func registerAndParseInvalidFlag[T any](
	t *testing.T,
	registry *FlagRegistry,
	cmd *cobra.Command,
	cfg *T,
	flagName, flagValue string,
) {
	t.Helper()

	registerAndParseFlags(t, registry, cmd, flagName, flagValue)

	err := registry.ParseFlags(cmd, cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func registerAndParseFlags(
	t *testing.T,
	registry *FlagRegistry,
	cmd *cobra.Command,
	flagName, flagValue string,
) {
	t.Helper()

	err := registry.RegisterFlags(cmd)
	if err != nil {
		t.Fatalf("expected no error registering flags, got: %v", err)
	}

	err = cmd.Flags().Set(flagName, flagValue)
	if err != nil {
		t.Fatalf("expected no error setting flag, got: %v", err)
	}
}

func noOpHandler() func(context.Context, *testConfig, NoFlags) error {
	return func(_ context.Context, _ *testConfig, _ NoFlags) error {
		return nil
	}
}
