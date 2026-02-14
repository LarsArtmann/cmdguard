package validation

import (
	"testing"

	"github.com/larsartmann/cmdguard/internal/config"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestValidator(t *testing.T) (*Validator, *Registry, do.Injector) {
	injector := do.New()

	// Provide config
	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return &config.Config{}, nil
	})

	// Provide registry
	do.Provide(injector, NewRegistry)

	// Provide validator
	do.Provide(injector, NewValidator)

	registry, err := do.Invoke[*Registry](injector)
	require.NoError(t, err)

	validator, err := do.Invoke[*Validator](injector)
	require.NoError(t, err)

	return validator, registry, injector
}

func TestNewValidator(t *testing.T) {
	validator, _, _ := setupTestValidator(t)
	assert.NotNil(t, validator)
	assert.NotNil(t, validator.registry)
	assert.NotNil(t, validator.cfg)
}

func TestValidator_ValidateCommands(t *testing.T) {
	validator, registry, _ := setupTestValidator(t)

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid command with handler",
			setup: func() {
				cmd := &cobra.Command{
					Use: "test",
					RunE: func(cmd *cobra.Command, args []string) error {
						return nil
					},
				}
				_ = registry.RegisterCommand(cmd)
			},
			wantErr: false,
		},
		{
			name: "command with subcommands is valid without handler",
			setup: func() {
				parent := &cobra.Command{Use: "parent"}
				child := &cobra.Command{
					Use: "child",
					RunE: func(cmd *cobra.Command, args []string) error {
						return nil
					},
				}
				parent.AddCommand(child)
				_ = registry.RegisterCommand(parent)
				_ = registry.RegisterSubcommand(parent, child)
			},
			wantErr: false,
		},
		{
			name: "command without handler fails",
			setup: func() {
				cmd := &cobra.Command{Use: "test"}
				_ = registry.RegisterCommand(cmd)
			},
			wantErr: true,
			errMsg:  "has no handler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := validator.ValidateCommands()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestValidator_ValidateFlags(t *testing.T) {
	validator, registry, _ := setupTestValidator(t)

	// Register command with unbound flag
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "default", "name flag")

	err := registry.RegisterCommand(cmd)
	require.NoError(t, err)

	// Flag is marked as bound by default, so validation passes
	err = validator.ValidateFlags()
	require.NoError(t, err)

	// Unmark the flag
	registry.UnmarkFlagBound("test", "name")

	// Now validation should fail
	err = validator.ValidateFlags()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not bound")
}

func TestValidator_ValidateCommandTree(t *testing.T) {
	validator, _, _ := setupTestValidator(t)

	// Create command tree
	root := &cobra.Command{
		Use: "root",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	child := &cobra.Command{
		Use: "child",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	root.AddCommand(child)

	// Validate entire tree
	err := validator.ValidateCommandTree(root)
	require.NoError(t, err)

	// Verify commands were registered
	assert.Equal(t, 2, validator.registry.CommandCount())
}

func TestValidator_IsStrictMode(t *testing.T) {
	tests := []struct {
		name     string
		strict   bool
		expected bool
	}{
		{
			name:     "strict mode disabled",
			strict:   false,
			expected: false,
		},
		{
			name:     "strict mode enabled",
			strict:   true,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injector := do.New()

			// Provide config with strict setting
			do.Provide(injector, func(i do.Injector) (*config.Config, error) {
				return &config.Config{StrictMode: tt.strict}, nil
			})

			// Provide registry
			do.Provide(injector, NewRegistry)

			// Provide validator
			do.Provide(injector, NewValidator)

			validator, err := do.Invoke[*Validator](injector)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, validator.IsStrictMode())
		})
	}
}

func TestNewFlagValidator(t *testing.T) {
	injector := do.New()

	// Provide config
	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return &config.Config{StrictMode: true}, nil
	})

	fv, err := NewFlagValidator(injector)
	require.NoError(t, err)
	assert.NotNil(t, fv)
	assert.True(t, fv.strict)
}

func TestFlagValidator_ValidateFlag(t *testing.T) {
	tests := []struct {
		name    string
		strict  bool
		value   any
		wantErr bool
	}{
		{
			name:    "non-strict mode allows nil",
			strict:  false,
			value:   nil,
			wantErr: false,
		},
		{
			name:    "strict mode rejects nil",
			strict:  true,
			value:   nil,
			wantErr: true,
		},
		{
			name:    "strict mode allows non-nil",
			strict:  true,
			value:   "test",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := &FlagValidator{strict: tt.strict}

			err := fv.ValidateFlag("test-flag", tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "is required in strict mode")
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestFlagValidator_ValidateFlagAccess(t *testing.T) {
	fv := &FlagValidator{}

	// Create command with flag
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "default", "name flag")

	// Valid flag access
	err := fv.ValidateFlagAccess(cmd, "name")
	require.NoError(t, err)

	// Invalid flag access
	err = fv.ValidateFlagAccess(cmd, "nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not registered")
}
