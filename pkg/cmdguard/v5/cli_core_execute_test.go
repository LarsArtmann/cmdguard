package v5_test

import (
	"context"
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func TestCLIExecute(t *testing.T) {
	t.Parallel()
	t.Run("executes command", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		executed := false
		cmd, err := v5.NewCommand(
			"run",
			v5.NoFlags{},
			func(_ context.Context, _ *testCLIConfig, _ v5.NoFlags) error {
				executed = true

				return nil
			},
			v5.WithShort("Run the command"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v5.AddCommand(cli, cmd)
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
