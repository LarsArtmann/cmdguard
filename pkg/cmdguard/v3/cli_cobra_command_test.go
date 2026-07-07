//nolint:fatcontext
package v3

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

func TestCLIToCobraCommand_DeeplyNested(t *testing.T) {
	t.Parallel()

	var executedCmd string

	leafCmd := Command[testAppConfig, NoFlags]{
		spec: commandSpec{use: "leaf"},
		runE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
			executedCmd = "leaf"

			return nil
		},
	}

	middleCmd := Command[testAppConfig, NoFlags]{
		spec:     commandSpec{use: "middle", long: "Middle level command"},
		commands: []Command[testAppConfig, NoFlags]{leafCmd},
	}

	topCmd := Command[testAppConfig, NoFlags]{
		spec:     commandSpec{use: "top", long: "Top level command"},
		commands: []Command[testAppConfig, NoFlags]{middleCmd},
	}

	cli, err := NewCLI[testAppConfig]("app", "App", testAppConfig{})
	testutil.AssertNoError(t, err)

	addCommand(t, cli, topCmd)

	err = cli.ExecuteWithArgs(t.Context(), []string{"top", "middle", "leaf"})
	testutil.AssertNoError(t, err)

	testutil.AssertFieldEqQuote(t, executedCmd, "leaf", "executedCmd")
}

func TestCLIToCobraCommand_PostRunEAfterSuccessfulRun(t *testing.T) {
	t.Parallel()

	var postRunCalled bool

	cli, err := NewCLI[testAppConfig]("app", "App", testAppConfig{})
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	cmd := Command[testAppConfig, NoFlags]{
		spec: commandSpec{
			use: "ok",
			postRunE: &typedHook[testAppConfig, NoFlags]{
				fn: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
					postRunCalled = true

					return nil
				},
			},
		},
		runE: noOpHandlerForTestAppConfig(),
	}

	addCommand(t, cli, cmd)

	err = cli.ExecuteWithArgs(t.Context(), []string{"ok"})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !postRunCalled {
		t.Error("PostRunE should be called after successful RunE")
	}
}

func TestCLIToCobraCommand_AllHooks(t *testing.T) {
	t.Parallel()

	var order []string

	cli, err := NewCLI[testAppConfig]("app", "App", testAppConfig{})
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	cmd := Command[testAppConfig, NoFlags]{
		spec: commandSpec{
			use:      "hooks",
			preRunE:  &typedHook[testAppConfig, NoFlags]{fn: makeHookRunE(&order, "pre")},
			postRunE: &typedHook[testAppConfig, NoFlags]{fn: makeHookRunE(&order, "post")},
		},
		runE: makeHookRunE(&order, "run"),
	}

	addCommand(t, cli, cmd)

	err = cli.ExecuteWithArgs(t.Context(), []string{"hooks"})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	want := []string{"pre", "run", "post"}
	if !slices.Equal(order, want) {
		t.Errorf("hook order = %v, want %v", order, want)
	}
}

func TestCLIToCobraCommand_NilContextFallsBack(t *testing.T) {
	t.Parallel()

	var gotCtx context.Context

	cli, err := NewCLI[testAppConfig]("app", "App", testAppConfig{})
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	cmd := Command[testAppConfig, NoFlags]{
		spec: commandSpec{use: "ctxcheck"},
		runE: func(ctx context.Context, _ *testAppConfig, _ NoFlags) error {
			gotCtx = ctx

			return nil
		},
	}

	addCommand(t, cli, cmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"ctxcheck"})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if gotCtx == nil {
		t.Error("expected non-nil context")
	}
}

func TestCLIToCobraCommand_SubcommandError(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI[testAppConfig]("app", "App", testAppConfig{})
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	invalidChild := Command[testAppConfig, NoFlags]{
		spec: commandSpec{use: ""}, // empty Use is invalid
		runE: noOpHandlerForTestAppConfig(),
	}

	parent := Command[testAppConfig, NoFlags]{
		spec:     commandSpec{use: "parent", long: "Parent command with invalid child"},
		commands: []Command[testAppConfig, NoFlags]{invalidChild},
	}

	err = AddCommand(cli, parent)
	if err == nil {
		t.Fatal("expected error for subcommand with empty Use")
	}

	if !errors.Is(err, ErrInvalidCommand) {
		t.Errorf("error = %v, want ErrInvalidCommand", err)
	}
}
