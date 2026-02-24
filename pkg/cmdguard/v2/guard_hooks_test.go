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
	tests := []struct {
		name     string
		hookName string
		setupCmd func(order *[]string) Command[TestAppConfig, NoFlags]
		want     []string
	}{
		{
			name:     "calls PreRunE before RunE",
			hookName: "pre",
			setupCmd: func(order *[]string) Command[TestAppConfig, NoFlags] {
				return Command[TestAppConfig, NoFlags]{
					Use: "test",
					PreRunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
						*order = append(*order, "pre")
						return nil
					},
					RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
						*order = append(*order, "run")
						return nil
					},
				}
			},
			want: []string{"pre", "run"},
		},
		{
			name:     "calls PostRunE after RunE",
			hookName: "post",
			setupCmd: func(order *[]string) Command[TestAppConfig, NoFlags] {
				return Command[TestAppConfig, NoFlags]{
					Use: "test",
					RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
						*order = append(*order, "run")
						return nil
					},
					PostRunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
						*order = append(*order, "post")
						return nil
					},
				}
			},
			want: []string{"run", "post"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var order []string

			g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
			require.NoError(t, err)

			cmd := tt.setupCmd(&order)
			require.NoError(t, g.AddCommand(cmd))

			err = g.ExecuteWithArgs(context.Background(), []string{"test"})
			require.NoError(t, err)
			assert.Equal(t, tt.want, order)
		})
	}

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
	tests := []struct {
		name           string
		use            string
		hidden         bool
		deprecated     string
		aliases        []string
		version        string
		wantHidden     bool
		wantDeprecated string
		wantAliases    []string
		wantVersion    string
	}{
		{
			name:       "hidden command",
			use:        "secret",
			hidden:     true,
			wantHidden: true,
		},
		{
			name:           "deprecated command",
			use:            "old",
			deprecated:     "use new-cmd instead",
			wantDeprecated: "use new-cmd instead",
		},
		{
			name:        "command with aliases",
			use:         "list",
			aliases:     []string{"ls", "l"},
			wantAliases: []string{"ls", "l"},
		},
		{
			name:        "command with version",
			use:         "versioned",
			version:     "v1.2.3",
			wantVersion: "v1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
			require.NoError(t, err)

			cmd := Command[TestAppConfig, NoFlags]{
				Use:        tt.use,
				Hidden:     tt.hidden,
				Deprecated: tt.deprecated,
				Aliases:    tt.aliases,
				Version:    tt.version,
				RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
					return nil
				},
			}
			require.NoError(t, g.AddCommand(cmd))

			cobraCmd := g.RootCommand().Commands()[0]
			assert.Equal(t, tt.wantHidden, cobraCmd.Hidden)
			assert.Equal(t, tt.wantDeprecated, cobraCmd.Deprecated)
			assert.Equal(t, tt.wantAliases, cobraCmd.Aliases)
			assert.Equal(t, tt.wantVersion, cobraCmd.Version)
		})
	}
}
