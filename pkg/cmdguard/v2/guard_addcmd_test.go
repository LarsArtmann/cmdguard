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

func TestAddAnyCommand(t *testing.T) {
	t.Run("adds command with different flag type", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		type GreetFlags struct {
			Name string
		}

		cmd := Command[TestAppConfig, *GreetFlags]{
			Use:   "greet",
			Short: "Greet someone",
			Flags: &GreetFlags{},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags *GreetFlags) error {
				return nil
			},
		}

		err = AddAnyCommand(g, cmd)
		require.NoError(t, err)

		rootCmd := g.RootCommand()
		require.Len(t, rootCmd.Commands(), 1)
	})

	t.Run("error when command has no handler", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		type OtherFlags struct {
			Value string
		}

		cmd := Command[TestAppConfig, *OtherFlags]{
			Use:   "invalid",
			Flags: &OtherFlags{},
			// No RunE
		}

		err = AddAnyCommand(g, cmd)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingHandler)
	})

	t.Run("error on duplicate command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		type OtherFlags struct{}

		cmd1 := Command[TestAppConfig, *OtherFlags]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags *OtherFlags) error {
				return nil
			},
		}

		cmd2 := Command[TestAppConfig, *OtherFlags]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags *OtherFlags) error {
				return nil
			},
		}

		err = AddAnyCommand(g, cmd1)
		require.NoError(t, err)

		err = AddAnyCommand(g, cmd2)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateCommand)
	})
}
