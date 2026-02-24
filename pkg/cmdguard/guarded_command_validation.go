package cmdguard

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

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
			return fmt.Errorf("%s %s: %w", parent.Name(), cmd.Name(), err)
		}
		if err := g.validateSubcommands(cmd); err != nil {
			return err
		}
	}
	return nil
}
