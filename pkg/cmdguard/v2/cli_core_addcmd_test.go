package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestCLIAddCommand(t *testing.T) {
	t.Parallel()
	t.Run("adds command with different flags type", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type greetFlags struct {
			Name string `default:"World" flag:"name" help:"Name to greet" short:"n"`
		}

		cmd, err := v2.NewCommand[testCLIConfig, greetFlags]("greet",
			func(_ context.Context, _ *testCLIConfig, _ greetFlags) error {
				return nil
			},
			v2.WithShort[testCLIConfig, greetFlags]("Greet someone"),
			v2.WithFlags[testCLIConfig, greetFlags](greetFlags{}),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}
	})

	t.Run("adds command with NoFlags", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICommandWithShort[testCLIConfig]("version", "Show version")

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}
	})

	t.Run("returns error for duplicate command", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICommandWithShort[testCLIConfig]("test", "Test command")

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("first AddCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, cmd)
		if err == nil {
			t.Fatal("expected error for duplicate command")
		}
	})

	t.Run("returns error for invalid command", func(t *testing.T) {
		t.Parallel()

		_, err := v2.NewCommand[testCLIConfig, v2.NoFlags]("",
			noOpRunE[testCLIConfig],
		)
		if err == nil {
			t.Fatal("expected error for invalid command")
		}
	})
}
