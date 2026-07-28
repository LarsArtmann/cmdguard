package v4_test

import (
	"context"
	"testing"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func TestCLIExecute(t *testing.T) {
	t.Parallel()
	t.Run("executes command", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		executed := false
		cmd, err := v4.NewCommand(
			"run",
			v4.NoFlags{},
			func(_ context.Context, _ *testCLIConfig, _ v4.NoFlags) error {
				executed = true

				return nil
			},
			v4.WithShort("Run the command"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v4.AddCommand(cli, cmd)
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
