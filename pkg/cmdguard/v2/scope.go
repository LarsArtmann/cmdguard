package v2

import (
	"context"
	"fmt"

	"github.com/samber/do/v2"
)

// Scope provides DI scope management for CLI applications.
// It wraps samber/do/v2 injector with convenience methods.
type Scope struct {
	injector do.Injector
	name     string
	parent   *Scope
}

// NewScope creates a new root scope.
func NewScope(name string) *Scope {
	return &Scope{
		injector: do.New(),
		name:     name,
		parent:   nil,
	}
}

// NewScopeFromInjector creates a Scope from an existing injector.
func NewScopeFromInjector(injector do.Injector, name string) *Scope {
	return &Scope{
		injector: injector,
		name:     name,
		parent:   nil,
	}
}

// Child creates a child scope that inherits services from this scope.
// Child scopes can override parent services but not affect them.
func (s *Scope) Child(name string) *Scope {
	return &Scope{
		injector: s.injector.Scope(name),
		name:     name,
		parent:   s,
	}
}

// Name returns the scope name.
func (s *Scope) Name() string {
	return s.name
}

// Parent returns the parent scope, or nil if this is a root scope.
func (s *Scope) Parent() *Scope {
	return s.parent
}

// Injector returns the underlying DI injector.
// Use this for direct samber/do/v2 operations.
func (s *Scope) Injector() do.Injector {
	return s.injector
}

// Provide registers a service provider in this scope.
// Returns an error if registration fails.
func Provide[T any](scope *Scope, provider func(do.Injector) (T, error)) error {
	if scope == nil {
		return fmt.Errorf("%w: scope is nil", ErrInvalidScope)
	}

	do.Provide(scope.injector, provider)

	return nil
}

// ProvideNamed registers a named service provider in this scope.
// Use this when you need to register multiple implementations of the same interface.
// Returns an error if registration fails.
func ProvideNamed[T any](scope *Scope, name string, provider func(do.Injector) (T, error)) error {
	if scope == nil {
		return fmt.Errorf("%w: scope is nil", ErrInvalidScope)
	}

	do.ProvideNamed(scope.injector, name, provider)

	return nil
}

// ProvideValue registers a value directly in this scope.
// Useful for registering already-constructed services.
func ProvideValue[T any](scope *Scope, value T) error {
	if scope == nil {
		return fmt.Errorf("%w: scope is nil", ErrInvalidScope)
	}

	do.ProvideValue(scope.injector, value)

	return nil
}

// Invoke retrieves a service from the scope.
// Returns an error if the service is not found or construction fails.
func Invoke[T any](scope *Scope) (T, error) {
	var zero T
	if scope == nil {
		return zero, fmt.Errorf("%w: scope is nil", ErrInvalidScope)
	}

	return do.Invoke[T](scope.injector)
}

// InvokeNamed retrieves a named service from the scope.
// Returns an error if the service is not found or construction fails.
func InvokeNamed[T any](scope *Scope, name string) (T, error) {
	var zero T
	if scope == nil {
		return zero, fmt.Errorf("%w: scope is nil", ErrInvalidScope)
	}

	return do.InvokeNamed[T](scope.injector, name)
}

// Shutdown gracefully shuts down all services in this scope.
// Services implementing the Shutdowner interface will be notified.
func (s *Scope) Shutdown(ctx context.Context) error {
	if s.injector == nil {
		return nil
	}

	report := s.injector.ShutdownWithContext(ctx)
	if report.Succeed {
		return nil
	}

	return report
}

// ShutdownAll shuts down this scope and all parent scopes.
func (s *Scope) ShutdownAll(ctx context.Context) error {
	var errs []error

	// Shutdown from child to parent
	current := s
	for current != nil {
		err := current.Shutdown(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("scope %q: %w", current.name, err))
		}

		current = current.parent
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %v", ErrServiceConstruction, errs)
	}

	return nil
}

// HealthCheck runs health checks on all services in this scope.
func (s *Scope) HealthCheck() error {
	if s.injector == nil {
		return nil
	}

	results := s.injector.HealthCheck()
	for _, err := range results {
		if err != nil {
			return err
		}
	}

	return nil
}

// HealthCheckWithContext runs health checks with context on all services.
// Services implementing HealthcheckerWithContext will use the provided context.
func (s *Scope) HealthCheckWithContext(ctx context.Context) error {
	if s.injector == nil {
		return nil
	}

	results := s.injector.HealthCheckWithContext(ctx)
	for _, err := range results {
		if err != nil {
			return err
		}
	}

	return nil
}

// ScopedProvider creates a provider that runs within a named child scope.
// Useful for plugins that need their own isolated scope.
func ScopedProvider[T any](
	parent *Scope,
	scopeName string,
	provider func(do.Injector) (T, error),
) func(do.Injector) (T, error) {
	return func(i do.Injector) (T, error) {
		childScope := parent.Child(scopeName)

		return provider(childScope.Injector())
	}
}

// RegisterInScope registers providers in a child scope.
// Returns the child scope for further operations.
func RegisterInScope(parent *Scope, name string, providers ...any) (*Scope, error) {
	if parent == nil {
		return nil, fmt.Errorf("%w: parent scope is nil", ErrInvalidScope)
	}

	child := parent.Child(name)

	for providerIndex, provider := range providers {
		switch fn := provider.(type) {
		case func(do.Injector) (any, error):
			do.Provide(child.injector, fn)
		default:
			return nil, fmt.Errorf(
				"%w: provider %d has invalid type %T, expected func(do.Injector) (T, error)",
				ErrServiceRegistration,
				providerIndex,
				provider,
			)
		}
	}

	return child, nil
}

// IsRoot returns true if this is a root scope (has no parent).
func (s *Scope) IsRoot() bool {
	return s.parent == nil
}

// Path returns the full scope path from root to this scope.
func (s *Scope) Path() []string {
	var path []string

	current := s
	for current != nil {
		path = append([]string{current.name}, path...)
		current = current.parent
	}

	return path
}
