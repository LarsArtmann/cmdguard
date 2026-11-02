package v4

import (
	"context"
	"encoding/json/v2"
	"errors"
	"strings"
	"testing"

	output "github.com/larsartmann/go-output"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

type explainDBConfig struct {
	Host string `flag:"db-host" default:"localhost" env:"DB_HOST" help:"Database host"`
}

type explainFixtureConfig struct {
	LogLevel string          `flag:"log-level" default:"info" env:"EXPLAIN_LOG_LEVEL" help:"Log level"`
	Port     int             `flag:"port"      default:"8080"                         help:"Port"`
	APIToken string          `flag:"api-token" default:""     env:"EXPLAIN_API_TOKEN" help:"API token" secret:"true"`
	Database explainDBConfig // no flag tag → nested config group
}

// newExplainCLI builds a CLI with a config file and a no-op leaf command.
func newExplainCLI(t *testing.T, configPath string, opts ...CLIOption) *CLI[explainFixtureConfig] {
	t.Helper()

	cliOpts := []CLIOption{WithFang(false)}
	if configPath != "" {
		cliOpts = append(cliOpts, WithConfigFile(configPath))
	}

	cliOpts = append(cliOpts, opts...)

	cli, err := NewCLI("app", "app", explainFixtureConfig{}, cliOpts...)
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"run",
		NoFlags{},
		func(_ context.Context, _ *explainFixtureConfig, _ NoFlags) error { return nil },
		WithShort("Run"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	return cli
}

// resolvedByName indexes an explanation by flag name.
func resolvedByName(explanation ConfigExplanation) map[string]ResolvedFlag {
	out := make(map[string]ResolvedFlag, len(explanation.Flags))
	for _, flag := range explanation.Flags {
		out[flag.Name] = flag
	}

	return out
}

func TestExplainConfig_DefaultSource(t *testing.T) {
	t.Parallel()

	cli := newExplainCLI(t, "")

	flags := resolvedByName(cli.ExplainConfig())
	port := flags["port"]
	testutil.AssertFieldEqString(t, string(port.Source), string(SourceDefault), "port source")
	testutil.AssertFieldEqString(t, port.Value, "8080", "port value")
	testutil.AssertFieldEqString(t, port.Default, "8080", "port declared default")
}

func TestExplainConfig_NestedKeyUsesCanonicalName(t *testing.T) {
	t.Parallel()

	cli := newExplainCLI(t, "")

	flags := resolvedByName(cli.ExplainConfig())
	host := flags["db-host"]
	testutil.AssertFieldEqString(t, host.Key, "Host", "db-host config key")
	testutil.AssertFieldEqString(t, host.Env, "DB_HOST", "db-host env")
}

func TestExplainConfig_ConfigFileSource(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "log-level: from-config\nport: 9000\n")

	cli := newExplainCLI(t, path, WithStrictConfigKeys())

	err := cli.ExecuteWithArgs(t.Context(), []string{"run"})
	testutil.AssertNoError(t, err)

	explanation := cli.ExplainConfig()
	flags := resolvedByName(explanation)

	testutil.AssertFieldEqString(t, string(flags["log-level"].Source), string(SourceConfigFile), "log-level source")
	testutil.AssertFieldEqString(t, flags["log-level"].Value, "from-config", "log-level value")
	testutil.AssertFieldEqString(t, flags["log-level"].Default, "info", "declared default survives")
	testutil.AssertFieldEqString(t, explanation.ConfigFile, path, "explanation config file")
}

func TestExplainConfig_FlagSource(t *testing.T) {
	t.Parallel()

	cli := newExplainCLI(t, "")

	err := cli.ExecuteWithArgs(t.Context(), []string{"run", "--log-level", "from-flag"})
	testutil.AssertNoError(t, err)

	flags := resolvedByName(cli.ExplainConfig())
	testutil.AssertFieldEqString(t, string(flags["log-level"].Source), string(SourceFlag), "log-level source")
	testutil.AssertFieldEqString(t, flags["log-level"].Value, "from-flag", "log-level value")
}

func TestExplainConfig_EnvironmentSource(t *testing.T) {
	t.Setenv("EXPLAIN_LOG_LEVEL", "from-env")

	cli := newExplainCLI(t, "")

	err := cli.ExecuteWithArgs(t.Context(), []string{"run"})
	testutil.AssertNoError(t, err)

	flags := resolvedByName(cli.ExplainConfig())
	testutil.AssertFieldEqString(t, string(flags["log-level"].Source), string(SourceEnvironment), "source")
	testutil.AssertFieldEqString(t, flags["log-level"].Value, "from-env", "value")
}

func TestExplainConfig_SecretIsAlwaysRedacted(t *testing.T) {
	t.Setenv("EXPLAIN_API_TOKEN", "super-secret-token")

	cli := newExplainCLI(t, "")

	err := cli.ExecuteWithArgs(t.Context(), []string{"run"})
	testutil.AssertNoError(t, err)

	explanation := cli.ExplainConfig()
	flags := resolvedByName(explanation)
	token := flags["api-token"]
	testutil.AssertBoolTrue(t, token.Secret, "api-token secret")
	testutil.AssertBoolTrue(t, token.Redacted, "api-token redacted")
	testutil.AssertFieldEqString(t, token.Value, RedactedValue, "redacted value")

	// The source is still reported: the point of redaction is to hide the
	// value, not to hide that the application is reading a secret from env.
	testutil.AssertFieldEqString(t, string(token.Source), string(SourceEnvironment), "token source")

	for _, flag := range explanation.Flags {
		if strings.Contains(flag.Value, "super-secret-token") {
			t.Fatalf("secret leaked in explanation of %q", flag.Name)
		}
	}
}

func TestExplainConfig_NestedFieldsAreExplained(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "database:\n  host: db.internal\n")

	cli := newExplainCLI(t, path, WithStrictConfigKeys())

	err := cli.ExecuteWithArgs(t.Context(), []string{"run"})
	testutil.AssertNoError(t, err)

	flags := resolvedByName(cli.ExplainConfig())
	testutil.AssertFieldEqString(t, string(flags["db-host"].Source), string(SourceConfigFile), "nested source")
	testutil.AssertFieldEqString(t, flags["db-host"].Value, "db.internal", "nested value")
}

