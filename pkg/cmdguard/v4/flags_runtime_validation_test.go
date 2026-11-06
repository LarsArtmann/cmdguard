package v4

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

func TestValidateFlags_StringMinMeansLength(t *testing.T) {
	t.Parallel()

	type lengthConfig struct {
		Name string `flag:"name" validate:"min=3"`
	}

	registry, cmd := setupFlagTest(t, lengthConfig{})

	setFlag(t, cmd, "name", "ab")

	err := registry.ValidateFlags(cmd)
	if !errors.Is(err, ErrValueTooShort) {
		t.Fatalf("expected ErrValueTooShort for a 2-character value, got: %v", err)
	}

	setFlag(t, cmd, "name", "abc")
	testutil.AssertNoError(t, registry.ValidateFlags(cmd))
}

// TestValidateFlags_DefaultValueIsValidated was retired: an invalid default
// now fails construction (TestNewFlagRegistry_RejectsInvalidEnumDefault in
// flag_defaults_test.go) instead of surviving to run time.

func TestValidateFlags_EnvValueIsValidated(t *testing.T) {
	type envConfig struct {
		Mode string `flag:"mode" values:"dev,prod" env:"VALIDATE_MODE"`
	}

	t.Setenv("VALIDATE_MODE", "bogus")

	registry, cmd := setupFlagTest(t, envConfig{})

	err := registry.ValidateFlags(cmd)
	if !errors.Is(err, ErrInvalidEnum) {
		t.Fatalf("expected an invalid env enum to fail validation, got: %v", err)
	}
}

func TestValidateFlags_EnvSatisfiesRequired(t *testing.T) {
	type requiredConfig struct {
		Token string `flag:"token" required:"true" env:"VALIDATE_TOKEN"`
	}

	t.Setenv("VALIDATE_TOKEN", "secret")

	registry, cmd := setupFlagTest(t, requiredConfig{})

	testutil.AssertNoError(t, registry.ValidateFlags(cmd))
}

func TestExecute_EnvSatisfiesRequired(t *testing.T) {
	type requiredRoot struct {
		Token string `flag:"token" required:"true" env:"EXECUTE_TOKEN" help:"Token"`
	}

	t.Setenv("EXECUTE_TOKEN", "from-env")

	cli, err := NewCLI("app", "App", requiredRoot{}, WithFang(false))
	testutil.AssertNoError(t, err)

	var gotToken string

	cmd, err := NewCommand(
		"run",
		NoFlags{},
		func(_ context.Context, cfg *requiredRoot, _ NoFlags) error {
			gotToken = cfg.Token

			return nil
		},
		WithShort("Run it"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	cli.RootCommand().SetErr(io.Discard)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"run"}))
	if gotToken != "from-env" {
		t.Fatalf("expected handler to see the env-sourced token, got %q", gotToken)
	}
}

func TestExecute_ConfigFileDoesNotSatisfyRequired(t *testing.T) {
	t.Parallel()

	type requiredRoot struct {
		Token string `flag:"token" required:"true" json:"token" help:"Token"`
	}

	configPath := filepath.Join(t.TempDir(), "config.json")
	const configJSON = `{"token": "from-file"}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0o600); err != nil {
		t.Fatal(err)
	}

	cli, err := NewCLI("app", "App", requiredRoot{}, WithConfigFile(configPath), WithFang(false))
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"run",
		NoFlags{},
		func(_ context.Context, _ *requiredRoot, _ NoFlags) error { return nil },
		WithShort("Run it"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	cli.RootCommand().SetErr(io.Discard)

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	if !errors.Is(err, ErrRequiredFlag) {
		t.Fatalf("expected config-file value NOT to satisfy required, got: %v", err)
	}

	cli.RootCommand().SetArgs([]string{"--token", "from-flag", "run"})
	testutil.AssertNoError(t, cli.Execute(context.Background()))
}

func TestExecute_RequiredFlagEnforced(t *testing.T) {
	t.Parallel()

	type requiredFlags struct {
		Name string `flag:"name" required:"true" help:"Name"`
	}

	cli, err := NewCLI("app", "App", testConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"run",
		requiredFlags{},
		func(_ context.Context, _ *testConfig, _ requiredFlags) error { return nil },
		WithShort("Run it"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	cli.RootCommand().SetErr(io.Discard)

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	if !errors.Is(err, ErrRequiredFlag) {
		t.Fatalf("expected ErrRequiredFlag, got: %v", err)
	}
}
