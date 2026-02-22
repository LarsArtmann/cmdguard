package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuardedCommand_AddCommand(t *testing.T) {
	t.Run("adds valid command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "greet",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}

		err = g.AddCommand(cmd)
		require.NoError(t, err)

		rootCmd := g.RootCommand()
		require.Len(t, rootCmd.Commands(), 1)
	})

	t.Run("error: invalid command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "greet",
			// No RunE
		}

		err = g.AddCommand(cmd)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingHandler)
	})

	t.Run("adds command with subcommands", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		subCmd := Command[TestAppConfig, NoFlags]{
			Use: "list",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}

		cmd := Command[TestAppConfig, NoFlags]{
			Use:      "greet",
			Commands: []Command[TestAppConfig, NoFlags]{subCmd},
		}

		err = g.AddCommand(cmd)
		require.NoError(t, err)

		rootCmd := g.RootCommand()
		require.Len(t, rootCmd.Commands(), 1)
		greetCmd := rootCmd.Commands()[0]
		require.Len(t, greetCmd.Commands(), 1)
	})

	t.Run("error: duplicate command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd1 := Command[TestAppConfig, NoFlags]{
			Use: "greet",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}

		cmd2 := Command[TestAppConfig, NoFlags]{
			Use: "greet",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}

		err = g.AddCommand(cmd1)
		require.NoError(t, err)

		err = g.AddCommand(cmd2)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateCommand)
		assert.Contains(t, err.Error(), "greet")
	})
}

func TestGuardedCommand_AddCommandFunc(t *testing.T) {
	t.Run("adds command via function", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		err = g.AddCommandFunc(func() Command[TestAppConfig, NoFlags] {
			return Command[TestAppConfig, NoFlags]{
				Use: "greet",
				RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
					return nil
				},
			}
		})
		require.NoError(t, err)

		rootCmd := g.RootCommand()
		require.Len(t, rootCmd.Commands(), 1)
	})
}