func TestConfigExplainCommand_RegisteredOnDemand(t *testing.T) {
	t.Parallel()

	cli := newExplainCLI(t, "", WithConfigExplainCommand())

	var buf strings.Builder

	cli.RootCommand().SetOut(&buf)

	err := cli.ExecuteWithArgs(t.Context(), []string{"config", "explain"})
	testutil.AssertNoError(t, err)

	rendered := buf.String()
	for _, want := range []string{"FLAG", "SOURCE", "log-level", "default", "EXPLAIN_LOG_LEVEL"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("table output missing %q:\n%s", want, rendered)
		}
	}

	if strings.Contains(rendered, "super-secret") {
		t.Errorf("table output leaked a secret:\n%s", rendered)
	}
}

func TestConfigExplainCommand_YieldsToConsumerCommand(t *testing.T) {
	t.Parallel()

	cli := newExplainCLI(t, "", WithConfigExplainCommand())

	consumerRan := false

	own, err := NewCommand(
		configExplainCommandName,
		NoFlags{},
		func(_ context.Context, _ *explainFixtureConfig, _ NoFlags) error {
			consumerRan = true

			return nil
		},
		WithShort("Consumer-owned config command"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, own))

	err = cli.ExecuteWithArgs(t.Context(), []string{"config", "explain"})
	testutil.AssertNoError(t, err)
	testutil.AssertBoolTrue(t, consumerRan, "consumer-owned config command runs instead of cmdguard's")
}

func TestConfigExplainCommand_JSONOutput(t *testing.T) {
	t.Parallel()

	cli := newExplainCLI(t, "", WithConfigExplainCommand(), WithOutputFormat(output.FormatJSON))

	var buf strings.Builder

	cli.RootCommand().SetOut(&buf)

	err := cli.ExecuteWithArgs(t.Context(), []string{"config", "explain"})
	testutil.AssertNoError(t, err)

	var explanation ConfigExplanation

	err = json.Unmarshal([]byte(buf.String()), &explanation)
	testutil.AssertNoError(t, err)
	testutil.AssertFieldEqString(t, explanation.SchemaVersion, descriptionSchemaVersion, "schema version")
	testutil.AssertBoolTrue(t, len(explanation.Flags) > 0, "flags present")

	flags := resolvedByName(explanation)
	testutil.AssertFieldEqString(t, string(flags["port"].Source), string(SourceDefault), "port source")
}

func TestConfigExplainCommand_NotRegisteredByDefault(t *testing.T) {
	t.Parallel()

	cli := newExplainCLI(t, "")

	err := cli.ExecuteWithArgs(t.Context(), []string{"config", "explain"})
	testutil.AssertExpectedError(t, err)
	testutil.AssertBoolTrue(
		t,
		errors.Is(err, ErrInvalidCommand) || strings.Contains(err.Error(), "unknown command"),
		"config command absent",
	)
}

func TestExplainConfig_UnknownFieldIsNotInvented(t *testing.T) {
	t.Parallel()

	cli := newExplainCLI(t, "")

	for _, flag := range cli.ExplainConfig().Flags {
		if flag.Source == "" {
			t.Errorf("flag %q has no source", flag.Name)
		}
	}
}
