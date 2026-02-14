// Package di provides dependency injection configuration using samber/do/v2.
package di

import (
	"context"

	"github.com/larsartmann/cmdguard/internal/commands"
	"github.com/larsartmann/cmdguard/internal/config"
	"github.com/larsartmann/cmdguard/internal/validation"
	"github.com/samber/do/v2"
)

// Module provides DI container setup for cmdguard.
type Module struct {
	injector do.Injector
}

// NewModule creates a new DI module with the root scope.
func NewModule() *Module {
	return &Module{
		injector: do.New(),
	}
}

// Injector returns the underlying injector.
func (m *Module) Injector() do.Injector {
	return m.injector
}

// ProvideServices registers all services with the DI container.
func (m *Module) ProvideServices() error {
	// Lazy services - only created when needed
	do.Provide(m.injector, config.NewConfig)
	do.Provide(m.injector, validation.NewRegistry)
	do.Provide(m.injector, validation.NewValidator)
	do.Provide(m.injector, commands.NewRegistry)

	// Transient services - new instance per injection
	do.ProvideTransient(m.injector, validation.NewFlagValidator)

	return nil
}

// CreateChildScope creates a child scope for a specific command.
func (m *Module) CreateChildScope(name string) do.Injector {
	return m.injector.Scope(name)
}

// InvokeConfig retrieves the Config service.
func (m *Module) InvokeConfig() (*config.Config, error) {
	return do.Invoke[*config.Config](m.injector)
}

// InvokeRegistry retrieves the CommandRegistry service.
func (m *Module) InvokeRegistry() (*commands.Registry, error) {
	return do.Invoke[*commands.Registry](m.injector)
}

// InvokeValidator retrieves the Validator service.
func (m *Module) InvokeValidator() (*validation.Validator, error) {
	return do.Invoke[*validation.Validator](m.injector)
}

// HealthCheck runs health checks on all services.
func (m *Module) HealthCheck() error {
	// Check config health
	if err := do.HealthCheck[*config.Config](m.injector); err != nil {
		return err
	}

	// Check registry health
	if err := do.HealthCheck[*validation.Registry](m.injector); err != nil {
		return err
	}

	// Check validator health
	if err := do.HealthCheck[*validation.Validator](m.injector); err != nil {
		return err
	}

	// Check commands registry health
	if err := do.HealthCheck[*commands.Registry](m.injector); err != nil {
		return err
	}

	return nil
}

// HealthCheckWithContext runs health checks with context.
func (m *Module) HealthCheckWithContext(ctx context.Context) error {
	return do.HealthCheckWithContext[*config.Config](ctx, m.injector)
}

// Shutdown gracefully shuts down all services.
func (m *Module) Shutdown() error {
	// Shutdown individual services in reverse order of dependency
	var errs []error

	if err := do.Shutdown[*commands.Registry](m.injector); err != nil {
		errs = append(errs, err)
	}

	if err := do.Shutdown[*validation.Validator](m.injector); err != nil {
		errs = append(errs, err)
	}

	if err := do.Shutdown[*validation.Registry](m.injector); err != nil {
		errs = append(errs, err)
	}

	if err := do.Shutdown[*config.Config](m.injector); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errs[0] // Return first error
	}

	return nil
}

// ShutdownWithContext gracefully shuts down with timeout.
func (m *Module) ShutdownWithContext(ctx context.Context) error {
	// For now, delegate to regular shutdown
	// Context handling would be added in production
	return m.Shutdown()
}

// MustInvokeConfig retrieves Config or panics.
func (m *Module) MustInvokeConfig() *config.Config {
	return do.MustInvoke[*config.Config](m.injector)
}

// MustInvokeRegistry retrieves Registry or panics.
func (m *Module) MustInvokeRegistry() *commands.Registry {
	return do.MustInvoke[*commands.Registry](m.injector)
}

// MustInvokeValidator retrieves Validator or panics.
func (m *Module) MustInvokeValidator() *validation.Validator {
	return do.MustInvoke[*validation.Validator](m.injector)
}
