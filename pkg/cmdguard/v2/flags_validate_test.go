package v2

import (
	"testing"

	"github.com/spf13/cobra"
)

func setupFlagTest[T any](t *testing.T, config T) (*FlagRegistry, *cobra.Command) {
	t.Helper()

	registry, err := NewFlagRegistry(config)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	cmd := &cobra.Command{Use: "test"}
	if err := registry.RegisterFlags(cmd); err != nil {
		t.Fatalf("expected no error registering flags, got: %v", err)
	}

	return registry, cmd
}

func setFlagAndAssertValid(
	t *testing.T,
	cmd *cobra.Command,
	registry *FlagRegistry,
	flagName, flagValue string,
) {
	t.Helper()

	err := cmd.Flags().Set(flagName, flagValue)
	if err != nil {
		t.Fatalf("expected no error setting flag, got: %v", err)
	}

	err = registry.ValidateFlags(cmd)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestFlagRegistry_ValidateFlags(t *testing.T) {
	t.Parallel()
	t.Run("valid values pass", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Mode string `default:"dev" flag:"mode" values:"dev,staging,prod"`
		}

		registry, cmd := setupFlagTest(t, TestConfig{})

		setFlag(t, cmd, "mode", "staging")

		err := registry.ValidateFlags(cmd)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("invalid value returns error", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Mode string `default:"dev" flag:"mode" values:"dev,staging,prod"`
		}

		registry, cmd := setupFlagTest(t, TestConfig{})

		// Manually set an invalid value (bypassing validation)
		setFlag(t, cmd, "mode", "invalid")

		err := registry.ValidateFlags(cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "mode")
	})

	t.Run("unchanged flag skips validation", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Mode string `default:"dev" flag:"mode" values:"dev,staging,prod"`
		}

		registry, cmd := setupFlagTest(t, TestConfig{})

		// Don't change the flag - should pass validation
		err := registry.ValidateFlags(cmd)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("flag without values skips validation", func(t *testing.T) {
		t.Parallel()

		type testConfigFlagWithoutValues struct {
			Name string `default:"default" flag:"name"`
		}

		registry, cmd := setupFlagTest(t, testConfigFlagWithoutValues{})
		setFlagAndAssertValid(t, cmd, registry, "name", "anything")
	})

	t.Run("required flag not set returns error", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name string `flag:"name" help:"required name" required:"true"`
		}

		registry, cmd := setupFlagTest(t, TestConfig{})

		// Don't set the flag - should fail validation
		err := registry.ValidateFlags(cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "name", "required")
	})

	t.Run("required flag set passes validation", func(t *testing.T) {
		t.Parallel()

		type testConfigRequiredFlagSet struct {
			Name string `flag:"name" help:"required name" required:"true"`
		}

		registry, cmd := setupFlagTest(t, testConfigRequiredFlagSet{})
		setFlagAndAssertValid(t, cmd, registry, "name", "test-value")
	})

	t.Run("required false does not enforce", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name string `flag:"name" help:"optional name" required:"false"`
		}

		registry, cmd := setupFlagTest(t, TestConfig{})

		// Don't set the flag - should pass since required:"false"
		err := registry.ValidateFlags(cmd)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestParseValidateRules_UnknownValidator(t *testing.T) {
	t.Parallel()

	t.Run("unknown validator name returns error", func(t *testing.T) {
		t.Parallel()

		_, err := parseValidateRules("emial,min=5")
		if err == nil {
			t.Fatal("expected error for unknown validator, got nil")
		}

		assertErrorContains(t, err, "unknown validator", "emial")
	})

	t.Run("known validators parse successfully", func(t *testing.T) {
		t.Parallel()

		rules, err := parseValidateRules("email,minlen=5")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(rules) != 2 {
			t.Fatalf("expected 2 rules, got %d", len(rules))
		}
	})
}

func TestValidatorErrors_InvalidParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		validator func(string) error
		param     string
	}{
		{name: "minlen with non-integer param", validator: validateMinLen, param: "abc:hello"},
		{name: "maxlen with non-integer param", validator: validateMaxLen, param: "abc:hello"},
		{name: "min with non-number param", validator: validateMin, param: "abc:10"},
		{name: "max with non-number param", validator: validateMax, param: "abc:10"},
		{name: "regex with invalid pattern", validator: validateRegex, param: "[invalid:hello"},
		{name: "minlen missing separator", validator: validateMinLen, param: "5"},
		{name: "min missing separator", validator: validateMin, param: "5"},
		{name: "regex missing separator", validator: validateRegex, param: "[a-z]+"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertValidatorError(t, tt.name, tt.validator, tt.param)
		})
	}
}
