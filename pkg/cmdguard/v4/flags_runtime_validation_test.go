package v4

import (
	"context"
	"errors"
	"io"
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

func TestValidateFlags_DefaultValueIsValidated(t *testing.T) {
	t.Parallel()

	type defaultConfig struct {
		Mode string `flag:"mode" default:"bogus" values:"dev,prod"`
	}

	registry, cmd := setupFlagTest(t, defaultConfig{})

	err := registry.ValidateFlags(cmd)
	if !errors.Is(err, ErrInvalidEnum) {
		t.Fatalf("expected an invalid default enum to fail validation, got: %v", err)
	}
}

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
