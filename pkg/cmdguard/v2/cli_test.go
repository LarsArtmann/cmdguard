package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type testCLIConfig struct {
	Verbose bool   `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
	Level   string `default:"info"  flag:"level"   help:"Log level"`
}

func TestNewCLI(t *testing.T) {
	t.Run("creates CLI with defaults", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli == nil {
			t.Fatal("cli is nil")
		}

		if cli.Name() != "test" {
			t.Errorf("Name() = %q, want %q", cli.Name(), "test")
		}

		if cli.Short() != "Test CLI" {
			t.Errorf("Short() = %q, want %q", cli.Short(), "Test CLI")
		}
	})

	t.Run("creates CLI with options", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig](
			"myapp",
			"My Application",
			testCLIConfig{},
			v2.WithCLIVersion[testCLIConfig]("1.0.0"),
			v2.WithCLILong[testCLIConfig]("This is a long description"),
		)
		if err != nil {
			t.Fatalf("NewCLI with options failed: %v", err)
		}

		if cli.Name() != "myapp" {
			t.Errorf("Name() = %q, want %q", cli.Name(), "myapp")
		}

		if cli.Long() != "This is a long description" {
			t.Errorf("Long() = %q, want %q", cli.Long(), "This is a long description")
		}
	})

	t.Run("returns error for empty name", func(t *testing.T) {
		_, err := v2.NewCLI[testCLIConfig]("", "short", testCLIConfig{})
		if err == nil {
			t.Fatal("expected error for empty name")
		}
	})
}

func TestCLIAddCommand(t *testing.T) {
	t.Run("adds command with different flags type", func(t *testing.T) {
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
			RunE: func(ctx context.Context, cfg *testCLIConfig, flags greetFlags) error {
				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}
	})

	t.Run("adds command with NoFlags", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:   "version",
			Short: "Show version",
			RunE: func(ctx context.Context, cfg *testCLIConfig, flags v2.NoFlags) error {
				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}
	})

	t.Run("returns error for duplicate command", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:   "test",
			Short: "Test command",
			RunE: func(ctx context.Context, cfg *testCLIConfig, flags v2.NoFlags) error {
				return nil
			},
		}

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
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use: "", // invalid
		}

		err = v2.AddCommand(cli, cmd)
		if err == nil {
			t.Fatal("expected error for invalid command")
		}
	})
}

func TestCLIExecute(t *testing.T) {
	t.Run("executes command", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		executed := false
		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:   "run",
			Short: "Run the command",
			RunE: func(ctx context.Context, cfg *testCLIConfig, flags v2.NoFlags) error {
				executed = true

				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		ctx := context.Background()

		err = cli.ExecuteWithArgs(ctx, []string{"run"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if !executed {
			t.Error("command was not executed")
		}
	})
}

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

		subCmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:   "list",
			Short: "List items",
			RunE: func(ctx context.Context, cfg *testCLIConfig, flags v2.NoFlags) error {
				return nil
			},
		}

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

func TestCLISetConfig(t *testing.T) {
	t.Run("updates config", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cli.SetConfig(testCLIConfig{Verbose: true, Level: "debug"})

		cfg := cli.Config()
		if cfg == nil {
			t.Fatal("Config() returned nil")
		}

		if !cfg.Verbose {
			t.Error("Verbose not updated")
		}

		if cfg.Level != "debug" {
			t.Errorf("Level = %q, want %q", cfg.Level, "debug")
		}
	})
}

func TestCLIShutdown(t *testing.T) {
	t.Run("shutdown succeeds", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.Shutdown(context.Background())
		if err != nil {
			t.Errorf("Shutdown failed: %v", err)
		}
	})
}

func TestCLIHealthCheck(t *testing.T) {
	t.Run("health check succeeds", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.HealthCheck()
		if err != nil {
			t.Errorf("HealthCheck failed: %v", err)
		}
	})
}

func TestCLIHealthCheckWithContext(t *testing.T) {
	t.Run("health check with context succeeds", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.HealthCheckWithContext(context.Background())
		if err != nil {
			t.Errorf("HealthCheckWithContext failed: %v", err)
		}
	})
}

func TestCLISetLong(t *testing.T) {
	t.Run("updates long description", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cli.SetLong("new long description")

		if cli.Long() != "new long description" {
			t.Errorf("Long() = %q, want %q", cli.Long(), "new long description")
		}

		if cli.RootCommand().Long != "new long description" {
			t.Error("RootCommand().Long not updated")
		}
	})
}

func TestCLISetVersion(t *testing.T) {
	t.Run("updates version", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cli.SetVersion("2.0.0")

		if cli.RootCommand().Version != "2.0.0" {
			t.Errorf("RootCommand().Version = %q, want %q", cli.RootCommand().Version, "2.0.0")
		}
	})
}

func TestCLIAddGlobalFlag(t *testing.T) {
	t.Run("adds global string flag", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cli.AddGlobalFlag("custom", "c", "default", "help text")

		flag := cli.RootCommand().PersistentFlags().Lookup("custom")
		if flag == nil {
			t.Error("flag not added")
		}
	})
}

func TestCLIAddGlobalBoolFlag(t *testing.T) {
	t.Run("adds global bool flag", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cli.AddGlobalBoolFlag("debug", "d", true, "enable debug")

		flag := cli.RootCommand().PersistentFlags().Lookup("debug")
		if flag == nil {
			t.Error("flag not added")
		}
	})
}

func TestCLIPrePostRunE(t *testing.T) {
	t.Run("calls PreRunE and PostRunE", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		preRan, postRan := false, false

		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use: "test",
			PreRunE: func(ctx context.Context, cfg *testCLIConfig, f v2.NoFlags) error {
				preRan = true
				return nil
			},
			RunE: func(ctx context.Context, cfg *testCLIConfig, f v2.NoFlags) error {
				return nil
			},
			PostRunE: func(ctx context.Context, cfg *testCLIConfig, f v2.NoFlags) error {
				postRan = true
				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(context.Background(), []string{"test"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if !preRan {
			t.Error("PreRunE not called")
		}

		if !postRan {
			t.Error("PostRunE not called")
		}
	})
}

func TestCLIPreRunEWithFlags(t *testing.T) {
	t.Run("PreRunE receives parsed flags", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type testFlags struct {
			Name string `flag:"name" default:"default" help:"name"`
		}

		var receivedName string

		cmd := v2.Command[testCLIConfig, testFlags]{
			Use:   "test",
			Flags: testFlags{},
			PreRunE: func(ctx context.Context, cfg *testCLIConfig, f testFlags) error {
				receivedName = f.Name
				return nil
			},
			RunE: func(ctx context.Context, cfg *testCLIConfig, f testFlags) error {
				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(context.Background(), []string{"test", "--name", "custom"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if receivedName != "custom" {
			t.Errorf("receivedName = %q, want %q", receivedName, "custom")
		}
	})
}

func TestCLIPostRunEWithFlags(t *testing.T) {
	t.Run("PostRunE receives parsed flags", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type testFlags struct {
			Value string `flag:"value" default:"default" help:"value"`
		}

		var receivedValue string

		cmd := v2.Command[testCLIConfig, testFlags]{
			Use:   "test",
			Flags: testFlags{},
			RunE: func(ctx context.Context, cfg *testCLIConfig, f testFlags) error {
				return nil
			},
			PostRunE: func(ctx context.Context, cfg *testCLIConfig, f testFlags) error {
				receivedValue = f.Value
				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(context.Background(), []string{"test", "--value", "postvalue"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if receivedValue != "postvalue" {
			t.Errorf("receivedValue = %q, want %q", receivedValue, "postvalue")
		}
	})
}

func TestWithCLIScope(t *testing.T) {
	t.Run("sets custom scope", func(t *testing.T) {
		customScope := v2.NewScope("custom")
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{},
			v2.WithCLIScope[testCLIConfig](customScope))
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli == nil {
			t.Fatal("cli is nil")
		}
	})
}

func TestCLIExecuteAndExit(t *testing.T) {
	t.Run("calls ExecuteAndExit successfully", func(t *testing.T) {
		cli, _ := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:   "run",
			Short: "Run",
			RunE: func(ctx context.Context, cfg *testCLIConfig, f v2.NoFlags) error {
				return nil
			},
		}
		if err := v2.AddCommand(cli, cmd); err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}
		cli.ExecuteAndExit(context.Background())
	})
}
