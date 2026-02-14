// cmdguard is a CLI validation library that ensures every command and flag
// is properly implemented and bound.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/larsartmann/cmdguard/internal/commands"
	"github.com/larsartmann/cmdguard/internal/di"
	"github.com/larsartmann/cmdguard/internal/validation"
	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		// Use fang for styled error output
		ctx := context.Background()
		errorCmd := &cobra.Command{
			RunE: func(c *cobra.Command, args []string) error {
				return err
			},
		}
		_ = fang.Execute(ctx, errorCmd)
		os.Exit(1)
	}
}

func run() error {
	// Create root scope
	module := di.NewModule()

	// Register services
	if err := module.ProvideServices(); err != nil {
		return fmt.Errorf("failed to provide services: %w", err)
	}

	// Get services
	registry := module.MustInvokeRegistry()
	validator := module.MustInvokeValidator()

	// Link validator to registry
	registry.SetValidator(validator)

	// Add subcommands
	setupCommands(registry, validator)

	// Run validation on startup
	if err := registry.Validate(); err != nil {
		return fmt.Errorf("startup validation failed: %w", err)
	}

	// Execute with fang styling
	ctx := context.Background()
	if err := registry.Execute(ctx); err != nil {
		return err
	}

	// Graceful shutdown
	if err := module.Shutdown(); err != nil {
		// Log shutdown error but don't fail
		fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
	}

	return nil
}

func setupCommands(registry *commands.Registry, validator *validation.Validator) {
	// Add validate command
	registry.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Run validation on the command tree",
		Long:  "Validates that all commands have handlers and all flags are properly bound.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validator.ValidateAll(); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "✓ All commands and flags validated successfully")
			return nil
		},
	})

	// Add version command
	registry.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "cmdguard version 0.1.0")
		},
	})

	// Add example command with flags
	exampleCmd := &cobra.Command{
		Use:   "example",
		Short: "An example command with flags",
		Long:  "Demonstrates how flags are validated.",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			count, _ := cmd.Flags().GetInt("count")

			fmt.Fprintf(cmd.OutOrStdout(), "Hello %s! Count: %d\n", name, count)
			return nil
		},
	}

	exampleCmd.Flags().String("name", "World", "name to greet")
	exampleCmd.Flags().Int("count", 1, "number of greetings")

	registry.AddCommand(exampleCmd)
}
