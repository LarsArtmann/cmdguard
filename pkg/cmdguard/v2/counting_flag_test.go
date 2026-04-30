package v2

import (
	"context"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

type countTestConfig struct {
	Quiet bool `default:"false" flag:"quiet" help:"Quiet mode"`
}

func TestCountingFlag_Integration(t *testing.T) {
	t.Parallel()

	t.Run("-vvv sets count to 3", func(t *testing.T) {
		t.Parallel()

		type verbFlags struct {
			Verbose int `flag:"verbose" short:"v" help:"verbosity level" count:"true"`
		}

		var result int

		cli, err := NewCLI[countTestConfig]("app", "test", countTestConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[countTestConfig, *verbFlags]("run",
			func(_ context.Context, _ *countTestConfig, flags *verbFlags) error {
				result = flags.Verbose

				return nil
			},
			WithShort[countTestConfig, *verbFlags]("Run"),
			WithFlags[countTestConfig, *verbFlags](&verbFlags{}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(t.Context(), []string{"run", "-vvv"})
		testutil.AssertNoError(t, err)
		if result != 3 {
			t.Errorf("verbose count = %d, want 3", result)
		}
	})

	t.Run("single -v sets count to 1", func(t *testing.T) {
		t.Parallel()

		type verbFlags struct {
			Verbose int `flag:"verbose" short:"v" help:"verbosity level" count:"true"`
		}

		var result int

		cli, err := NewCLI[countTestConfig]("app", "test", countTestConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[countTestConfig, *verbFlags]("run",
			func(_ context.Context, _ *countTestConfig, flags *verbFlags) error {
				result = flags.Verbose

				return nil
			},
			WithShort[countTestConfig, *verbFlags]("Run"),
			WithFlags[countTestConfig, *verbFlags](&verbFlags{}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(t.Context(), []string{"run", "-v"})
		testutil.AssertNoError(t, err)
		if result != 1 {
			t.Errorf("verbose count = %d, want 1", result)
		}
	})

	t.Run("no flag sets count to 0", func(t *testing.T) {
		t.Parallel()

		type verbFlags struct {
			Verbose int `flag:"verbose" short:"v" help:"verbosity level" count:"true"`
		}

		var result int

		cli, err := NewCLI[countTestConfig]("app", "test", countTestConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[countTestConfig, *verbFlags]("run",
			func(_ context.Context, _ *countTestConfig, flags *verbFlags) error {
				result = flags.Verbose

				return nil
			},
			WithShort[countTestConfig, *verbFlags]("Run"),
			WithFlags[countTestConfig, *verbFlags](&verbFlags{}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(t.Context(), []string{"run"})
		testutil.AssertNoError(t, err)
		if result != 0 {
			t.Errorf("verbose count = %d, want 0", result)
		}
	})
}
