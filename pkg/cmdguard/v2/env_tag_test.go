package v2

import (
	"context"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

type envTestConfig struct {
	Debug bool `default:"false" flag:"debug" help:"Debug mode"`
}

func TestEnvTag_Integration(t *testing.T) {
	t.Run("env var provides default when flag not set", func(t *testing.T) {
		type dbFlags struct {
			Host string `flag:"host" env:"DB_HOST" default:"localhost" help:"Database host"`
		}

		var result string

		t.Setenv("DB_HOST", "db.example.com")

		cli, err := NewCLI[envTestConfig]("app", "test", envTestConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[envTestConfig, *dbFlags](
			"connect",
			func(_ context.Context, _ *envTestConfig, flags *dbFlags) error {
				result = flags.Host

				return nil
			},
			WithShort[envTestConfig, *dbFlags]("Connect"),
			WithFlags[envTestConfig, *dbFlags](&dbFlags{}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(t.Context(), []string{"connect"})
		testutil.AssertNoError(t, err)
		if result != "db.example.com" {
			t.Errorf("host = %q, want %q", result, "db.example.com")
		}
	})

	t.Run("explicit flag overrides env var", func(t *testing.T) {
		type dbFlags struct {
			Host string `flag:"host" env:"DB_HOST" default:"localhost" help:"Database host"`
		}

		var result string

		t.Setenv("DB_HOST", "db.example.com")

		cli, err := NewCLI[envTestConfig]("app", "test", envTestConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[envTestConfig, *dbFlags](
			"connect",
			func(_ context.Context, _ *envTestConfig, flags *dbFlags) error {
				result = flags.Host

				return nil
			},
			WithShort[envTestConfig, *dbFlags]("Connect"),
			WithFlags[envTestConfig, *dbFlags](&dbFlags{}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(
			t.Context(),
			[]string{"connect", "--host", "explicit.example.com"},
		)
		testutil.AssertNoError(t, err)
		if result != "explicit.example.com" {
			t.Errorf("host = %q, want %q", result, "explicit.example.com")
		}
	})

	t.Run("default used when env var not set", func(t *testing.T) {
		type dbFlags struct {
			Host string `flag:"host" env:"CMDGUARD_TEST_DB_HOST_UNUSED" default:"localhost" help:"Database host"`
		}

		var result string

		cli, err := NewCLI[envTestConfig]("app", "test", envTestConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[envTestConfig, *dbFlags](
			"connect",
			func(_ context.Context, _ *envTestConfig, flags *dbFlags) error {
				result = flags.Host

				return nil
			},
			WithShort[envTestConfig, *dbFlags]("Connect"),
			WithFlags[envTestConfig, *dbFlags](&dbFlags{}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(t.Context(), []string{"connect"})
		testutil.AssertNoError(t, err)
		if result != "localhost" {
			t.Errorf("host = %q, want %q (default)", result, "localhost")
		}
	})

	t.Run("env prefix is applied to command flags", func(t *testing.T) {
		type dbFlags struct {
			Host string `flag:"host" env:"DB_HOST" default:"localhost" help:"Database host"`
		}

		var result string

		t.Setenv("MYAPP_DB_HOST", "prefixed.example.com")

		cli, err := NewCLI[envTestConfig](
			"app", "test", envTestConfig{},
			WithEnvPrefix[envTestConfig]("MYAPP_"),
		)
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[envTestConfig, *dbFlags](
			"connect",
			func(_ context.Context, _ *envTestConfig, flags *dbFlags) error {
				result = flags.Host

				return nil
			},
			WithShort[envTestConfig, *dbFlags]("Connect"),
			WithFlags[envTestConfig, *dbFlags](&dbFlags{}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(t.Context(), []string{"connect"})
		testutil.AssertNoError(t, err)
		if result != "prefixed.example.com" {
			t.Errorf("host = %q, want %q", result, "prefixed.example.com")
		}
	})
}
