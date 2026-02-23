package v2

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagRegistry_ValidateFlags(t *testing.T) {
	t.Run("valid values pass", func(t *testing.T) {
		type TestConfig struct {
			Mode string `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("mode", "staging"))

		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})

	t.Run("invalid value returns error", func(t *testing.T) {
		type TestConfig struct {
			Mode string `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		// Manually set an invalid value (bypassing validation)
		require.NoError(t, cmd.PersistentFlags().Set("mode", "invalid"))

		err = registry.ValidateFlags(cmd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mode")
	})

	t.Run("unchanged flag skips validation", func(t *testing.T) {
		type TestConfig struct {
			Mode string `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		// Don't change the flag - should pass validation
		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})

	t.Run("flag without values skips validation", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" default:"default"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("name", "anything"))

		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})

	t.Run("required flag not set returns error", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" required:"true" help:"required name"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		// Don't set the flag - should fail validation
		err = registry.ValidateFlags(cmd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name")
		assert.Contains(t, err.Error(), "required")
	})

	t.Run("required flag set passes validation", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" required:"true" help:"required name"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("name", "test-value"))

		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})

	t.Run("required false does not enforce", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" required:"false" help:"optional name"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		// Don't set the flag - should pass since required:"false"
		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})
}
