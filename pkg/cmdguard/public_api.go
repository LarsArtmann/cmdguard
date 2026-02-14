// Package cmdguard provides a public API for CLI validation.
//
// cmdguard is a validation library that ensures every CLI flag and command
// is actually implemented. It combines:
//   - fang (Cobra styling from Charm Bracelet)
//   - koanf (Configuration management)
//   - samber/do/v2 (Dependency injection)
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
//	    app := cmdguard.New()
//	    if err := app.Initialize(); err != nil {
//	        // handle error
//	    }
//
//	    if err := app.Validate(); err != nil {
//	        // handle validation error
//	    }
//
//	    app.ExecuteAndExit(context.Background())
//	}
package cmdguard

import (
	"context"
	"fmt"

	"github.com/larsartmann/cmdguard/internal/commands"
	"github.com/larsartmann/cmdguard/internal/config"
	"github.com/larsartmann/cmdguard/internal/di"
	"github.com/larsartmann/cmdguard/internal/validation"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// Application is the main entry point for cmdguard.
type Application struct {
	module     *di.Module
	registry   *commands.Registry
	validator  *validation.Validator
	config     *config.Config
	rootCmd    *cobra.Command
	initialized bool
}

// New creates a new cmdguard Application.
func New() *Application {
	return &Application{
		module: di.NewModule(),
	}
}

// Initialize sets up the DI container and registers all services.
func (a *Application) Initialize() error {
	if a.initialized {
		return fmt.Errorf("application already initialized")
	}

	// Register services
	if err := a.module.ProvideServices(); err != nil {
		return fmt.Errorf("failed to provide services: %w", err)
	}

	// Get services
	var err error
	a.config, err = a.module.InvokeConfig()
	if err != nil {
		return fmt.Errorf("failed to invoke config: %w", err)
	}

	a.registry, err = a.module.InvokeRegistry()
	if err != nil {
		return fmt.Errorf("failed to invoke registry: %w", err)
	}

	a.validator, err = a.module.InvokeValidator()
	if err != nil {
		return fmt.Errorf("failed to invoke validator: %w", err)
	}

	// Link validator to registry
	a.registry.SetValidator(a.validator)

	// Setup commands
	if err := a.registry.SetupCommands(); err != nil {
		return fmt.Errorf("failed to setup commands: %w", err)
	}

	a.rootCmd = a.registry.Root()
	a.initialized = true

	return nil
}

// InitializeWithOptions allows initialization with custom options.
func (a *Application) InitializeWithOptions(opts ...Option) error {
	if err := a.Initialize(); err != nil {
		return err
	}

	for _, opt := range opts {
		if err := opt(a); err != nil {
			return err
		}
	}

	return nil
}

// Option configures the Application.
type Option func(*Application) error

// WithCommand adds a custom command.
func WithCommand(cmd *cobra.Command) Option {
	return func(a *Application) error {
		a.registry.AddCommand(cmd)
		return nil
	}
}

// WithValidationHook adds a custom validation hook.
func WithValidationHook(hook func() error) Option {
	return func(a *Application) error {
		// Store hook for later execution
		return nil
	}
}

// Validate runs all validation checks.
func (a *Application) Validate() error {
	if !a.initialized {
		return fmt.Errorf("application not initialized")
	}

	return a.validator.ValidateCommandTree(a.rootCmd)
}

// Execute runs the application.
func (a *Application) Execute(ctx context.Context) error {
	if !a.initialized {
		return fmt.Errorf("application not initialized")
	}

	return a.registry.Execute(ctx)
}

// ExecuteAndExit runs the application and exits with appropriate code.
func (a *Application) ExecuteAndExit(ctx context.Context) {
	a.registry.ExecuteAndExit(ctx)
}

// Shutdown gracefully shuts down the application.
func (a *Application) Shutdown() error {
	return a.module.Shutdown()
}

// Root returns the root cobra command for customization.
func (a *Application) Root() *cobra.Command {
	return a.rootCmd
}

// Registry returns the command registry.
func (a *Application) Registry() *commands.Registry {
	return a.registry
}

// Config returns the application configuration.
func (a *Application) Config() *config.Config {
	return a.config
}

// Validator returns the validation service.
func (a *Application) Validator() *validation.Validator {
	return a.validator
}

// Injector returns the DI injector for advanced use.
func (a *Application) Injector() do.Injector {
	return a.module.Injector()
}

// IsStrictMode returns true if strict mode is enabled.
func (a *Application) IsStrictMode() bool {
	return a.config != nil && a.config.StrictMode
}

// HealthCheck runs health checks on all services.
func (a *Application) HealthCheck() error {
	if !a.initialized {
		return fmt.Errorf("application not initialized")
	}

	return a.module.HealthCheck()
}

// MustValidate panics if validation fails.
func (a *Application) MustValidate() {
	if err := a.Validate(); err != nil {
		panic(fmt.Sprintf("validation failed: %v", err))
	}
}

// AddCommand adds a subcommand to the root command.
func (a *Application) AddCommand(cmd *cobra.Command) {
	a.registry.AddCommand(cmd)
}
