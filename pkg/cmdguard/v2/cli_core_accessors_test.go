package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

func TestCLIAccessors(t *testing.T) {
	t.Parallel()
	t.Run("Scope returns DI scope", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli.Scope() == nil {
			t.Error("Scope() returned nil")
		}
	})

	t.Run("Config returns config pointer", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
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

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
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

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		subCmd := newTestCLICommandWithShort[testCLIConfig]("list", "List items")
		parentCmd := newTestParentCommand[testCLIConfig](
			"items",
			"Item management",
			"Item management long description",
			subCmd,
		)
		cmd := newTestParentCommand[testCLIConfig](
			"admin",
			"Admin commands",
			"Admin commands long description",
			parentCmd,
		)

		err = v2.AddCommand(cli, cmd)
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

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli.Injector() == nil {
			t.Error("Injector() returned nil")
		}
	})
}
