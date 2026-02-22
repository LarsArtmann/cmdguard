package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

func TestGuardedCommand_PreRunE_PostRunE(t *testing.T) {
	t.Run("calls PreRunE before RunE", func(t *testing.T) {
		var order []string

		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "test",
			PreRunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				order = append(order, "pre")
				return nil
			},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				order = append(order, "run")
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"test"})
		require.NoError(t, err)
		assert.Equal(t, []string{"pre", "run"}, order)
	})

	t.Run("calls PostRunE after RunE", func(t *testing.T) {
		var order []string

		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				order = append(order, "run")
				return nil
			},
			PostRunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				order = append(order, "post")
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"test"})
		require.NoError(t, err)
		assert.Equal(t, []string{"run", "post"}, order)
	})

	t.Run("PreRunE error stops execution", func(t *testing.T) {
		called := false

		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "test",
			PreRunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return errTest
			},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				called = true
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"test"})
		require.Error(t, err)
		assert.False(t, called)
	})
}

func TestGuardedCommand_CommandOptions(t *testing.T) {
	t.Run("hidden command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use:    "secret",
			Hidden: true,
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		cobraCmd := g.RootCommand().Commands()[0]
		assert.True(t, cobraCmd.Hidden)
	})

	t.Run("deprecated command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use:        "old",
			Deprecated: "use new-cmd instead",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		cobraCmd := g.RootCommand().Commands()[0]
		assert.Equal(t, "use new-cmd instead", cobraCmd.Deprecated)
	})

	t.Run("command with aliases", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use:     "list",
			Aliases: []string{"ls", "l"},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		cobraCmd := g.RootCommand().Commands()[0]
		assert.Equal(t, []string{"ls", "l"}, cobraCmd.Aliases)
	})

	t.Run("command with version", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use:     "versioned",
			Version: "v1.2.3",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		cobraCmd := g.RootCommand().Commands()[0]
		assert.Equal(t, "v1.2.3", cobraCmd.Version)
	})
}
