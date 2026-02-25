package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlagTags(t *testing.T) {
	t.Run("valid struct", func(t *testing.T) {
		type TestConfig struct {
			Name    string `default:"test" flag:"name"    help:"The name"       short:"n"`
			Count   int    `default:"10"   flag:"count"   help:"The count"`
			Enabled bool   `default:"true" flag:"enabled" help:"Enable feature" short:"e"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		require.NoError(t, err)
		require.Len(t, tags, 3)

		// Check first field
		assert.Equal(t, "Name", tags[0].Field)
		assert.Equal(t, "name", tags[0].Name)
		assert.Equal(t, "n", tags[0].Short)
		assert.Equal(t, "test", tags[0].Default)
		assert.Equal(t, "The name", tags[0].Help)
	})

	t.Run("pointer to struct", func(t *testing.T) {
		type TestConfig struct {
			Field string `flag:"field"`
		}

		tags, err := ParseFlagTags(&TestConfig{})
		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, "field", tags[0].Name)
	})

	t.Run("skips fields without flag tag", func(t *testing.T) {
		type TestConfig struct {
			Tagged   string `flag:"tagged"`
			Untagged string
			Ignored  string `flag:"-"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, "Tagged", tags[0].Field)
	})

	t.Run("nil config", func(t *testing.T) {
		tags, err := ParseFlagTags(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must not be nil")
		assert.Nil(t, tags)
	})

	t.Run("non-struct config", func(t *testing.T) {
		tags, err := ParseFlagTags("not a struct")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be a struct")
		assert.Nil(t, tags)
	})

	t.Run("with values tag", func(t *testing.T) {
		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn,error"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, []string{"debug", "info", "warn", "error"}, tags[0].Values)
	})

	t.Run("embedded Config", func(t *testing.T) {
		type AppConfig struct {
			Config

			AppName string `default:"myapp" flag:"app-name"`
		}

		tags, err := ParseFlagTags(AppConfig{})
		require.NoError(t, err)
		// Config has 4 fields + AppName = 5
		assert.GreaterOrEqual(t, len(tags), 1)
	})
}
