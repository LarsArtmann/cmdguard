package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestAppConfig struct {
	Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
	Output  string `flag:"output" short:"o" default:"-" help:"Output file"`
}

func TestVersion(t *testing.T) {
	t.Run("version is set", func(t *testing.T) {
		assert.NotEmpty(t, Version)
		assert.Equal(t, "2.0.0", Version)
	})
}

func TestNew(t *testing.T) {
	t.Run("creates GuardedCommand", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI application", defaults)
		require.NoError(t, err)
		require.NotNil(t, g)

		assert.Equal(t, "myapp", g.Name())
		assert.Equal(t, "My CLI application", g.Short())
		assert.Equal(t, "", g.Long())
	})

	t.Run("error: empty name", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := New[TestAppConfig, NoFlags]("", "My CLI", defaults)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidCommand)
		assert.Nil(t, g)
	})

	t.Run("registers config in scope", func(t *testing.T) {
		defaults := TestAppConfig{Verbose: true}
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", defaults)
		require.NoError(t, err)

		cfg := g.Config()
		require.NotNil(t, cfg)
		assert.True(t, cfg.Verbose)
	})

	t.Run("creates scope", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", defaults)
		require.NoError(t, err)

		scope := g.ScopeStruct()
		require.NotNil(t, scope)
		assert.Equal(t, "myapp", scope.Name())
	})
}

func TestNewWithLong(t *testing.T) {
	t.Run("creates GuardedCommand with long description", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := NewWithLong[TestAppConfig, NoFlags]("myapp", "short", "long description", defaults)
		require.NoError(t, err)
		require.NotNil(t, g)

		assert.Equal(t, "myapp", g.Name())
		assert.Equal(t, "short", g.Short())
		assert.Equal(t, "long description", g.Long())
	})

	t.Run("error: empty name", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := NewWithLong[TestAppConfig, NoFlags]("", "short", "long", defaults)
		require.Error(t, err)
		assert.Nil(t, g)
	})
}

func TestNew_FlagTypeValidation(t *testing.T) {
	t.Run("rejects invalid flag type in New", func(t *testing.T) {
		g, err := New[TestAppConfig, int]("myapp", "My CLI", TestAppConfig{})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Nil(t, g)
	})

	t.Run("accepts NoFlags in New", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)
		require.NotNil(t, g)
	})

	t.Run("accepts pointer to struct in New", func(t *testing.T) {
		type CmdFlags struct {
			Name string `flag:"name"`
		}
		g, err := New[TestAppConfig, *CmdFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)
		require.NotNil(t, g)
	})
}
