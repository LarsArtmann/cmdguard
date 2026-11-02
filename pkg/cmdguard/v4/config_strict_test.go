package v4

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

type strictServerConfig struct {
	Host string `flag:"host" default:"localhost" help:"Server host"`
	Port int    `flag:"port" default:"8080"      help:"Server port"`
}

type strictFixtureConfig struct {
	LogLevel string               `flag:"log-level" default:"info" env:"LOG_LEVEL" help:"Log level"`
	Servers  []strictServerConfig // no flag tag → array of config groups
	Server   strictServerConfig   // no flag tag → nested config group
}

// writeConfigFile writes content to a uniquely named file inside a temp dir.
func writeConfigFile(t *testing.T, name, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	err := os.WriteFile(path, []byte(content), 0o600)
	testutil.AssertNoError(t, err)

	return path
}

// loadStrict loads a config file through the same path the CLI uses.
func loadStrict(t *testing.T, cfg any, path string, strict bool) ([]string, error) {
	t.Helper()

	loader := NewKoanfLoader(path)
	loader.SetStrictKeys(strict)

	return loader.Load(nil, cfg)
}

func TestStrictConfigKeys_UnknownTopLevelKey(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "log-leve1: debug\n")

	cfg := strictFixtureConfig{}

	_, err := loadStrict(t, &cfg, path, true)
	testutil.AssertExpectedError(t, err)
	testutil.AssertBoolTrue(t, errors.Is(err, ErrConfigUnknownKey), "errors.Is(ErrConfigUnknownKey)")

	if !strings.Contains(err.Error(), "log-leve1") {
		t.Errorf("error must name the offending key, got: %v", err)
	}

	if !strings.Contains(err.Error(), `did you mean "LogLevel"`) {
		t.Errorf("error should suggest the accepted spelling, got: %v", err)
	}
}

func TestStrictConfigKeys_ErrorNamesTheFile(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.json", `{"nope":1}`)

	cfg := strictFixtureConfig{}

	_, err := loadStrict(t, &cfg, path, true)
	testutil.AssertExpectedError(t, err)

	if !strings.Contains(err.Error(), path) {
		t.Errorf("error must name the file that contains the key, got: %v", err)
	}
}

func TestStrictConfigKeys_UnknownNestedKey(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "server:\n  hostt: localhost\n")

	cfg := strictFixtureConfig{}

	_, err := loadStrict(t, &cfg, path, true)
	testutil.AssertExpectedError(t, err)

	if !strings.Contains(strings.ToLower(err.Error()), "server.hostt") {
		t.Errorf("nested unknown key must carry its dotted path, got: %v", err)
	}
}

func TestStrictConfigKeys_UnknownKeyInArrayElement(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "servers:\n  - host: a\n    portt: 1\n")

	cfg := strictFixtureConfig{}

	_, err := loadStrict(t, &cfg, path, true)
	testutil.AssertExpectedError(t, err)

	if !strings.Contains(strings.ToLower(err.Error()), "servers[0].portt") {
		t.Errorf("array element key must carry its indexed path, got: %v", err)
	}
}

func TestStrictConfigKeys_AcceptsEveryFieldSpelling(t *testing.T) {
	t.Parallel()

	for _, key := range []string{"log-level", "LOG-LEVEL", "log_level", "LogLevel", "loglevel"} {
		path := writeConfigFile(t, "config.yaml", key+": debug\n")

		cfg := strictFixtureConfig{}

		setFields, err := loadStrict(t, &cfg, path, true)
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cfg.LogLevel, "debug", "LogLevel via "+key)
		testutil.AssertFieldLen(t, setFields, 1, "setFields via "+key)

		if setFields[0] != "LogLevel" {
			t.Errorf("setFields = %v, want [LogLevel] for key %q", setFields, key)
		}
	}
}

func TestStrictConfigKeys_TolerantModeStillLoads(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "log-level: debug\nnope: 1\n")

	cfg := strictFixtureConfig{}

	setFields, err := loadStrict(t, &cfg, path, false)
	testutil.AssertNoError(t, err)
	testutil.AssertFieldEqString(t, cfg.LogLevel, "debug", "LogLevel")
	testutil.AssertFieldLen(t, setFields, 1, "setFields")
}

func TestStrictConfigKeys_ReportsEveryUnknownKey(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "alpha: 1\nbeta: 2\n")

	cfg := strictFixtureConfig{}

	_, err := loadStrict(t, &cfg, path, true)
	testutil.AssertExpectedError(t, err)

	for _, key := range []string{"alpha", "beta"} {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("error should report %q, got: %v", key, err)
		}
	}
}

