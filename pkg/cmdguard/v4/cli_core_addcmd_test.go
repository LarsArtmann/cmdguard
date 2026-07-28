package v4_test

import (
	"context"
	"testing"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func TestCLIAddCommand(t *testing.T) {
	t.Parallel()
	t.Run("adds command with different flags type", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type greetFlags struct {
			Name string `default:"World" flag:"name" help:"Name to greet" short:"n"`
		}

		cmd, err := v4.NewCommand(
			"greet",
			greetFlags{},
			func(_ context.Context, _ *testCLIConfig, _ greetFlags) error {
				return nil
			},
			v4.WithShort("Greet someone"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v4.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}
	})

	t.Run("adds command with NoFlags", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICommandWithShort[testCLIConfig](t, "version", "Show version")

		err = v4.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}
	})

	t.Run("returns error for duplicate command", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICommandWithShort[testCLIConfig](t, "test", "Test command")

		err = v4.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("first AddCommand failed: %v", err)
		}

		err = v4.AddCommand(cli, cmd)
		if err == nil {
			t.Fatal("expected error for duplicate command")
		}
	})
}
