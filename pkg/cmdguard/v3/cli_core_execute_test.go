package v3_test

import (
	"context"
	"testing"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

func TestCLIExecute(t *testing.T) {
	t.Parallel()
	t.Run("executes command", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		executed := false
		cmd, err := v3.NewCommand(
			"run",
			v3.NoFlags{},
			func(_ context.Context, _ *testCLIConfig, _ v3.NoFlags) error {
				executed = true

				return nil
			},
			v3.WithShort("Run the command"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v3.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		ctx := t.Context()

		err = cli.ExecuteWithArgs(ctx, []string{"run"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if !executed {
			t.Error("command was not executed")
		}
	})
}
