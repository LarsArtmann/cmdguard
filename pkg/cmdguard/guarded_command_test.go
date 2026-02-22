package cmdguard

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// hasSubcommand checks if a command with the given name exists in the GuardedCommand's subcommands.
func hasSubcommand(g *GuardedCommand, name string) bool {
	for _, cmd := range g.cmd.Commands() {
		if cmd.Name() == name {
			return true
		}
	}
	return false
}

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

		assert.True(t, hasSubcommand(g, "version"), "version command should be added by default")
	})

	t.Run("validate command is added", func(t *testing.T) {
		g := New("testapp", "Test")

		assert.True(t, hasSubcommand(g, "validate"), "validate command should be added by default")
	})
}

func TestVersion(t *testing.T) {
	t.Run("returns version string", func(t *testing.T) {
		v := Version()
		assert.NotEmpty(t, v)
		assert.Equal(t, "dev", v) // Default value without ldflags
	})
}

func TestGuardedCommand_ValidateCommandTree(t *testing.T) {
	t.Run("validates command tree successfully", func(t *testing.T) {
		g := New("testapp", "Test")

		// Add a valid command with subcommand
		parent := &cobra.Command{Use: "parent"}
		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		parent.AddCommand(child)
		g.AddCommand(parent)

		err := g.validateCommandTree()
		require.NoError(t, err)
	})

	t.Run("returns error for invalid nested command", func(t *testing.T) {
		g := New("testapp", "Test")

		// Create a command tree with an invalid nested subcommand
		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		// Manually add an invalid child (bypassing guard validation)
		invalidChild := &cobra.Command{Use: "invalid"} // No handler
		parent.AddCommand(invalidChild)
		g.cmd.AddCommand(parent)

		err := g.validateCommandTree()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid")
	})
}

func TestGuardedCommand_ValidateSubcommands(t *testing.T) {
	t.Run("validates valid subcommands", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		parent.AddCommand(child)

		err := g.validateSubcommands(parent)
		require.NoError(t, err)
	})

	t.Run("returns error for invalid subcommand", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		// Manually add invalid child (bypassing guard)
		invalidChild := &cobra.Command{Use: "invalid"} // No handler
		parent.AddCommand(invalidChild)

		err := g.validateSubcommands(parent)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "parent invalid")
	})

	t.Run("validates nested subcommands", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		grandchild := &cobra.Command{
			Use: "grandchild",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		child.AddCommand(grandchild)
		parent.AddCommand(child)

		err := g.validateSubcommands(parent)
		require.NoError(t, err)
	})
}

func TestGuardedCommand_DefaultCommands_Execution(t *testing.T) {
	t.Run("version command executes", func(t *testing.T) {
		g := New("testapp", "Test")

		// Find version command
		var versionCmd *cobra.Command
		for _, cmd := range g.cmd.Commands() {
			if cmd.Name() == "version" {
				versionCmd = cmd
				break
			}
		}
		require.NotNil(t, versionCmd)

		// Execute it
		versionCmd.Run(versionCmd, []string{})
		// If we get here without panic, it worked
	})

	t.Run("validate command succeeds with valid tree", func(t *testing.T) {
		g := New("testapp", "Test")

		// Add a valid command
		cmd := &cobra.Command{
			Use: "testcmd",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		g.AddCommand(cmd)

		// Find validate command
		var validateCmd *cobra.Command
		for _, c := range g.cmd.Commands() {
			if c.Name() == "validate" {
				validateCmd = c
				break
			}
		}
		require.NotNil(t, validateCmd)

		// Execute it
		err := validateCmd.RunE(validateCmd, []string{})
		require.NoError(t, err)
	})

	t.Run("validate command fails with invalid tree", func(t *testing.T) {
		g := New("testapp", "Test")

		// Manually add invalid command (bypassing guard)
		invalidCmd := &cobra.Command{Use: "invalid"} // No handler
		g.cmd.AddCommand(invalidCmd)

		// Find validate command
		var validateCmd *cobra.Command
		for _, c := range g.cmd.Commands() {
			if c.Name() == "validate" {
				validateCmd = c
				break
			}
		}
		require.NotNil(t, validateCmd)

		// Execute it - should error
		err := validateCmd.RunE(validateCmd, []string{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})
}
