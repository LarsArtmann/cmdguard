package validation

import (
	"fmt"
	"strings"

	"github.com/larsartmann/cmdguard/internal/config"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// Validator performs comprehensive validation on commands and flags.
type Validator struct {
	registry *Registry
	cfg      *config.Config
}

// NewValidator creates a new Validator instance.
func NewValidator(i do.Injector) (*Validator, error) {
	registry, err := do.Invoke[*Registry](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke registry: %w", err)
	}

	cfg, err := do.Invoke[*config.Config](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke config: %w", err)
	}

	return &Validator{
		registry: registry,
		cfg:      cfg,
	}, nil
}

// ValidateAll performs all validation checks.
func (v *Validator) ValidateAll() error {
	if err := v.ValidateCommands(); err != nil {
		return err
	}

	if err := v.ValidateFlags(); err != nil {
		return err
	}

	return nil
}

// ValidateCommands checks that all registered commands have proper handlers.
func (v *Validator) ValidateCommands() error {
	commands := v.registry.GetCommands()

	var errors []string
	for name, cmd := range commands {
		if err := v.validateCommand(name, cmd); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("command validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// ValidateFlags checks that all declared flags are properly bound.
func (v *Validator) ValidateFlags() error {
	commands := v.registry.GetCommands()

	var errors []string
	for cmdName := range commands {
		flags, _ := v.registry.GetFlags(cmdName)
		for _, flag := range flags {
			if !flag.IsBound {
				errors = append(errors, fmt.Sprintf("flag %q on command %q is not bound", flag.Name, cmdName))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("flag validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// ValidateCommandTree validates an entire cobra command tree.
func (v *Validator) ValidateCommandTree(root *cobra.Command) error {
	// Register the root command
	if err := v.registry.RegisterCommand(root); err != nil {
		return fmt.Errorf("failed to register root command: %w", err)
	}

	// Register all subcommands recursively
	if err := v.registerSubcommands(root); err != nil {
		return fmt.Errorf("failed to register subcommands: %w", err)
	}

	// Now validate
	return v.ValidateAll()
}

// registerSubcommands recursively registers all subcommands.
func (v *Validator) registerSubcommands(parent *cobra.Command) error {
	for _, child := range parent.Commands() {
		if err := v.registry.RegisterSubcommand(parent, child); err != nil {
			return err
		}
		// Recursively register child's subcommands
		if err := v.registerSubcommands(child); err != nil {
			return err
		}
	}
	return nil
}

// validateCommand validates a single command.
func (v *Validator) validateCommand(name string, cmd *CommandInfo) error {
	// Commands with subcommands don't need a handler
	if cmd.HasSubcommands {
		return nil
	}

	// Leaf commands must have a handler
	if cmd.Handler == nil {
		return fmt.Errorf("command %q has no handler", name)
	}

	return nil
}

// IsStrictMode returns true if strict validation is enabled.
func (v *Validator) IsStrictMode() bool {
	return v.cfg != nil && v.cfg.StrictMode
}

// HealthCheck implements the Healthchecker interface.
func (v *Validator) HealthCheck() error {
	return v.ValidateAll()
}

// FlagValidator provides per-validation flag checking.
type FlagValidator struct {
	strict bool
}

// NewFlagValidator creates a new transient FlagValidator.
func NewFlagValidator(i do.Injector) (*FlagValidator, error) {
	cfg, err := do.Invoke[*config.Config](i)
	if err != nil {
		return nil, err
	}

	return &FlagValidator{
		strict: cfg.StrictMode,
	}, nil
}

// ValidateFlag checks if a flag value is valid.
func (v *FlagValidator) ValidateFlag(name string, value interface{}) error {
	if v.strict && value == nil {
		return fmt.Errorf("flag %q is required in strict mode", name)
	}
	return nil
}

// ValidateFlagAccess checks if a flag is properly registered before access.
func (v *FlagValidator) ValidateFlagAccess(cmd *cobra.Command, flagName string) error {
	if cmd.Flags().Lookup(flagName) == nil {
		return fmt.Errorf("flag %q is not registered on command %q", flagName, cmd.Name())
	}
	return nil
}
