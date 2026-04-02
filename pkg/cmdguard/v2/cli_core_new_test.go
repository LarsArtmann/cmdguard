package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type testCLIConfig struct {
	Verbose bool   `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
	Level   string `default:"info"  flag:"level"   help:"Log level"`
}

func TestNewCLI(t *testing.T) {
	t.Parallel()
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
