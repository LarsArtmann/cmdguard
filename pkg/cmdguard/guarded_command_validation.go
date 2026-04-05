package cmdguard

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	ErrCommandNoName          = errors.New("command has no name")
	ErrCommandNoHandler       = errors.New("command has no handler (Run or RunE) and no subcommands")
	ErrStrictModeRequiresRunE = errors.New("strict mode requires RunE handler that returns error")
	ErrValidationFailed       = errors.New("validation failed")
	ErrInvalidLogLevel        = errors.New("invalid log level")
)

// validateCommand checks if a command is valid.
// Returns error if command has no handler and no subcommands.
func (g *GuardedCommand) validateCommand(cmd *cobra.Command) error {
	if cmd.Name() == "" {
		return ErrCommandNoName
	}

	// Commands with subcommands don't need a handler
	if len(cmd.Commands()) > 0 {
		return nil
	}

	// Check for handler
	hasRun := cmd.Run != nil
	hasRunE := cmd.RunE != nil

	if !hasRun && !hasRunE {
		return ErrCommandNoHandler
	}

	// In strict mode, require RunE (returns error)
	if g.strictMode && !hasRunE {
		return ErrStrictModeRequiresRunE
	}

	return nil
}

// validateCommandTree validates all commands in the tree.
func (g *GuardedCommand) validateCommandTree() error {
	var errs []string

	for _, cmd := range g.cmd.Commands() {
		err := g.validateCommand(cmd)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", cmd.Name(), err))
		}
		// Recursively validate subcommands
		err = g.validateSubcommands(cmd)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w:\n  - %s", ErrValidationFailed, strings.Join(errs, "\n  - "))
	}

	return nil
}

// checkDuplicateSubcommands checks for duplicate subcommand names within a command.
// PANICS if duplicates are found.
func (g *GuardedCommand) checkDuplicateSubcommands(parent *cobra.Command) {
	seen := make(map[string]bool)
	for _, cmd := range parent.Commands() {
		if seen[cmd.Name()] {
			panic(
				fmt.Sprintf(
					"cmdguard: duplicate subcommand %q in command %q",
					cmd.Name(),
					parent.Name(),
				),
			)
		}

		seen[cmd.Name()] = true
	}
}

// validateSubcommands recursively validates subcommands.
func (g *GuardedCommand) validateSubcommands(parent *cobra.Command) error {
	for _, cmd := range parent.Commands() {
		err := g.validateCommand(cmd)
		if err != nil {
			return fmt.Errorf("%s %s: %w", parent.Name(), cmd.Name(), err)
		}

		err = g.validateSubcommands(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}
