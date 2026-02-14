package cmdguard

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates GuardedCommand with defaults", func(t *testing.T) {
		_ = os.Unsetenv("CMDGUARD_LOG_LEVEL")
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")

		g := New("testapp", "Test application")

		require.NotNil(t, g)
		assert.NotNil(t, g.cmd)
		assert.NotNil(t, g.cfg)
		assert.Equal(t, "testapp", g.cmd.Use)
		assert.Equal(t, "Test application", g.cmd.Short)
	})

	t.Run("loads config from environment", func(t *testing.T) {
		_ = os.Setenv("CMDGUARD_LOG_LEVEL", "debug")
		_ = os.Setenv("CMDGUARD_STRICT_MODE", "true")
		defer func() {
			_ = os.Unsetenv("CMDGUARD_LOG_LEVEL")
			_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
		}()

		g := New("testapp", "Test")

		require.NotNil(t, g)
		assert.Equal(t, "debug", g.cfg.LogLevel)
		assert.True(t, g.cfg.StrictMode)
		assert.True(t, g.strictMode)
	})
}

func TestGuardedCommand_AddCommand(t *testing.T) {
	t.Run("accepts valid command with Run", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		assert.NotPanics(t, func() {
			g.AddCommand(cmd)
		})
	})

	t.Run("accepts valid command with RunE", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		assert.NotPanics(t, func() {
			g.AddCommand(cmd)
		})
	})

	t.Run("panics on command without handler", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "invalid",
		}

		assert.Panics(t, func() {
			g.AddCommand(cmd)
		})
	})

	t.Run("panics on command without name", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{}

		assert.Panics(t, func() {
			g.AddCommand(cmd)
		})
	})

	t.Run("panics after Execute called", func(t *testing.T) {
		g := New("testapp", "Test")
		g.validated = true // Simulate post-execute state

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		assert.Panics(t, func() {
			g.AddCommand(cmd)
		})
	})
}

func TestGuardedCommand_AddSubcommand(t *testing.T) {
	t.Run("adds subcommand to parent", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		g.AddCommand(parent)

		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		assert.NotPanics(t, func() {
			g.AddSubcommand(parent, child)
		})

		found := false
		for _, c := range parent.Commands() {
			if c.Name() == "child" {
				found = true
				break
			}
		}
		assert.True(t, found, "child command should be added to parent")
	})

	t.Run("panics on invalid child", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		g.AddCommand(parent)

		child := &cobra.Command{
			Use: "invalid-child",
			// No Run or RunE
		}

		assert.Panics(t, func() {
			g.AddSubcommand(parent, child)
		})
	})
}

func TestGuardedCommand_Execute(t *testing.T) {
	t.Run("executes command successfully", func(t *testing.T) {
		g := New("testapp", "Test")
		g.cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return nil
		}

		err := g.Execute(context.Background())
		require.NoError(t, err)
		assert.True(t, g.validated)
	})
}

func TestGuardedCommand_Accessors(t *testing.T) {
	t.Run("Command returns underlying cobra command", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := g.Command()

		require.NotNil(t, cmd)
		assert.Equal(t, "testapp", cmd.Use)
	})

	t.Run("Config returns config", func(t *testing.T) {
		g := New("testapp", "Test")

		cfg := g.Config()

		require.NotNil(t, cfg)
		assert.Equal(t, "info", cfg.LogLevel)
	})

	t.Run("IsStrictMode returns correct value", func(t *testing.T) {
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
		g := New("testapp", "Test")
		assert.False(t, g.IsStrictMode())

		_ = os.Setenv("CMDGUARD_STRICT_MODE", "true")
		defer func() { _ = os.Unsetenv("CMDGUARD_STRICT_MODE") }()

		g2 := New("testapp2", "Test2")
		assert.True(t, g2.IsStrictMode())
	})
}

func TestGuardedCommand_validateCommand(t *testing.T) {
	t.Run("command with Run is valid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		err := g.validateCommand(cmd)
		require.NoError(t, err)
	})

	t.Run("command with RunE is valid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		err := g.validateCommand(cmd)
		require.NoError(t, err)
	})

	t.Run("command without name is invalid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{}

		err := g.validateCommand(cmd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no name")
	})

	t.Run("command without handler is invalid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "invalid",
		}

		err := g.validateCommand(cmd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no handler")
	})

	t.Run("parent command with subcommands does not need handler", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{Use: "parent"}
		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		parent.AddCommand(child)

		err := g.validateCommand(parent)
		require.NoError(t, err)
	})

	t.Run("strict mode requires RunE", func(t *testing.T) {
		_ = os.Setenv("CMDGUARD_STRICT_MODE", "true")
		defer func() { _ = os.Unsetenv("CMDGUARD_STRICT_MODE") }()

		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		err := g.validateCommand(cmd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "strict mode requires RunE")
	})
}

func TestGuardedCommand_DefaultCommands(t *testing.T) {
	t.Run("version command is added", func(t *testing.T) {
		g := New("testapp", "Test")

		found := false
		for _, cmd := range g.cmd.Commands() {
			if cmd.Name() == "version" {
				found = true
				break
			}
		}
		assert.True(t, found, "version command should be added by default")
	})

	t.Run("validate command is added", func(t *testing.T) {
		g := New("testapp", "Test")

		found := false
		for _, cmd := range g.cmd.Commands() {
			if cmd.Name() == "validate" {
				found = true
				break
			}
		}
		assert.True(t, found, "validate command should be added by default")
	})
}