func TestStrictConfigKeys_NoKeySuggestionWhenNothingIsClose(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "completely-unrelated: 1\n")

	cfg := strictFixtureConfig{}

	_, err := loadStrict(t, &cfg, path, true)
	testutil.AssertExpectedError(t, err)

	if strings.Contains(err.Error(), "did you mean") {
		t.Errorf("no suggestion expected for an unrelated key, got: %v", err)
	}
}

func TestWithBatteries_EnablesStrictConfigKeys(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `flag:"name"`
	}

	plain, err := NewCLI("app", "app", cfg{})
	testutil.AssertNoError(t, err)
	testutil.AssertBoolFalse(t, plain.StrictConfigKeys(), "default StrictConfigKeys")

	batteries, err := NewCLI("app", "app", cfg{}, WithBatteries())
	testutil.AssertNoError(t, err)
	testutil.AssertBoolTrue(t, batteries.StrictConfigKeys(), "WithBatteries StrictConfigKeys")
}

func TestStrictConfigKeys_EndToEndFailsAtConstruction(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "log-leve1: debug\n")

	_, err := NewCLI(
		"app",
		"app",
		strictFixtureConfig{},
		WithFang(false),
		WithConfigFile(path),
		WithStrictConfigKeys(),
	)
	testutil.AssertExpectedError(t, err)
	testutil.AssertBoolTrue(t, errors.Is(err, ErrConfigUnknownKey), "errors.Is(ErrConfigUnknownKey)")
}

// TestStrictConfigKeys_ConfigValueSurvivesFlagParsing pins the split brain this
// workstream exists to remove: the config file populated the struct, but the
// set-field list did not see the key, so parsing fell back to the tag default
// and silently overwrote the file's value.
func TestStrictConfigKeys_ConfigValueSurvivesFlagParsing(t *testing.T) {
	t.Parallel()

	for _, key := range []string{"log-level", "LogLevel", "log_level"} {
		path := writeConfigFile(t, "config.yaml", key+": debug\n")

		cli, err := NewCLI(
			"app",
			"app",
			strictFixtureConfig{},
			WithFang(false),
			WithConfigFile(path),
			WithStrictConfigKeys(),
		)
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand(
			"run",
			NoFlags{},
			func(_ context.Context, _ *strictFixtureConfig, _ NoFlags) error { return nil },
			WithShort("Run"),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(t.Context(), []string{"run"})
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cli.Config().LogLevel, "debug", "config file value via "+key)
		testutil.AssertFieldEqString(t, cli.ConfigFilePath(), path, "ConfigFilePath")
	}
}

// TestConfigPrecedence_FlagBeatsEnvBeatsConfigBeatsDefault pins the four-layer
// order end to end so the provenance work has a verified contract to describe.
//
//nolint:paralleltest // t.Setenv cannot be combined with t.Parallel
func TestConfigPrecedence_FlagBeatsEnvBeatsConfigBeatsDefault(t *testing.T) {
	path := writeConfigFile(t, "config.yaml", "log-level: from-config\n")

	t.Run("config beats default", func(t *testing.T) {
		cli := newPrecedenceCLI(t, path)
		err := cli.ExecuteWithArgs(t.Context(), []string{"run"})
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cli.Config().LogLevel, "from-config", "log level")
	})

	t.Run("env beats config", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "from-env")

		cli := newPrecedenceCLI(t, path)
		err := cli.ExecuteWithArgs(t.Context(), []string{"run"})
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cli.Config().LogLevel, "from-env", "log level")
	})

	t.Run("flag beats env and config", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "from-env")

		cli := newPrecedenceCLI(t, path)
		err := cli.ExecuteWithArgs(t.Context(), []string{"run", "--log-level", "from-flag"})
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cli.Config().LogLevel, "from-flag", "log level")
	})
}

// newPrecedenceCLI builds a CLI with a config file and a no-op leaf command.
func newPrecedenceCLI(t *testing.T, configPath string) *CLI[strictFixtureConfig] {
	t.Helper()

	cli, err := NewCLI(
		"app",
		"app",
		strictFixtureConfig{},
		WithFang(false),
		WithConfigFile(configPath),
		WithStrictConfigKeys(),
	)
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"run",
		NoFlags{},
		func(_ context.Context, _ *strictFixtureConfig, _ NoFlags) error { return nil },
		WithShort("Run"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	return cli
}
