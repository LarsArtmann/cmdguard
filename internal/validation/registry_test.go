package validation

import (
	"testing"

	"github.com/larsartmann/cmdguard/internal/config"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRegistry(t *testing.T) (*Registry, do.Injector) {
	injector := do.New()

	// Provide config
	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return &config.Config{}, nil
	})

	registry, err := NewRegistry(injector)
	require.NoError(t, err)

	return registry, injector
}

func TestNewRegistry(t *testing.T) {
	injector := do.New()

	// Provide config
	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return &config.Config{}, nil
	})

	registry, err := NewRegistry(injector)
	require.NoError(t, err)
	assert.NotNil(t, registry)
	assert.NotNil(t, registry.commands)
	assert.NotNil(t, registry.flags)
}

func TestRegistry_RegisterCommand(t *testing.T) {
	registry, _ := setupTestRegistry(t)

	tests := []struct {
		name    string
		cmd     *cobra.Command
		wantErr bool
	}{
		{
			name: "register command with RunE handler",
			cmd: &cobra.Command{
				Use: "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "register command with Run handler",
			cmd: &cobra.Command{
				Use: "test2",
				Run: func(cmd *cobra.Command, args []string) {
				},
			},
			wantErr: false,
		},
		{
			name: "register command without handler",
			cmd: &cobra.Command{
				Use: "test3",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.RegisterCommand(tt.cmd)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify command was registered
			info, ok := registry.GetCommand(tt.cmd.Name())
			require.True(t, ok)
			assert.Equal(t, tt.cmd.Name(), info.Name)
		})
	}
}

func TestRegistry_GetCommand(t *testing.T) {
	registry, _ := setupTestRegistry(t)

	// Register a command
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	err := registry.RegisterCommand(cmd)
	require.NoError(t, err)

	// Test getting existing command
	info, ok := registry.GetCommand("test")
	require.True(t, ok)
	assert.Equal(t, "test", info.Name)

	// Test getting non-existent command
	_, ok = registry.GetCommand("nonexistent")
	assert.False(t, ok)
}

func TestRegistry_CommandCount(t *testing.T) {
	registry, _ := setupTestRegistry(t)

	// Initially empty
	assert.Equal(t, 0, registry.CommandCount())

	// Register commands
	cmd1 := &cobra.Command{Use: "cmd1"}
	cmd2 := &cobra.Command{Use: "cmd2"}

	err := registry.RegisterCommand(cmd1)
	require.NoError(t, err)
	assert.Equal(t, 1, registry.CommandCount())

	err = registry.RegisterCommand(cmd2)
	require.NoError(t, err)
	assert.Equal(t, 2, registry.CommandCount())
}

func TestRegistry_FlagCount(t *testing.T) {
	registry, _ := setupTestRegistry(t)

	// Initially empty
	assert.Equal(t, 0, registry.FlagCount())

	// Register command with flags
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "default", "name flag")
	cmd.Flags().Int("count", 1, "count flag")

	err := registry.RegisterCommand(cmd)
	require.NoError(t, err)

	// Should have 2 flags
	assert.Equal(t, 2, registry.FlagCount())
}

func TestRegistry_HealthCheck(t *testing.T) {
	registry, _ := setupTestRegistry(t)

	// Empty registry should pass health check
	err := registry.HealthCheck()
	require.NoError(t, err)

	// Register command without handler (no subcommands)
	cmd := &cobra.Command{
		Use: "test",
	}
	err = registry.RegisterCommand(cmd)
	require.NoError(t, err)

	// Health check should fail - command has no handler
	err = registry.HealthCheck()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "has no handler")
}

func TestRegistry_MarkFlagBound(t *testing.T) {
	registry, _ := setupTestRegistry(t)

	// Register command with flag
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "default", "name flag")

	err := registry.RegisterCommand(cmd)
	require.NoError(t, err)

	// Initially flag should be bound (default)
	flags, _ := registry.GetFlags("test")
	require.Len(t, flags, 1)
	assert.True(t, flags[0].IsBound)

	// Unmark flag
	registry.UnmarkFlagBound("test", "name")
	flags, _ = registry.GetFlags("test")
	assert.False(t, flags[0].IsBound)

	// Mark flag again
	registry.MarkFlagBound("test", "name")
	flags, _ = registry.GetFlags("test")
	assert.True(t, flags[0].IsBound)
}
