package v2

import (
	"strings"
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

	err := registry.ValidateFlags(cmd)
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

		if err := cmd.Flags().Set("mode", "staging"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

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
		if err := cmd.Flags().Set("mode", "invalid"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err := registry.ValidateFlags(cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "mode") {
			t.Errorf("expected error to contain 'mode', got: %v", err)
		}
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

	// flag without values skips validation
	{
		type TestConfig struct {
			Name string `default:"default" flag:"name"`
		}

		registry, cmd := setupFlagTest(t, TestConfig{})
		setFlagAndAssertValid(t, cmd, registry, "name", "anything")
	}

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

		if !strings.Contains(err.Error(), "name") {
			t.Errorf("expected error to contain 'name', got: %v", err)
		}

		if !strings.Contains(err.Error(), "required") {
			t.Errorf("expected error to contain 'required', got: %v", err)
		}
	})

	// required flag set passes validation
	{
		type TestConfig struct {
			Name string `flag:"name" help:"required name" required:"true"`
		}

		registry, cmd := setupFlagTest(t, TestConfig{})
		setFlagAndAssertValid(t, cmd, registry, "name", "test-value")
	}

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
