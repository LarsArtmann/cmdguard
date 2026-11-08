package v5_test

import (
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

type testCLIConfig struct {
	Verbose bool   `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
	Level   string `default:"info"  flag:"level"   help:"Log level"`
}

func TestNewCLI(t *testing.T) {
	t.Parallel()
	t.Run("creates CLI with defaults", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI("test", "Test CLI", testCLIConfig{})
		testutil.AssertNoError(t, err)

		testutil.AssertNotNil(t, cli)

		if cli.Name() != "test" {
			t.Errorf("Name() = %q, want %q", cli.Name(), "test")
		}

		if cli.Short() != "Test CLI" {
			t.Errorf("Short() = %q, want %q", cli.Short(), "Test CLI")
		}
	})

	t.Run("creates CLI with options", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI(
			"myapp",
			"My Application",
			testCLIConfig{},
			v5.WithCLIVersion("1.0.0"),
			v5.WithCLILong("This is a long description"),
		)
		testutil.AssertNoError(t, err)

		if cli.Name() != "myapp" {
			t.Errorf("Name() = %q, want %q", cli.Name(), "myapp")
		}

		if cli.Long() != "This is a long description" {
			t.Errorf("Long() = %q, want %q", cli.Long(), "This is a long description")
		}
	})

	t.Run("returns error for empty name", func(t *testing.T) {
		t.Parallel()

		_, err := v5.NewCLI("", "short", testCLIConfig{})
		testutil.AssertExpectedError(t, err)
	})

	t.Run("defaults SilenceUsage to true (no usage-on-error footgun)", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI("test", "Test CLI", testCLIConfig{})
		testutil.AssertNoError(t, err)

		testutil.AssertBoolTrue(t, cli.RootCommand().SilenceUsage, "default SilenceUsage")
	})

	t.Run("WithSilenceErrors sets SilenceErrors", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v5.WithSilenceErrors(),
		)
		testutil.AssertNoError(t, err)

		testutil.AssertBoolTrue(t, cli.RootCommand().SilenceErrors, "SilenceErrors")
	})

	t.Run("WithSilenceUsage sets SilenceUsage", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v5.WithSilenceUsage(),
		)
		testutil.AssertNoError(t, err)

		testutil.AssertBoolTrue(t, cli.RootCommand().SilenceUsage, "SilenceUsage")
	})

	t.Run("WithFang false disables fang", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v5.WithFang(false),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})

	t.Run("WithFang true keeps fang enabled", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v5.WithFang(true),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})

	t.Run("combines multiple options", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v5.WithSilenceErrors(),
			v5.WithSilenceUsage(),
			v5.WithCLIVersion("2.0.0"),
		)
		testutil.AssertNoError(t, err)

		testutil.AssertBoolTrue(t, cli.RootCommand().SilenceErrors, "SilenceErrors")
		testutil.AssertBoolTrue(t, cli.RootCommand().SilenceUsage, "SilenceUsage")

		rootCmd := cli.RootCommand()
		if rootCmd.Version != "2.0.0" {
			t.Errorf("Version = %q, want %q", rootCmd.Version, "2.0.0")
		}
	})

	t.Run("WithoutSilenceUsage re-enables usage on error", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v5.WithoutSilenceUsage(),
		)
		testutil.AssertNoError(t, err)

		testutil.AssertBoolFalse(
			t,
			cli.RootCommand().SilenceUsage,
			"WithoutSilenceUsage should set SilenceUsage to false",
		)
	})
}
