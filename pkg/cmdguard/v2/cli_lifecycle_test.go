package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func noOpRunEForTestCLIConfigWithFlags[F any]() func(context.Context, *testCLIConfig, F) error {
	return func(_ context.Context, _ *testCLIConfig, _ F) error {
		return nil
	}
}

func TestCLISetLong(t *testing.T) {
	t.Parallel()
	t.Run("updates long description", func(t *testing.T) {
		t.Parallel()

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
	t.Parallel()
	t.Run("updates version", func(t *testing.T) {
		t.Parallel()

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
	t.Parallel()
	t.Run("adds global string flag", func(t *testing.T) {
		t.Parallel()

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
	t.Parallel()
	t.Run("adds global bool flag", func(t *testing.T) {
		t.Parallel()

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
	t.Parallel()
	t.Run("calls PreRunE and PostRunE", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		preRan, postRan := false, false

		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use: "test",
			PreRunE: func(_ context.Context, _ *testCLIConfig, _ v2.NoFlags) error {
				preRan = true

				return nil
			},
			RunE: noOpRunE[testCLIConfig],
			PostRunE: func(_ context.Context, _ *testCLIConfig, _ v2.NoFlags) error {
				postRan = true

				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
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

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type testFlags struct {
			Name string `default:"default" flag:"name" help:"name"`
		}

		var receivedName string

		cmd := v2.Command[testCLIConfig, testFlags]{
			Use:   "test",
			Flags: testFlags{},
			PreRunE: func(_ context.Context, _ *testCLIConfig, f testFlags) error {
				receivedName = f.Name

				return nil
			},
			RunE: noOpRunEForTestCLIConfigWithFlags[testFlags](),
		}

		err = v2.AddCommand(cli, cmd)
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

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type testFlags struct {
			Value string `default:"default" flag:"value" help:"value"`
		}

		var receivedValue string

		cmd := v2.Command[testCLIConfig, testFlags]{
			Use:   "test",
			Flags: testFlags{},
			RunE:  noOpRunEForTestCLIConfigWithFlags[testFlags](),
			PostRunE: func(_ context.Context, _ *testCLIConfig, f testFlags) error {
				receivedValue = f.Value

				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
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
	t.Parallel()
	t.Run("calls ExecuteAndExit successfully", func(t *testing.T) {
		t.Parallel()

		cli, _ := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})

		cmd := newTestCLICommand[testCLIConfig]("run")

		err := v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		cli.ExecuteAndExit(t.Context())
	})
}
