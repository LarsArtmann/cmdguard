package v4

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

func buildFlagCommand(t *testing.T, registry *FlagRegistry) *cobra.Command {
	t.Helper()

	cmd := &cobra.Command{Use: "test"}
	registerFlags(t, registry, cmd)

	return cmd
}

func MustNewFlagRegistryForTest[T any](t *testing.T, config T) *FlagRegistry {
	t.Helper()

	registry, err := NewFlagRegistry(config)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	return registry
}

// Construction-time default validation (M07): an impossible default must fail
// NewFlagRegistry — where the mistake is made — not the first run.

func TestNewFlagRegistry_RejectsInvalidEnumDefault(t *testing.T) {
	t.Parallel()

	type enumDefault struct {
		Mode string `flag:"mode" default:"bogus" values:"dev,prod"`
	}

	_, err := NewFlagRegistry(enumDefault{})
	if !errors.Is(err, ErrInvalidFlagDefault) {
		t.Fatalf("expected ErrInvalidFlagDefault for an out-of-enum default, got: %v", err)
	}

	if !errors.Is(err, ErrInvalidEnum) {
		t.Fatalf("expected the enum failure to stay identifiable, got: %v", err)
	}
}

func TestNewFlagRegistry_RejectsOutOfRangeNumericDefault(t *testing.T) {
	t.Parallel()

	type rangeDefault struct {
		Port int `flag:"port" default:"99999" validate:"min=1,max=65535"`
	}

	_, err := NewFlagRegistry(rangeDefault{})
	if !errors.Is(err, ErrInvalidFlagDefault) {
		t.Fatalf("expected ErrInvalidFlagDefault for an out-of-range default, got: %v", err)
	}

	if !errors.Is(err, ErrValueTooLarge) {
		t.Fatalf("expected the range failure to stay identifiable, got: %v", err)
	}
}

func TestNewFlagRegistry_RejectsRegexMismatchedDefault(t *testing.T) {
	t.Parallel()

	type regexDefault struct {
		ID string `flag:"id" default:"abc" validate:"regex=^[0-9]+$"`
	}

	_, err := NewFlagRegistry(regexDefault{})
	if !errors.Is(err, ErrInvalidFlagDefault) {
		t.Fatalf("expected ErrInvalidFlagDefault for a regex-mismatched default, got: %v", err)
	}

	if !errors.Is(err, ErrValuePatternMismatch) {
		t.Fatalf("expected the pattern failure to stay identifiable, got: %v", err)
	}
}

func TestNewFlagRegistry_RejectsTooShortStringDefault(t *testing.T) {
	t.Parallel()

	type lengthDefault struct {
		Name string `flag:"name" default:"ab" validate:"min=3"`
	}

	_, err := NewFlagRegistry(lengthDefault{})
	if !errors.Is(err, ErrInvalidFlagDefault) {
		t.Fatalf("expected ErrInvalidFlagDefault for a too-short default, got: %v", err)
	}

	if !errors.Is(err, ErrValueTooShort) {
		t.Fatalf("expected the length failure to stay identifiable, got: %v", err)
	}
}

func TestNewFlagRegistry_RejectsInvalidDurationDefaultBounds(t *testing.T) {
	t.Parallel()

	type durationDefault struct {
		Wait Duration `flag:"wait" default:"90s" validate:"min=10s,max=1m"`
	}

	_, err := NewFlagRegistry(durationDefault{})
	if !errors.Is(err, ErrInvalidFlagDefault) {
		t.Fatalf("expected ErrInvalidFlagDefault for an out-of-bounds duration default, got: %v", err)
	}

	if !errors.Is(err, ErrValueTooLarge) {
		t.Fatalf("expected the duration-bound failure to stay identifiable, got: %v", err)
	}
}

func TestNewFlagRegistry_AcceptsValidDefaults(t *testing.T) {
	t.Parallel()

	type validDefaults struct {
		Mode  string   `flag:"mode"  default:"dev"      values:"dev,prod"`
		Port  int      `flag:"port"  default:"8080"     validate:"min=1,max=65535"`
		ID    string   `flag:"id"    default:"123"      validate:"regex=^[0-9]+$"`
		Wait  Duration `flag:"wait"  default:"30s"      validate:"min=10s,max=1m"`
		Plain string   `flag:"plain" default:"anything"`
	}

	_, err := NewFlagRegistry(validDefaults{})
	testutil.AssertNoError(t, err)
}

