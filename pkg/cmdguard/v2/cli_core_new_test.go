package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
	"github.com/larsartmann/cmdguard/pkg/testutil"
)

type testCLIConfig struct {
	Verbose bool   `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
	Level   string `default:"info"  flag:"level"   help:"Log level"`
}

func TestNewCLI(t *testing.T) {
	t.Parallel()
	t.Run("creates CLI with defaults", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
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

		cli, err := v2.NewCLI[testCLIConfig](
			"myapp",
			"My Application",
			testCLIConfig{},
			v2.WithCLIVersion[testCLIConfig]("1.0.0"),
			v2.WithCLILong[testCLIConfig]("This is a long description"),
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

		_, err := v2.NewCLI[testCLIConfig]("", "short", testCLIConfig{})
		testutil.AssertExpectedError(t, err)
	})

	t.Run("WithSilenceErrors sets SilenceErrors", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig](
			"test", "Test CLI", testCLIConfig{},
			v2.WithSilenceErrors[testCLIConfig](),
		)
		testutil.AssertNoError(t, err)

		testutil.AssertBoolTrue(t, cli.RootCommand().SilenceErrors, "SilenceErrors")
	})

	t.Run("WithSilenceUsage sets SilenceUsage", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig](
			"test", "Test CLI", testCLIConfig{},
			v2.WithSilenceUsage[testCLIConfig](),
		)
		testutil.AssertNoError(t, err)

		testutil.AssertBoolTrue(t, cli.RootCommand().SilenceUsage, "SilenceUsage")
	})

	t.Run("WithColor false disables fang", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig](
			"test", "Test CLI", testCLIConfig{},
			v2.WithColor[testCLIConfig](false),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})

	t.Run("WithColor true keeps fang enabled", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig](
			"test", "Test CLI", testCLIConfig{},
			v2.WithColor[testCLIConfig](true),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})

	t.Run("combines multiple options", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig](
			"test", "Test CLI", testCLIConfig{},
			v2.WithSilenceErrors[testCLIConfig](),
			v2.WithSilenceUsage[testCLIConfig](),
			v2.WithCLIVersion[testCLIConfig]("2.0.0"),
		)
		testutil.AssertNoError(t, err)

		testutil.AssertBoolTrue(t, cli.RootCommand().SilenceErrors, "SilenceErrors")
		testutil.AssertBoolTrue(t, cli.RootCommand().SilenceUsage, "SilenceUsage")

		rootCmd := cli.RootCommand()
		if rootCmd.Version != "2.0.0" {
			t.Errorf("Version = %q, want %q", rootCmd.Version, "2.0.0")
		}
	})
}
