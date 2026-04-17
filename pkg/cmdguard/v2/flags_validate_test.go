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

	t.Run("minlen with non-integer param", func(t *testing.T) {
		t.Parallel()

		err := validateMinLen("abc:hello")
		if err == nil {
			t.Fatal("expected error for non-integer minlen param, got nil")
		}

		assertErrorContains(t, err, "invalid")
	})

	t.Run("maxlen with non-integer param", func(t *testing.T) {
		t.Parallel()

		err := validateMaxLen("abc:hello")
		if err == nil {
			t.Fatal("expected error for non-integer maxlen param, got nil")
		}

		assertErrorContains(t, err, "invalid")
	})

	t.Run("min with non-number param", func(t *testing.T) {
		t.Parallel()

		err := validateMin("abc:10")
		if err == nil {
			t.Fatal("expected error for non-number min param, got nil")
		}

		assertErrorContains(t, err, "invalid")
	})

	t.Run("max with non-number param", func(t *testing.T) {
		t.Parallel()

		err := validateMax("abc:10")
		if err == nil {
			t.Fatal("expected error for non-number max param, got nil")
		}

		assertErrorContains(t, err, "invalid")
	})

	t.Run("regex with invalid pattern", func(t *testing.T) {
		t.Parallel()

		err := validateRegex("[invalid:hello")
		if err == nil {
			t.Fatal("expected error for invalid regex pattern, got nil")
		}

		assertErrorContains(t, err, "invalid")
	})

	t.Run("minlen missing separator", func(t *testing.T) {
		t.Parallel()

		err := validateMinLen("5")
		if err == nil {
			t.Fatal("expected error for minlen without separator, got nil")
		}

		assertErrorContains(t, err, "invalid")
	})
}