func TestNewFlagRegistry_LateRegisteredValidatorKeepsWorking(t *testing.T) {
	t.Parallel()

	type lateValidator struct {
		Name string `flag:"name" default:"x" validate:"custom"`
	}

	registry, err := NewFlagRegistry(lateValidator{})
	testutil.AssertNoError(t, err)

	registry.RegisterFlagValidator("custom", func(_ string) error { return nil })

	cmd := buildFlagCommand(t, registry)
	testutil.AssertNoError(t, registry.ValidateFlags(cmd))
}

// Parsed-type bounds (M08): min/max resolve in the flag type's domain, so a
// duration never compares nanoseconds against a seconds bound and textual
// types fall back to length semantics instead of failing to parse as numbers.

func TestValidateFlags_DurationBoundsUseDurationDomain(t *testing.T) {
	t.Parallel()

	type durationFlags struct {
		Wait Duration `flag:"wait" validate:"min=10s,max=1m"`
	}

	registry := MustNewFlagRegistryForTest(t, durationFlags{})

	cmd := buildFlagCommand(t, registry)

	setFlag(t, cmd, "wait", "5s")
	if err := registry.ValidateFlags(cmd); !errors.Is(err, ErrValueTooSmall) {
		t.Fatalf("expected ErrValueTooSmall for a 5s wait against min=10s, got: %v", err)
	}

	setFlag(t, cmd, "wait", "90s")
	if err := registry.ValidateFlags(cmd); !errors.Is(err, ErrValueTooLarge) {
		t.Fatalf("expected ErrValueTooLarge for a 90s wait against max=1m, got: %v", err)
	}

	setFlag(t, cmd, "wait", "30s")
	testutil.AssertNoError(t, registry.ValidateFlags(cmd))
}

func TestValidateFlags_TimeDurationBoundsUseDurationDomain(t *testing.T) {
	t.Parallel()

	type goDurationFlags struct {
		Wait time.Duration `flag:"wait" validate:"min=10s,max=1m"`
	}

	registry := MustNewFlagRegistryForTest(t, goDurationFlags{})
	registry.RegisterGoDurationHandler()

	cmd := buildFlagCommand(t, registry)

	setFlag(t, cmd, "wait", "5s")
	if err := registry.ValidateFlags(cmd); !errors.Is(err, ErrValueTooSmall) {
		t.Fatalf("expected ErrValueTooSmall for a 5s wait against min=10s, got: %v", err)
	}

	setFlag(t, cmd, "wait", "90s")
	if err := registry.ValidateFlags(cmd); !errors.Is(err, ErrValueTooLarge) {
		t.Fatalf("expected ErrValueTooLarge for a 90s wait against max=1m, got: %v", err)
	}

	setFlag(t, cmd, "wait", "30s")
	testutil.AssertNoError(t, registry.ValidateFlags(cmd))
}

func TestValidateFlags_TextualParsedTypeBoundsMeanLength(t *testing.T) {
	t.Parallel()

	type urlFlags struct {
		Site URL `flag:"site" validate:"min=5,max=20"`
	}

	registry := MustNewFlagRegistryForTest(t, urlFlags{})

	cmd := buildFlagCommand(t, registry)

	setFlag(t, cmd, "site", "https://ex.com")
	testutil.AssertNoError(t, registry.ValidateFlags(cmd))

	setFlag(t, cmd, "site", "https://example.com/very/long/path")
	if err := registry.ValidateFlags(cmd); !errors.Is(err, ErrValueTooLong) {
		t.Fatalf("expected ErrValueTooLong for an over-long URL, got: %v", err)
	}

	setFlag(t, cmd, "site", "h.io")
	if err := registry.ValidateFlags(cmd); !errors.Is(err, ErrValueTooShort) {
		t.Fatalf("expected ErrValueTooShort for a too-short URL, got: %v", err)
	}
}

func TestValidateConfig_DurationBoundsUseDurationDomain(t *testing.T) {
	t.Parallel()

	type durationConfig struct {
		Wait Duration `flag:"wait" validate:"min=10s,max=1m"`
	}

	cfg := durationConfig{Wait: FromDuration(5 * time.Second)}
	if err := ValidateConfig(cfg); !errors.Is(err, ErrValueTooSmall) {
		t.Fatalf("expected ErrValueTooSmall validating a 5s config value, got: %v", err)
	}

	ok := durationConfig{Wait: FromDuration(30 * time.Second)}
	testutil.AssertNoError(t, ValidateConfig(ok))
}
