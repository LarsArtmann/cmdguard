package v2

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagRegistry_ParseFlags(t *testing.T) {
	t.Run("parse string flag", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" default:"default"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		// Set flag value
		require.NoError(t, cmd.PersistentFlags().Set("name", "custom"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		assert.Equal(t, "custom", cfg.Name)
	})

	t.Run("parse bool flag", func(t *testing.T) {
		type TestConfig struct {
			Verbose bool `flag:"verbose" default:"false"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("verbose", "true"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		assert.True(t, cfg.Verbose)
	})

	t.Run("parse int flag", func(t *testing.T) {
		type TestConfig struct {
			Count int `flag:"count" default:"0"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("count", "42"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		assert.Equal(t, 42, cfg.Count)
	})

	t.Run("parse float64 flag", func(t *testing.T) {
		type TestConfig struct {
			Rate float64 `flag:"rate" default:"0.0"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("rate", "3.14159"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		assert.InDelta(t, 3.14159, cfg.Rate, 0.00001)
	})

	t.Run("parse Duration flag", func(t *testing.T) {
		type TestConfig struct {
			Timeout Duration `flag:"timeout" default:"1m"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("timeout", "5m30s"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		expected := FromDuration(5*time.Minute + 30*time.Second)
		assert.Equal(t, expected, cfg.Timeout)
	})

	t.Run("parse LogLevel flag", func(t *testing.T) {
		type TestConfig struct {
			Level LogLevel `flag:"level" default:"info"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("level", "debug"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		assert.Equal(t, LogLevelDebug, cfg.Level)
	})

	t.Run("parse invalid LogLevel returns error", func(t *testing.T) {
		type TestConfig struct {
			Level LogLevel `flag:"level" default:"info"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("level", "invalid"))

		err = registry.ParseFlags(cmd, cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "level")
	})

	t.Run("parse Enum flag", func(t *testing.T) {
		type TestConfig struct {
			Mode Enum `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("mode", "prod"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		assert.Equal(t, "prod", cfg.Mode.String())
	})

	t.Run("parse invalid Enum returns error", func(t *testing.T) {
		type TestConfig struct {
			Mode Enum `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("mode", "invalid"))

		err = registry.ParseFlags(cmd, cfg)
		require.Error(t, err)
	})
}

func TestFlagRegistry_FlagNotFound(t *testing.T) {
	t.Run("missing flag returns error", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		// Don't register flags on command
		cmd := &cobra.Command{Use: "test"}

		err = registry.ParseFlags(cmd, cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
