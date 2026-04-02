package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestCLIAccessors(t *testing.T) {
	t.Run("Scope returns DI scope", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli.Scope() == nil {
			t.Error("Scope() returned nil")
		}
	})

	t.Run("Config returns config pointer", func(t *testing.T) {
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
	t.Run("adds nested subcommands", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		subCmd := newTestCLICmd("list")
		subCmd.Short = "List items"

		parentCmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:      "items",
			Short:    "Item management",
			Commands: []v2.Command[testCLIConfig, v2.NoFlags]{subCmd},
		}

		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:      "admin",
			Short:    "Admin commands",
			Commands: []v2.Command[testCLIConfig, v2.NoFlags]{parentCmd},
		}

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
	t.Run("returns DI injector", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli.Injector() == nil {
			t.Error("Injector() returned nil")
		}
	})
}
