package v4_test

import (
	"testing"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func TestCLIAccessors(t *testing.T) {
	t.Parallel()
	t.Run("Scope returns DI scope", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli.Scope() == nil {
			t.Error("Scope() returned nil")
		}
	})

	t.Run("Config returns config pointer", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cfg := cli.Config()
		if cfg == nil {
			t.Error("Config() returned nil")
		}
	})

	t.Run("RootCommand returns cobra command", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		rootCmd := cli.RootCommand()
		if rootCmd == nil {
			t.Error("RootCommand() returned nil")

			return
		}

		if rootCmd.Use != "test" {
			t.Errorf("RootCommand().Use = %q, want %q", rootCmd.Use, "test")
		}
	})
}

func TestCLISubcommands(t *testing.T) {
	t.Parallel()
	t.Run("adds nested subcommands", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		subCmd := newTestCLICommandWithShort[testCLIConfig](t, "list", "List items")
		parentCmd := newTestParentCommand(
			t,
			"items",
			"Item management",
			"Item management long description",
			subCmd,
		)
		cmd := newTestParentCommand(
			t,
			"admin",
			"Admin commands",
			"Admin commands long description",
			parentCmd,
		)

		err = v4.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		rootCmd := cli.RootCommand()
		if len(rootCmd.Commands()) != 1 {
			t.Errorf("expected 1 command, got %d", len(rootCmd.Commands()))
		}
	})
}

func TestCLIInjector(t *testing.T) {
	t.Parallel()
	t.Run("returns DI injector", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli.Injector() == nil {
			t.Error("Injector() returned nil")
		}
	})
}
