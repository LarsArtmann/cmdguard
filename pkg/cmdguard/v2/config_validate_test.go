package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		type TestConfig struct {
			Name  string `flag:"name"`
			Count int    `flag:"count"`
		}

		err := ValidateConfig(TestConfig{Name: "test", Count: 10})
		require.NoError(t, err)
	})

	t.Run("nil config", func(t *testing.T) {
		err := ValidateConfig(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must not be nil")
	})

	t.Run("non-struct config", func(t *testing.T) {
		err := ValidateConfig("not a struct")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected struct")
	})

	t.Run("valid enum value", func(t *testing.T) {
		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn"`
		}

		err := ValidateConfig(TestConfig{Level: "info"})
		require.NoError(t, err)
	})

	t.Run("invalid enum value", func(t *testing.T) {
		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn"`
		}

		err := ValidateConfig(TestConfig{Level: "invalid"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config validation")
	})

	t.Run("pointer to config", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name"`
		}

		err := ValidateConfig(&TestConfig{Name: "test"})
		require.NoError(t, err)
	})

	t.Run("LogLevel field with values", func(t *testing.T) {
		type TestConfig struct {
			Level LogLevel `flag:"level" values:"debug,info,warn,error"`
		}

		cfg := TestConfig{Level: LogLevelInfo}
		err := ValidateConfig(cfg)
		require.NoError(t, err)
	})

	t.Run("LogFormat field with values", func(t *testing.T) {
		type TestConfig struct {
			Format LogFormat `flag:"format" values:"text,json"`
		}

		cfg := TestConfig{Format: LogFormatText}
		err := ValidateConfig(cfg)
		require.NoError(t, err)
	})
}
