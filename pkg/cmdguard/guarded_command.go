// Package cmdguard provides compile-time guarded CLI construction.
//
// The guard approach panics at construction time if commands are invalid,
// ensuring errors are caught immediately rather than at runtime.
//
// Basic usage:
//
//	package main
//
//	import (
//	    "context"
//	    "github.com/larsartmann/cmdguard/pkg/cmdguard"
//	)
//
//	func main() {
//	    // Single-step initialization - panics on invalid
//	    root := cmdguard.New("myapp", "My application")
//
//	    // This will panic if command has no handler
//	    root.AddCommand(&cobra.Command{
//	    Use:   "sub",
//	        Short: "Subcommand",
//	        Run: func(cmd *cobra.Command, args []string) {
//	            // handler
//	        },
//	    })
//
//	    // Execute
//	    root.Execute(context.Background())
//	}
package cmdguard

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/larsartmann/cmdguard/internal/config"
	"github.com/larsartmann/cmdguard/internal/logging"
	"github.com/spf13/cobra"
)

// version is set at build time via ldflags:
// go build -ldflags "-X github.com/larsartmann/cmdguard/pkg/cmdguard.version=X"
var version = "dev"

// GuardedCommand wraps a cobra.Command with compile-time validation.
// It panics on construction if commands are invalid, ensuring errors
// are caught immediately at startup rather than at runtime.
type GuardedCommand struct {
	cmd            *cobra.Command
	cfg            *config.Config
	logger         *slog.Logger
	validated      bool
	strictMode     bool
	registeredCmds map[string]bool // tracks registered command names for duplicate detection
}

// New creates a new GuardedCommand with the given name and description.
// This is the single entry point for creating a guarded CLI application.
//
// Example:
//
//	root := cmdguard.New("myapp", "My application description")
//	root.Execute(context.Background())
func New(name, short string) *GuardedCommand {
	// Load configuration early
	cfg := config.Load()

	// Initialize logger
	logger := logging.NewLogger(cfg.LogFormat, cfg.LogLevel)
	slog.SetDefault(logger)

	// Create root command
	cmd := &cobra.Command{
		Use:           name,
		Short:         short,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add global flags
	cmd.PersistentFlags().StringP("config", "c", "", "Config file path")
	cmd.PersistentFlags().StringP("log-level", "l", cfg.LogLevel, "Log level: debug, info, warn, error")
	cmd.PersistentFlags().BoolP("strict", "s", cfg.StrictMode, "Enable strict mode validation")

	// Validate log-level in PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		level, _ := cmd.Flags().GetString("log-level")
		validLevels := []string{"debug", "info", "warn", "error"}
		if !slices.Contains(validLevels, level) {
			return fmt.Errorf("invalid --log-level %q: must be one of: debug, info, warn, error", level)
		}
		return nil
	}

	g := &GuardedCommand{
		cmd:            cmd,
		cfg:            cfg,
		logger:         logger,
		strictMode:     cfg.StrictMode,
		registeredCmds: make(map[string]bool),
	}

	// Add default commands
	g.addDefaultCommands()

	return g
}

// AddCommand adds a subcommand to the guarded command.
// PANICS if the command is invalid (no handler and no subcommands).
// PANICS if a command with the same name already exists.
//
// This is intentional - it ensures errors are caught at startup
// rather than when the command is invoked.
func (g *GuardedCommand) AddCommand(cmd *cobra.Command) {
	if g.validated {
		panic("cmdguard: cannot add commands after execution")
	}

	// Check for duplicate command name
	if g.registeredCmds[cmd.Name()] {
		panic(fmt.Sprintf("cmdguard: duplicate command %q", cmd.Name()))
	}

	// Check for duplicate subcommand names within this command
	g.checkDuplicateSubcommands(cmd)

	// Validate command before adding
	if err := g.validateCommand(cmd); err != nil {
		panic(fmt.Sprintf("cmdguard: invalid command %q: %v", cmd.Name(), err))
	}

	g.cmd.AddCommand(cmd)
	g.registeredCmds[cmd.Name()] = true
	g.logger.Debug("added command", "command", cmd.Name())
}

