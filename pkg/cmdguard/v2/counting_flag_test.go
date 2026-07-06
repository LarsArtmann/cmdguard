package v2

import (
	"context"
	"testing"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

type countTestConfig struct {
	Quiet bool `default:"false" flag:"quiet" help:"Quiet mode"`
}

func TestCountingFlag_Integration(t *testing.T) {
	t.Parallel()

	type verbFlags struct {
		Verbose int `flag:"verbose" short:"v" help:"verbosity level" count:"true"`
	}

	tests := []struct {
		name        string
		args        []string
		wantVerbose int
	}{
		{"-vvv sets count to 3", []string{"run", "-vvv"}, 3},
		{"single -v sets count to 1", []string{"run", "-v"}, 1},
		{"no flag sets count to 0", []string{"run"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var result int

			cli, err := NewCLI[countTestConfig]("app", "test", countTestConfig{})
			testutil.AssertNoError(t, err)

			cmd, err := NewCommand(
				"run",
				&verbFlags{},
				func(_ context.Context, _ *countTestConfig, flags *verbFlags) error {
					result = flags.Verbose

					return nil
				},
				WithShort("Run"),
			)
			testutil.AssertNoError(t, err)
			testutil.AssertNoError(t, AddCommand(cli, cmd))

			err = cli.ExecuteWithArgs(t.Context(), tt.args)
			testutil.AssertNoError(t, err)
			if result != tt.wantVerbose {
				t.Errorf("verbose count = %d, want %d", result, tt.wantVerbose)
			}
		})
	}
}
