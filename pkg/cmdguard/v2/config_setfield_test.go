package v2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetField(t *testing.T) {
	t.Run("set string field", func(t *testing.T) {
		cfg := &struct {
			Name string
		}{}
		err := SetField(cfg, "Name", "test")
		require.NoError(t, err)
		assert.Equal(t, "test", cfg.Name)
	})

	t.Run("set int field", func(t *testing.T) {
		cfg := &struct {
			Count int
		}{}
		err := SetField(cfg, "Count", 42)
		require.NoError(t, err)
		assert.Equal(t, 42, cfg.Count)
	})

	t.Run("set bool field", func(t *testing.T) {
		cfg := &struct {
			Enabled bool
		}{}
		err := SetField(cfg, "Enabled", true)
		require.NoError(t, err)
		assert.True(t, cfg.Enabled)
	})

	t.Run("non-pointer config", func(t *testing.T) {
		cfg := struct{ Name string }{}
		err := SetField(cfg, "Name", "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pointer to struct")
	})

	t.Run("field not found", func(t *testing.T) {
		cfg := &struct{ Name string }{}
		err := SetField(cfg, "Missing", "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("string to LogLevel", func(t *testing.T) {
		cfg := &struct {
			Level LogLevel
		}{}
		err := SetField(cfg, "Level", "debug")
		require.NoError(t, err)
		assert.Equal(t, LogLevelDebug, cfg.Level)
	})

	t.Run("string to LogFormat", func(t *testing.T) {
		cfg := &struct {
			Format LogFormat
		}{}
		err := SetField(cfg, "Format", "json")
		require.NoError(t, err)
		assert.Equal(t, LogFormatJSON, cfg.Format)
	})

	t.Run("string to Duration", func(t *testing.T) {
		cfg := &struct {
			Timeout Duration
		}{}
		err := SetField(cfg, "Timeout", "30s")
		require.NoError(t, err)
		assert.Equal(t, 30*time.Second, cfg.Timeout.Duration())
	})

	t.Run("time.Duration to Duration", func(t *testing.T) {
		cfg := &struct {
			Timeout Duration
		}{}
		err := SetField(cfg, "Timeout", 45*time.Second)
		require.NoError(t, err)
		assert.Equal(t, 45*time.Second, cfg.Timeout.Duration())
	})

	t.Run("invalid LogLevel", func(t *testing.T) {
		cfg := &struct {
			Level LogLevel
		}{}
		err := SetField(cfg, "Level", "invalid")
		require.Error(t, err)
	})

	t.Run("incompatible types", func(t *testing.T) {
		cfg := &struct {
			Name string
		}{}
		// Use a struct type which is truly incompatible with string
		err := SetField(cfg, "Name", Duration{duration: 5 * time.Minute})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot convert")
	})
}
