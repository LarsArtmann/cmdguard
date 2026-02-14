// Package commands provides Cobra command tree setup for cmdguard.
package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/larsartmann/cmdguard/internal/config"
	"github.com/larsartmann/cmdguard/internal/logging"
	"github.com/larsartmann/cmdguard/internal/validation"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// Registry manages the Cobra command tree.
type Registry struct {
	root      *cobra.Command
	cfg       *config.Config
	validator *validation.Validator
	logger    *slog.Logger
}

// NewRegistry creates a new command registry with the root command.
func NewRegistry(i do.Injector) (*Registry, error) {
	cfg, err := do.Invoke[*config.Config](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke config: %w", err)
	}

	root := &cobra.Command{
		Use:   "cmdguard",
		Short: "CLI with validated commands",
		Long: `cmdguard is a CLI validation library that ensures every command and flag
is properly implemented and bound. It combines fang (styling), koanf (config),
and samber/do/v2 (DI) with compile-time and runtime validation.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add global flags with short versions and clear defaults
	root.PersistentFlags().StringP("config", "c", "", "Config file path (default: config.yaml)")
	root.PersistentFlags().StringP("log-level", "l", "info", "Log level: debug, info, warn, error")
	root.PersistentFlags().BoolP("strict", "s", false, "Enable strict mode validation")

	// Validate log-level enum values
	root.PreRunE = func(cmd *cobra.Command, args []string) error {
		level, _ := cmd.Flags().GetString("log-level")
		validLevels := []string{"debug", "info", "warn", "error"}
		for _, valid := range validLevels {
			if level == valid {
				return nil
			}
		}
		return fmt.Errorf("invalid --log-level %q: must be one of: debug, info, warn, error", level)
	}

	// Initialize logger based on config
	logger := logging.NewLogger(cfg.LogLevel)
	slog.SetDefault(logger)

	return &Registry{
		root:   root,
		cfg:    cfg,
		logger: logger,
	}, nil
}

// Root returns the root cobra command.
func (r *Registry) Root() *cobra.Command {
	return r.root
}

// AddCommand adds a subcommand to the root.
func (r *Registry) AddCommand(cmd *cobra.Command) {
	r.root.AddCommand(cmd)
}

// Execute runs the command with fang styling.
func (r *Registry) Execute(ctx context.Context) error {
	return fang.Execute(ctx, r.root)
}

// ExecuteAndExit runs the command and exits with appropriate code.
func (r *Registry) ExecuteAndExit(ctx context.Context) {
	if err := r.Execute(ctx); err != nil {
		// fang handles the error styling
		os.Exit(1)
	}
}

// Validate runs validation on the command tree.
func (r *Registry) Validate() error {
	if r.validator == nil {
		return fmt.Errorf("validator not set")
	}
	return r.validator.ValidateCommandTree(r.root)
}

// SetValidator sets the validator for this registry.
func (r *Registry) SetValidator(v *validation.Validator) {
	r.validator = v
}

// HealthCheck implements the Healthchecker interface.
func (r *Registry) HealthCheck() error {
	if r.validator != nil {
		return r.validator.ValidateCommandTree(r.root)
	}
	return nil
}

// SetupCommands adds standard cmdguard commands to the registry.
func (r *Registry) SetupCommands() error {
	// Add validate command
	r.AddCommand(r.createValidateCommand())

	// Add version command
	r.AddCommand(r.createVersionCommand())

	return nil
}

// createValidateCommand creates the validation command.
func (r *Registry) createValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Run validation on the command tree",
		Long:  "Validates that all commands have handlers and all flags are properly bound.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if r.validator == nil {
				return fmt.Errorf("validator not initialized")
			}

			if err := r.validator.ValidateAll(); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			slog.Info("All commands and flags validated successfully")
			fmt.Fprintln(cmd.OutOrStdout(), "✓ All commands and flags validated successfully")
			return nil
		},
	}
}

// createVersionCommand creates the version command.
func (r *Registry) createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("version command executed", "version", "0.1.0")
			fmt.Fprintln(cmd.OutOrStdout(), "cmdguard version 0.1.0")
		},
	}
}

// GetConfig returns the configuration.
func (r *Registry) GetConfig() *config.Config {
	return r.cfg
}

// IsStrictMode returns true if strict mode is enabled.
func (r *Registry) IsStrictMode() bool {
	return r.cfg != nil && r.cfg.StrictMode
}
