package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestCLIExecute(t *testing.T) {
	t.Parallel()
	t.Run("executes command", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		executed := false
		cmd, err := v2.NewCommand[testCLIConfig, v2.NoFlags]("run",
			func(_ context.Context, _ *testCLIConfig, _ v2.NoFlags) error {
				executed = true

				return nil
			},
			v2.WithShort[testCLIConfig, v2.NoFlags]("Run the command"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, cmd)
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
