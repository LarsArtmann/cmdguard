package v4

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

var errCleanupRunFailed = errors.New("run failed")

var errCleanupHookFailed = errors.New("cleanup failed")

func TestWithCleanup_FiresOnSuccess(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `flag:"name" default:"world" help:"name"`
	}

	var (
		fired      bool
		seenRunErr error
		seenCfg    *cfg
	)

	cli, err := NewCLI(
		"test", "Test", cfg{},
		WithCleanup(func(_ *cobra.Command, c *cfg, runErr error) error {
			fired = true
			seenRunErr = runErr
			seenCfg = c

			return nil
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{Use: "sub", RunE: func(*cobra.Command, []string) error { return nil }}
	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub", "--name", "alice"})
	testutil.AssertNoError(t, err)

	if !fired {
		t.Fatal("cleanup hook did not fire after a successful RunE")
	}

	if seenRunErr != nil {
		t.Errorf("runErr = %v, want nil on success", seenRunErr)
	}

	if seenCfg == nil || seenCfg.Name != "alice" {
		t.Errorf("cfg is not the resolved config: got %+v", seenCfg)
	}
}

func TestWithCleanup_FiresOnRunEError(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	var received error

	cli, err := NewCLI(
		"test", "Test", cfg{},
		WithCleanup(func(_ *cobra.Command, _ *cfg, runErr error) error {
			received = runErr

			return nil
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{
		Use: "sub",
		RunE: func(*cobra.Command, []string) error {
			return errCleanupRunFailed
		},
	}
	cli.RootCommand().AddCommand(subCmd)

	_ = cli.ExecuteWithArgs(context.Background(), []string{"sub"})

	if !errors.Is(received, errCleanupRunFailed) {
		t.Errorf("cleanup received runErr = %v, want errCleanupRunFailed", received)
	}
}

func TestWithCleanup_MultipleHooksInOrder(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	var order []int

	cli, err := NewCLI(
		"test", "Test", cfg{},
		WithCleanup(func(*cobra.Command, *cfg, error) error {
			order = append(order, 1)

			return nil
		}),
		WithCleanup(func(*cobra.Command, *cfg, error) error {
			order = append(order, 2)

			return nil
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{Use: "sub", RunE: func(*cobra.Command, []string) error { return nil }}
	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub"})
	testutil.AssertNoError(t, err)

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Errorf("cleanup hooks ran in order %v, want [1 2]", order)
	}
}

func TestWithCleanup_HookErrorOnSuccessReturned(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	cli, err := NewCLI(
		"test", "Test", cfg{},
		WithCleanup(func(*cobra.Command, *cfg, error) error {
			return errCleanupHookFailed
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{Use: "sub", RunE: func(*cobra.Command, []string) error { return nil }}
	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub"})
	testutil.AssertExpectedError(t, err)

	if !errors.Is(err, errCleanupHookFailed) {
		t.Errorf("err = %v, want it to wrap errCleanupHookFailed", err)
	}
}

func TestWithCleanup_RunEErrorNotSwallowed(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	cli, err := NewCLI(
		"test", "Test", cfg{},
		WithCleanup(func(*cobra.Command, *cfg, error) error {
			return errCleanupHookFailed
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{
		Use: "sub",
		RunE: func(*cobra.Command, []string) error {
			return errCleanupRunFailed
		},
	}
	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub"})
	testutil.AssertExpectedError(t, err)

	if !errors.Is(err, errCleanupRunFailed) {
		t.Error("original RunE error was swallowed by cleanup")
	}

	if !errors.Is(err, errCleanupHookFailed) {
		t.Error("cleanup hook error was lost")
	}
}

// TestWithCleanup_FiresOncePerExecute guards against double-wrapping when
// Execute is called more than once on the same CLI: without the wiring guard
// the second Execute would wrap the already-wrapped RunE, firing cleanup twice
// per run (fires would be 3, not 2).
func TestWithCleanup_FiresOncePerExecute(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	var fires int

	cli, err := NewCLI(
		"test", "Test", cfg{},
		WithCleanup(func(*cobra.Command, *cfg, error) error {
			fires++

			return nil
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{Use: "sub", RunE: func(*cobra.Command, []string) error { return nil }}
	cli.RootCommand().AddCommand(subCmd)

	_ = cli.ExecuteWithArgs(context.Background(), []string{"sub"})
	_ = cli.ExecuteWithArgs(context.Background(), []string{"sub"})

	if fires != 2 {
		t.Errorf("cleanup fired %d times, want 2 (once per Execute, no double-wrap)", fires)
	}
}

func TestWithCleanup_FiresForRootRunE(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	var fired bool

	cli, err := NewCLI(
		"test", "Test", cfg{},
		WithCleanup(func(*cobra.Command, *cfg, error) error {
			fired = true

			return nil
		}),
	)
	testutil.AssertNoError(t, err)

	cli.RootCommand().RunE = func(*cobra.Command, []string) error { return nil }

	err = cli.ExecuteWithArgs(context.Background(), nil)
	testutil.AssertNoError(t, err)

	if !fired {
		t.Error("cleanup did not fire for the root command's RunE")
	}
}

func TestWithCleanup_NoHooksLeavesRunEUnchanged(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	ran := false

	cli, err := NewCLI("test", "Test", cfg{})
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{
		Use: "sub",
		RunE: func(*cobra.Command, []string) error {
			ran = true

			return nil
		},
	}
	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub"})
	testutil.AssertNoError(t, err)

	if !ran {
		t.Error("RunE did not run when no cleanup hooks are registered")
	}
}
