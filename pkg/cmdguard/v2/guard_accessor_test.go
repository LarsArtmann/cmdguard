package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuardedCommand_Scope(t *testing.T) {
	t.Run("returns injector", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		injector := g.Scope()
		require.NotNil(t, injector)
	})

	t.Run("returns scope struct", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		scope := g.ScopeStruct()
		require.NotNil(t, scope)
		assert.Equal(t, "myapp", scope.Name())
	})
}

func TestGuardedCommand_Config(t *testing.T) {
	t.Run("returns config", func(t *testing.T) {
		defaults := TestAppConfig{Verbose: true, Output: "/tmp/out"}
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", defaults)
		require.NoError(t, err)

		cfg := g.Config()
		require.NotNil(t, cfg)
		assert.True(t, cfg.Verbose)
		assert.Equal(t, "/tmp/out", cfg.Output)
	})

	t.Run("SetConfig updates config", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		newCfg := TestAppConfig{Verbose: true, Output: "/new/path"}
		g.SetConfig(newCfg)

		cfg := g.Config()
		assert.True(t, cfg.Verbose)
		assert.Equal(t, "/new/path", cfg.Output)
	})
}

func TestGuardedCommand_RootCommand(t *testing.T) {
	t.Run("returns cobra command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		rootCmd := g.RootCommand()
		require.NotNil(t, rootCmd)
		assert.Equal(t, "myapp", rootCmd.Use)
		assert.Equal(t, "My CLI", rootCmd.Short)
	})
}

func TestGuardedCommand_Metadata(t *testing.T) {
	t.Run("Name returns name", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)
		assert.Equal(t, "myapp", g.Name())
	})

	t.Run("Short returns short description", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)
		assert.Equal(t, "My CLI", g.Short())
	})

	t.Run("Long returns long description", func(t *testing.T) {
		g, err := NewWithLong[TestAppConfig, NoFlags](
			"myapp",
			"short",
			"long desc",
			TestAppConfig{},
		)
		require.NoError(t, err)
		assert.Equal(t, "long desc", g.Long())
	})

	t.Run("SetLong updates long description", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		g.SetLong("new long description")
		assert.Equal(t, "new long description", g.Long())
		assert.Equal(t, "new long description", g.RootCommand().Long)
	})

	t.Run("SetVersion sets version", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		g.SetVersion("v1.0.0")
		assert.Equal(t, "v1.0.0", g.RootCommand().Version)
	})
}

func TestGuardedCommand_AddGlobalFlag(t *testing.T) {
	t.Run("adds global string flag", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		g.AddGlobalFlag("config", "c", "/etc/config.yaml", "Config file path")

		flag := g.RootCommand().PersistentFlags().Lookup("config")
		require.NotNil(t, flag)
		assert.Equal(t, "c", flag.Shorthand)
		assert.Equal(t, "/etc/config.yaml", flag.DefValue)
	})

	t.Run("adds global bool flag", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		g.AddGlobalBoolFlag("debug", "d", true, "Enable debug mode")

		flag := g.RootCommand().PersistentFlags().Lookup("debug")
		require.NotNil(t, flag)
		assert.Equal(t, "d", flag.Shorthand)
		assert.Equal(t, "true", flag.DefValue)
	})
}
