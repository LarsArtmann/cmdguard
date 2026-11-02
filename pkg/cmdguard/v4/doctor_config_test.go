package v4

import (
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

// newDoctorCLI builds a CLI with the doctor command explicitly registered.
func newDoctorCLI(
	t *testing.T,
	configPath string,
	opts ...DoctorOption[strictFixtureConfig],
) *CLI[strictFixtureConfig] {
	t.Helper()

	cliOpts := []CLIOption{WithFang(false)}
	if configPath != "" {
		cliOpts = append(cliOpts, WithConfigFile(configPath))
	}

	cli, err := NewCLI("app", "app", strictFixtureConfig{}, cliOpts...)
	testutil.AssertNoError(t, err)

	doctor, err := DoctorCommand(cli, opts...)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, doctor))

	return cli
}

// runDoctor executes the doctor command and returns its output and error.
func runDoctor(t *testing.T, cli *CLI[strictFixtureConfig]) (string, error) {
	t.Helper()

	var buf strings.Builder

	cli.RootCommand().SetOut(&buf)

	err := cli.ExecuteWithArgs(t.Context(), []string{"doctor"})

	return buf.String(), err
}

func TestDoctor_NoConfigFileConfiguredAddsNoConfigChecks(t *testing.T) {
	t.Parallel()

	cli := newDoctorCLI(t, "")

	out, err := runDoctor(t, cli)
	testutil.AssertNoError(t, err)

	if strings.Contains(out, DiagnosticConfigPath) {
		t.Errorf("doctor must stay quiet when no config file is configured:\n%s", out)
	}
}

func TestDoctor_CleanConfigPassesEveryConfigCheck(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "log-level: debug\n")

	cli := newDoctorCLI(t, path)

	out, err := runDoctor(t, cli)
	testutil.AssertNoError(t, err)

	for _, name := range []string{DiagnosticConfigPath, DiagnosticConfigParse, DiagnosticConfigKeys} {
		assertDoctorPassed(t, out, name)
	}
}

func TestDoctor_UnknownConfigKeyFailsTheKeysCheck(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "log-leve1: debug\n")

	// Tolerant mode at startup: the typo is ignored, which is exactly the
	// situation doctor has to surface.
	cli := newDoctorCLI(t, path)

	out, err := runDoctor(t, cli)
	testutil.AssertExpectedError(t, err)
	testutil.AssertBoolTrue(t, errors.Is(err, ErrDoctorFailed), "errors.Is(ErrDoctorFailed)")

	assertDoctorFailed(t, out, DiagnosticConfigKeys)

	if !strings.Contains(out, "log-leve1") {
		t.Errorf("doctor must name the offending key:\n%s", out)
	}

	assertDoctorPassed(t, out, DiagnosticConfigPath)
}

func TestDoctor_MissingConfigFileFailsThePathCheck(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI(
		"app",
		"app",
		strictFixtureConfig{},
		WithFang(false),
		WithConfigFile("/nonexistent/cmdguard-doctor-test/config.yaml"),
	)
	testutil.AssertNoError(t, err)

	doctor, err := DoctorCommand(cli)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, doctor))

	out, runErr := runDoctor(t, cli)
	testutil.AssertExpectedError(t, runErr)
	assertDoctorFailed(t, out, DiagnosticConfigPath)

	exitErr, ok := errors.AsType[*ExitError](runErr)
	if !ok {
		t.Fatalf("doctor failure must carry an exit code, got %T: %v", runErr, runErr)
	}

	if got := exitErr.ExitCode(); got != 1 {
		t.Errorf("doctor exit code = %d, want 1", got)
	}
}

func TestDoctor_SchemaValidatorReceivesTheResolvedPath(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "log-level: debug\n")

	var received string

	cli := newDoctorCLI(t, path, WithDoctorSchemaValidator[strictFixtureConfig](
		func(configPath string) error {
			received = configPath

			return nil
		},
	))

	out, err := runDoctor(t, cli)
	testutil.AssertNoError(t, err)
	testutil.AssertFieldEqString(t, received, path, "schema validator path")
	assertDoctorPassed(t, out, DiagnosticConfigSchema)
}

func TestDoctor_SchemaValidatorFailureIsReported(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "config.yaml", "log-level: debug\n")

	schemaErr := errors.New("config does not match schema")

	cli := newDoctorCLI(t, path, WithDoctorSchemaValidator[strictFixtureConfig](
		func(_ string) error { return schemaErr },
	))

	out, err := runDoctor(t, cli)
	testutil.AssertExpectedError(t, err)
	testutil.AssertBoolTrue(t, strings.Contains(out, schemaErr.Error()), "schema failure reported")
}

func TestDoctor_SchemaValidatorIsSkippedWithoutAConfigFile(t *testing.T) {
	t.Parallel()

	called := false

	cli, err := NewCLI(
		"app",
		"app",
		strictFixtureConfig{},
		WithFang(false),
		WithConfigFile("/nonexistent/cmdguard-doctor-test/config.yaml"),
	)
	testutil.AssertNoError(t, err)

	doctor, err := DoctorCommand(cli, WithDoctorSchemaValidator[strictFixtureConfig](
		func(_ string) error {
			called = true

			return nil
		},
	))
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, doctor))

	_, runErr := runDoctor(t, cli)
	testutil.AssertExpectedError(t, runErr)
	testutil.AssertBoolFalse(t, called, "schema validator without a config file")
}

// assertDoctorPassed checks that the named check is reported as passing.
func assertDoctorPassed(t *testing.T, out, name string) {
	t.Helper()

	if !strings.Contains(out, "✓ "+name) {
		t.Errorf("expected %q to pass, doctor output:\n%s", name, out)
	}
}

// assertDoctorFailed checks that the named check is reported as failing.
func assertDoctorFailed(t *testing.T, out, name string) {
	t.Helper()

	if !strings.Contains(out, "✗ "+name) {
		t.Errorf("expected %q to fail, doctor output:\n%s", name, out)
	}
}