// AddSubcommand adds a subcommand to a parent command.
// PANICS if the subcommand is invalid.
// PANICS if a subcommand with the same name already exists under the parent.
func (g *GuardedCommand) AddSubcommand(parent, child *cobra.Command) {
	if g.validated {
		panic("cmdguard: cannot add commands after execution")
	}

	// Check for duplicate subcommand name under this parent
	for _, existing := range parent.Commands() {
		if existing.Name() == child.Name() {
			panic(fmt.Sprintf("cmdguard: duplicate subcommand %q in command %q", child.Name(), parent.Name()))
		}
	}

	// Validate child before adding
	if err := g.validateCommand(child); err != nil {
		panic(fmt.Sprintf("cmdguard: invalid subcommand %q: %v", child.Name(), err))
	}

	parent.AddCommand(child)
	g.logger.Debug("added subcommand",
		"parent", parent.Name(),
		"child", child.Name(),
	)
}

// Execute runs the command with the given context.
func (g *GuardedCommand) Execute(ctx context.Context) error {
	g.validated = true
	return fang.Execute(ctx, g.cmd)
}

// ExecuteAndExit runs the command and exits with appropriate exit code.
func (g *GuardedCommand) ExecuteAndExit(ctx context.Context) {
	if err := g.Execute(ctx); err != nil {
		// fang handles error styling
		os.Exit(1)
	}
}

// Command returns the underlying cobra command for advanced customization.
// Use with caution - modifications bypass guard validation.
func (g *GuardedCommand) Command() *cobra.Command {
	return g.cmd
}

// Config returns the application configuration.
func (g *GuardedCommand) Config() *config.Config {
	return g.cfg
}

// IsStrictMode returns true if strict mode is enabled.
func (g *GuardedCommand) IsStrictMode() bool {
	return g.strictMode
}

// Version returns the current version string.
func Version() string {
	return version
}

// validateCommand checks if a command is valid.
// Returns error if command has no handler and no subcommands.
func (g *GuardedCommand) validateCommand(cmd *cobra.Command) error {
	// Check for command name
	if cmd.Name() == "" {
		return fmt.Errorf("command has no name")
	}

	// Commands with subcommands don't need a handler
	if len(cmd.Commands()) > 0 {
		return nil
	}

	// Check for handler
	hasRun := cmd.Run != nil
	hasRunE := cmd.RunE != nil

	if !hasRun && !hasRunE {
		return fmt.Errorf("command has no handler (Run or RunE) and no subcommands")
	}

	// In strict mode, require RunE (returns error)
	if g.strictMode && !hasRunE {
		return fmt.Errorf("strict mode requires RunE handler that returns error")
	}

	return nil
}

// addDefaultCommands adds the built-in commands.
func (g *GuardedCommand) addDefaultCommands() {
	// Add version command
	g.cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "cmdguard version "+version)
		},
	})

	// Add validate command (self-validation)
	g.cmd.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validate command tree",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := g.validateCommandTree(); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), "✓ All commands validated successfully"); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			return nil
		},
	})
}

// validateCommandTree validates all commands in the tree.
func (g *GuardedCommand) validateCommandTree() error {
	var errors []string

	for _, cmd := range g.cmd.Commands() {
		if err := g.validateCommand(cmd); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", cmd.Name(), err))
		}
		// Recursively validate subcommands
		if err := g.validateSubcommands(cmd); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// checkDuplicateSubcommands checks for duplicate subcommand names within a command.
// PANICS if duplicates are found.
func (g *GuardedCommand) checkDuplicateSubcommands(parent *cobra.Command) {
	seen := make(map[string]bool)
	for _, cmd := range parent.Commands() {
		if seen[cmd.Name()] {
			panic(fmt.Sprintf("cmdguard: duplicate subcommand %q in command %q", cmd.Name(), parent.Name()))
		}
		seen[cmd.Name()] = true
	}
}

// validateSubcommands recursively validates subcommands.
func (g *GuardedCommand) validateSubcommands(parent *cobra.Command) error {
	for _, cmd := range parent.Commands() {
		if err := g.validateCommand(cmd); err != nil {
			return fmt.Errorf("%s %s: %v", parent.Name(), cmd.Name(), err)
		}
		if err := g.validateSubcommands(cmd); err != nil {
			return err
		}
	}
	return nil
}
