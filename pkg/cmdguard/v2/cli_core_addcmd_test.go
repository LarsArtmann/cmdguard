package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func newTestCLICmd(use string) v2.Command[testCLIConfig, v2.NoFlags] {
	return v2.Command[testCLIConfig, v2.NoFlags]{
		Use: use,
		RunE: func(_ context.Context, _ *testCLIConfig, _ v2.NoFlags) error {
			return nil
		},
	}
}

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

		cmd := v2.Command[testCLIConfig, greetFlags]{
			Use:   "greet",
			Short: "Greet someone",
			Flags: greetFlags{},
			RunE: func(_ context.Context, _ *testCLIConfig, _ greetFlags) error {
				return nil
			},
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

		cmd := newTestCLICmd("version")
		cmd.Short = "Show version"

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

		cmd := newTestCLICmd("test")
		cmd.Short = "Test command"

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
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use: "",
		}

		err = v2.AddCommand(cli, cmd)
		if err == nil {
			t.Fatal("expected error for invalid command")
		}
	})
}
