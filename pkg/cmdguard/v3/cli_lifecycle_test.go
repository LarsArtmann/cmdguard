package v3_test

import (
	"context"
	"os"
	"testing"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

func TestCLISetLong(t *testing.T) {
	t.Parallel()
	t.Run("updates long description", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
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
	t.Parallel()
	t.Run("updates version", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
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
	t.Parallel()
	t.Run("adds global string flag", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cli.AddGlobalFlag("custom", "c", "default", "help text")

		flag := cli.RootCommand().PersistentFlags().Lookup("custom")
		if flag == nil {
			t.Fatal("flag not added")
		}

		if flag.Shorthand != "c" {
			t.Errorf("Shorthand = %q, want %q", flag.Shorthand, "c")
		}

		if flag.DefValue != "default" {
			t.Errorf("DefValue = %q, want %q", flag.DefValue, "default")
		}

		if flag.Usage != "help text" {
			t.Errorf("Usage = %q, want %q", flag.Usage, "help text")
		}
	})
}

func TestCLIAddGlobalBoolFlag(t *testing.T) {
	t.Parallel()
	t.Run("adds global bool flag", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cli.AddGlobalBoolFlag("debug", "d", true, "enable debug")

		flag := cli.RootCommand().PersistentFlags().Lookup("debug")
		if flag == nil {
			t.Fatal("flag not added")
		}

		if flag.Shorthand != "d" {
			t.Errorf("Shorthand = %q, want %q", flag.Shorthand, "d")
		}

		if flag.DefValue != "true" {
			t.Errorf("DefValue = %q, want %q", flag.DefValue, "true")
		}

		if flag.Usage != "enable debug" {
			t.Errorf("Usage = %q, want %q", flag.Usage, "enable debug")
		}
	})
}

func TestCLIPrePostRunE(t *testing.T) {
	t.Parallel()
	t.Run("calls PreRunE and PostRunE", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		preRan, postRan := false, false

		cmd, err := v3.NewCommand(
			"test",
			v3.NoFlags{},
			noOpRunE[testCLIConfig],
			v3.WithPreRunE(
				func(_ context.Context, _ *testCLIConfig, _ v3.NoFlags) error {
					preRan = true

					return nil
				},
			),
			v3.WithPostRunE(
				func(_ context.Context, _ *testCLIConfig, _ v3.NoFlags) error {
					postRan = true

					return nil
				},
			),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v3.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"test"})
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
	t.Parallel()
	t.Run("PreRunE receives parsed flags", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type testFlags struct {
			Name string `default:"default" flag:"name" help:"name"`
		}

		var receivedName string

		cmd, err := v3.NewCommand(
			"test",
			testFlags{},
			NoOpRunEWithFlags[testCLIConfig, testFlags](),
			v3.WithPreRunE(
				func(_ context.Context, _ *testCLIConfig, f testFlags) error {
					receivedName = f.Name

					return nil
				},
			),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v3.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"test", "--name", "custom"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if receivedName != "custom" {
			t.Errorf("receivedName = %q, want %q", receivedName, "custom")
		}
	})
}

func TestCLIPostRunEWithFlags(t *testing.T) {
	t.Parallel()
	t.Run("PostRunE receives parsed flags", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type testFlags struct {
			Value string `default:"default" flag:"value" help:"value"`
		}

		var receivedValue string

		cmd, err := v3.NewCommand(
			"test",
			testFlags{},
			NoOpRunEWithFlags[testCLIConfig, testFlags](),
			v3.WithPostRunE(
				func(_ context.Context, _ *testCLIConfig, f testFlags) error {
					receivedValue = f.Value

					return nil
				},
			),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v3.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"test", "--value", "postvalue"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if receivedValue != "postvalue" {
			t.Errorf("receivedValue = %q, want %q", receivedValue, "postvalue")
		}
	})
}

func TestWithCLIScope(t *testing.T) {
	t.Parallel()
	t.Run("sets custom scope", func(t *testing.T) {
		t.Parallel()

		customScope := v3.NewScope("custom")

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{},
			v3.WithCLIScope(customScope))
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli == nil {
			t.Fatal("cli is nil")
		}

		if cli.Scope() != customScope {
			t.Error("cli.Scope() should be the custom scope passed via WithCLIScope")
		}
	})
}

func TestCLIExecuteAndExit(t *testing.T) {
	t.Parallel()
	t.Run("calls ExecuteAndExit successfully", func(t *testing.T) {
		t.Parallel()

		cli, _ := v3.NewCLI("test", "Test", testCLIConfig{})

		cmd := newTestCLICommand[testCLIConfig](t, "run")

		err := v3.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		cli.ExecuteAndExit(t.Context())
	})
}

func TestCLINoColor(t *testing.T) {
	t.Parallel()

	t.Run("--no-color flag is registered", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		flag := cli.RootCommand().PersistentFlags().Lookup("no-color")
		if flag == nil {
			t.Fatal("no-color flag not registered")
		}

		if flag.DefValue != "false" {
			t.Errorf("DefValue = %q, want %q", flag.DefValue, "false")
		}

		if flag.Usage == "" {
			t.Error("Usage should not be empty")
		}
	})

	t.Run("NoColor returns false by default", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		if cli.NoColor() {
			t.Error("NoColor() should return false when --no-color is not set")
		}
	})

	t.Run("NoColor returns true after --no-color is parsed", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{},
			v3.WithFang(false))
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICommand[testCLIConfig](t, "run")
		addCommand(t, cli, cmd)

		err = cli.ExecuteWithArgs(t.Context(), []string{"run", "--no-color"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if !cli.NoColor() {
			t.Error("NoColor() should return true after --no-color is parsed")
		}
	})
}

func TestCLINoColorEnvVar(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cli, err := v3.NewCLI("test", "Test", testCLIConfig{},
		v3.WithFang(false))
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	cmd := newTestCLICommand[testCLIConfig](t, "run")
	addCommand(t, cli, cmd)

	err = cli.ExecuteWithArgs(t.Context(), []string{"run"})
	if err != nil {
		t.Fatalf("ExecuteWithArgs failed: %v", err)
	}
}

func TestCLINoColorRestoresEnvVar(t *testing.T) {
	//nolint:paralleltest // mutates process-wide env var
	t.Run("restores NO_COLOR after execution", func(t *testing.T) {
		os.Unsetenv("NO_COLOR")

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{},
			v3.WithFang(false))
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICommand[testCLIConfig](t, "run")
		addCommand(t, cli, cmd)

		err = cli.ExecuteWithArgs(t.Context(), []string{"--no-color", "run"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		_, set := os.LookupEnv("NO_COLOR")
		if set {
			t.Error("NO_COLOR should be unset after execution (was not set before)")
		}
	})

	t.Run("restores previous NO_COLOR value", func(t *testing.T) {
		t.Setenv("NO_COLOR", "0")

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{},
			v3.WithFang(false))
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := newTestCLICommand[testCLIConfig](t, "run")
		addCommand(t, cli, cmd)

		err = cli.ExecuteWithArgs(t.Context(), []string{"--no-color", "run"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		got := os.Getenv("NO_COLOR")
		if got != "0" {
			t.Errorf("NO_COLOR = %q, want %q", got, "0")
		}
	})
}
