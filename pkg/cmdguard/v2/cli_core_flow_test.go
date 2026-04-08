package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestCLIFlowContext(t *testing.T) {
	t.Parallel()
	t.Run("FlowContext is nil before Execute", func(t *testing.T) {
		t.Parallel()
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli.FlowContext() != nil {
			t.Error("FlowContext should be nil before Execute")
		}
	})

	t.Run("FlowContext is set after Execute", func(t *testing.T) {
		t.Parallel()
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICmd("run")
		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"run"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if cli.FlowContext() == nil {
			t.Error("FlowContext should not be nil after Execute")
		}
	})

	t.Run("FlowContext is root context", func(t *testing.T) {
		t.Parallel()
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICmd("run")
		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"run"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		fc := cli.FlowContext()
		if fc == nil {
			t.Fatal("FlowContext is nil")
		}

		if !fc.IsRoot() {
			t.Error("FlowContext should be root after Execute")
		}
	})
}

func TestCLIFlowContextIntegration(t *testing.T) {
	t.Parallel()
	t.Run("command can access flow context", func(t *testing.T) {
		t.Parallel()
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		var accessedFlowCtx bool
		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:   "check",
			Short: "Check flow context",
			RunE: func(ctx context.Context, _ *testCLIConfig, _ v2.NoFlags) error {
				bfc, ok := v2.GetBranchingFlowContext(ctx)
				accessedFlowCtx = ok && bfc != nil
				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"check"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if !accessedFlowCtx {
			t.Error("command should be able to access flow context")
		}
	})

	t.Run("flow context has correct root name", func(t *testing.T) {
		t.Parallel()
		cli, err := v2.NewCLI[testCLIConfig]("myapp", "My App", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICmd("run")
		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"run"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		fc := cli.FlowContext()
		if fc == nil {
			t.Fatal("FlowContext is nil")
		}

		if fc.PathString() != "" {
			t.Errorf("root path should be empty, got %q", fc.PathString())
		}
	})
}
