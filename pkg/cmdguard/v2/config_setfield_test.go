package v2

import (
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func TestSetField(t *testing.T) {
	t.Parallel()
	t.Run("set string field", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Name string
		}{}

		err := SetField(cfg, "Name", "test")
		testutil.AssertNoError(t, err)

		testutil.AssertFieldEqString(t, cfg.Name, "test", "Name")
	})

	t.Run("set int field", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Count int
		}{}

		err := SetField(cfg, "Count", 42)
		testutil.AssertNoError(t, err)

		testutil.AssertFieldEq(t, cfg.Count, 42, "Count")
	})

	t.Run("set bool field", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Enabled bool
		}{}

		err := SetField(cfg, "Enabled", true)
		testutil.AssertNoError(t, err)

		testutil.AssertBoolTrue(t, cfg.Enabled, "Enabled")
	})

	t.Run("non-pointer config", func(t *testing.T) {
		t.Parallel()

		cfg := struct{ Name string }{}

		err := SetField(cfg, "Name", "test")
		testutil.AssertExpectedError(t, err)

		assertErrorContains(t, err, "pointer to struct")
	})

	t.Run("field not found", func(t *testing.T) {
		t.Parallel()

		cfg := &struct{ Name string }{}

		err := SetField(cfg, "Missing", "test")
		testutil.AssertExpectedError(t, err)

		assertErrorContains(t, err, "not found")
	})

	t.Run("string to LogLevel", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Level LogLevel
		}{}

		err := SetField(cfg, "Level", "debug")
		testutil.AssertNoError(t, err)

		assertEnumString(t, cfg.Level.String(), "debug", "Level")
	})

	t.Run("string to LogFormat", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Format LogFormat
		}{}

		err := SetField(cfg, "Format", "json")
		testutil.AssertNoError(t, err)

		assertEnumString(t, cfg.Format.String(), "json", "Format")
	})

	t.Run("string to Duration", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Timeout Duration
		}{}

		err := SetField(cfg, "Timeout", "30s")
		testutil.AssertNoError(t, err)

		assertDurationField(t, cfg.Timeout, 30*time.Second)
	})

	t.Run("time.Duration to Duration", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Timeout Duration
		}{}

		err := SetField(cfg, "Timeout", 45*time.Second)
		testutil.AssertNoError(t, err)

		assertDurationField(t, cfg.Timeout, 45*time.Second)
	})

	t.Run("invalid LogLevel", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Level LogLevel
		}{}

		err := SetField(cfg, "Level", "invalid")
		testutil.AssertExpectedError(t, err)
	})

	t.Run("invalid LogFormat", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Format LogFormat
		}{}

		err := SetField(cfg, "Format", "invalid-format")
		testutil.AssertExpectedError(t, err)
	})

	t.Run("invalid Duration", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Timeout Duration
		}{}

		err := SetField(cfg, "Timeout", "not-a-duration")
		testutil.AssertExpectedError(t, err)
	})

	t.Run("incompatible types", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Name string
		}{}
		// Use a struct type which is truly incompatible with string
		err := SetField(cfg, "Name", Duration{duration: 5 * time.Minute})
		testutil.AssertExpectedError(t, err)

		assertErrorContains(t, err, "cannot convert")
	})

	t.Run("string to Enum with allowed values", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Mode Enum
		}{
			Mode: Enum{value: "dev", allowed: []string{"dev", "staging", "prod"}},
		}

		err := SetField(cfg, "Mode", "prod")
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cfg.Mode.String(), "prod", "Mode")
	})

	t.Run("string to Enum with invalid value rejects", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Mode Enum
		}{
			Mode: Enum{value: "dev", allowed: []string{"dev", "staging", "prod"}},
		}

		err := SetField(cfg, "Mode", "invalid")
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "invalid")
	})

	t.Run("string to Enum without allowed values accepts any", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Mode Enum
		}{}

		err := SetField(cfg, "Mode", "anything")
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cfg.Mode.String(), "anything", "Mode")
	})
}
