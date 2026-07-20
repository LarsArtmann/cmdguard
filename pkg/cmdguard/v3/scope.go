package v3

import (
	"context"
	"errors"
	"fmt"
	"slices"

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

// NewScopeWithOpts creates a new root scope with custom injector options.
// Use this for DI logging, lifecycle hooks, health check timeouts, etc.
func NewScopeWithOpts(name string, opts *do.InjectorOpts) *Scope {
	return &Scope{
		injector: do.NewWithOpts(opts),
		name:     name,
		parent:   nil,
	}
}

// NewScopeFromInjector creates a Scope from an existing injector.
// Returns an error if injector is nil.
func NewScopeFromInjector(injector do.Injector, name string) (*Scope, error) {
	if injector == nil {
		return nil, fmt.Errorf("%w: injector is nil, name=%q", ErrInvalidScope, name)
	}

	return &Scope{
		injector: injector,
		name:     name,
		parent:   nil,
	}, nil
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

// RootScope returns the root scope by navigating up the parent chain.
// For a root scope, returns itself.
func (s *Scope) RootScope() *Scope {
	return walkRoot(s, scopeParentOf)
}

// walkRoot walks up the parent chain to the topmost ancestor. hasParent
// returns the next node and true if a parent exists; false otherwise.
// Used by Scope.RootScope and BranchingFlowContext.Root.
func walkRoot[T any](node T, hasParent func(T) (T, bool)) T {
	current := node
	for {
		next, ok := hasParent(current)
		if !ok {
			return current
		}

		current = next
	}
}

func scopeParentOf(s *Scope) (*Scope, bool) {
	return s.parent, s.parent != nil
}

// safeProvide calls fn with panic recovery, returning an error on panic.
// This ensures the zero-panic guarantee holds even when samber/do panics
// on duplicate registrations.
func safeProvide(fn func(), context string) error {
	var panicVal any

	func() {
		defer func() {
			panicVal = recover()
		}()

		fn()
	}()

	if panicVal == nil {
		return nil
	}

	return fmt.Errorf("%w: %s: %v", ErrServiceRegistration, context, panicVal)
}

// Provide registers a service provider in this scope.
// The provider is invoked lazily on first Invoke call.
// Returns an error if scope is nil or if registration panics (e.g. duplicate service).
func Provide[T any](scope *Scope, provider func(do.Injector) (T, error)) error {
	if scope == nil {
		return fmt.Errorf("%w: scope is nil, provider=%T", ErrInvalidScope, provider)
	}

	return safeProvide(func() { do.Provide(scope.injector, provider) },
		fmt.Sprintf("Provide[%T]", provider))
}

// ProvideNamed registers a named service provider in this scope.
// Use this when you need to register multiple implementations of the same interface.
// Returns an error if scope is nil or if registration panics (e.g. duplicate service).
func ProvideNamed[T any](scope *Scope, name string, provider func(do.Injector) (T, error)) error {
	if scope == nil {
		return fmt.Errorf("%w: scope is nil, name=%q, provider=%T", ErrInvalidScope, name, provider)
	}

	return safeProvide(func() { do.ProvideNamed(scope.injector, name, provider) },
		fmt.Sprintf("ProvideNamed[%T](%q)", provider, name))
}

// ProvideValue registers a value directly in this scope.
// Useful for registering already-constructed services.
// Returns an error if scope is nil or if registration panics (e.g. duplicate service).
func ProvideValue[T any](scope *Scope, value T) error {
	if scope == nil {
		return fmt.Errorf("%w: scope is nil, value type=%T", ErrInvalidScope, value)
	}

	return safeProvide(func() { do.ProvideValue(scope.injector, value) },
		fmt.Sprintf("ProvideValue[%T]", value))
}

// Invoke retrieves a service from the scope.
// Returns an error if the service is not found or construction fails.
// When the service is missing, the error chain includes ErrServiceNotFound so
// callers can distinguish not-found from construction failures via errors.Is.
func Invoke[T any](scope *Scope) (T, error) {
	var zero T
	if scope == nil {
		return zero, fmt.Errorf("%w: scope is nil, result type=%T", ErrInvalidScope, zero)
	}

	v, err := do.Invoke[T](scope.injector)
	if err != nil && errors.Is(err, do.ErrServiceNotFound) {
		return zero, fmt.Errorf("%w: type=%T, detail: %w", ErrServiceNotFound, zero, err)
	}

	return v, err
}

// InvokeNamed retrieves a named service from the scope.
// Returns an error if the service is not found or construction fails.
// When the service is missing, the error chain includes ErrServiceNotFound so
// callers can distinguish not-found from construction failures via errors.Is.
func InvokeNamed[T any](scope *Scope, name string) (T, error) {
	var zero T
	if scope == nil {
		return zero, fmt.Errorf(
			"%w: scope is nil, name=%q, result type=%T",
			ErrInvalidScope,
			name,
			zero,
		)
	}

	v, err := do.InvokeNamed[T](scope.injector, name)
	if err != nil && errors.Is(err, do.ErrServiceNotFound) {
		return zero, fmt.Errorf("%w: name=%q, type=%T, detail: %w", ErrServiceNotFound, name, zero, err)
	}

	return v, err
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

	return fmt.Errorf("%w: shutdown of scope %q: %w", ErrServiceConstruction, s.name, report)
}

// ShutdownAll shuts down this scope and all parent scopes.
func (s *Scope) ShutdownAll(ctx context.Context) error {
	var errs []error

	// Shutdown from child to parent
	current := s
	for current != nil {
		err := current.Shutdown(ctx)
		if err != nil {
			errs = append(errs, err)
		}

		current = current.parent
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %w", ErrServiceConstruction, errors.Join(errs...))
	}

	return nil
}

// firstHealthCheckError returns the first error from results, or nil.
func firstHealthCheckError(results map[string]error, scopeName string) error {
	for _, err := range results {
		if err != nil {
			return fmt.Errorf(
				"%w: health check in scope %q: %w",
				ErrServiceConstruction,
				scopeName,
				err,
			)
		}
	}

	return nil
}

// HealthCheck runs health checks on all services in this scope.
// Returns the first error found, or nil if all services are healthy.
// For per-service results, use HealthCheckResults instead.
func (s *Scope) HealthCheck() error {
	if s.injector == nil {
		return nil
	}

	results := s.injector.HealthCheck()

	return firstHealthCheckError(results, s.name)
}

// HealthCheckResults runs health checks and returns per-service results.
// The returned map keys are service names and values are their errors (nil = healthy).
func (s *Scope) HealthCheckResults() map[string]error {
	if s.injector == nil {
		return map[string]error{}
	}

	return s.injector.HealthCheck()
}

// HealthCheckWithContext runs health checks with context on all services.
// Services implementing HealthcheckerWithContext will use the provided context.
// Returns the first error found, or nil if all services are healthy.
// For per-service results, use HealthCheckResultsWithContext instead.
func (s *Scope) HealthCheckWithContext(ctx context.Context) error {
	if s.injector == nil {
		return nil
	}

	results := s.injector.HealthCheckWithContext(ctx)

	return firstHealthCheckError(results, s.name)
}

// HealthCheckResultsWithContext runs health checks with context and returns
// per-service results. The returned map keys are service names and values
// are their errors (nil = healthy).
func (s *Scope) HealthCheckResultsWithContext(ctx context.Context) map[string]error {
	if s.injector == nil {
		return map[string]error{}
	}

	return s.injector.HealthCheckWithContext(ctx)
}

// Override replaces a service provider in this scope.
// Useful for testing — replace real services with mocks in a cloned scope.
// Returns an error only if scope is nil.
func Override[T any](scope *Scope, replacement func(do.Injector) (T, error)) error {
	if scope == nil {
		return fmt.Errorf("%w: scope is nil, replacement=%T", ErrInvalidScope, replacement)
	}

	do.Override(scope.injector, replacement)

	return nil
}

// OverrideValue replaces a pre-constructed value in this scope.
// Useful for testing — inject config or mock values into a cloned scope.
// Returns an error only if scope is nil.
func OverrideValue[T any](scope *Scope, value T) error {
	if scope == nil {
		return fmt.Errorf("%w: scope is nil, value type=%T", ErrInvalidScope, value)
	}

	do.OverrideValue(scope.injector, value)

	return nil
}

// CloneScope creates a copy of the scope with the same service registrations
// but without the invoked service state. Use with Override/OverrideValue
// for test isolation.
func CloneScope(scope *Scope) *Scope {
	cloned := scope.injector.RootScope().Clone()

	return &Scope{
		injector: cloned,
		name:     scope.name,
		parent:   nil,
	}
}

// ScopedProvider creates a provider that runs within a named child scope.
// Useful for plugins that need their own isolated scope.
func ScopedProvider[T any](
	parent *Scope,
	scopeName string,
	provider func(do.Injector) (T, error),
) func(do.Injector) (T, error) {
	return func(_ do.Injector) (T, error) {
		childScope := parent.Child(scopeName)

		return provider(childScope.Injector())
	}
}

// RegisterInScope registers providers in a child scope.
// Returns the child scope for further operations.
func RegisterInScope(parent *Scope, name string, providers ...any) (*Scope, error) {
	if parent == nil {
		return nil, fmt.Errorf(
			"%w: parent scope is nil, name=%q, providers=%d",
			ErrInvalidScope,
			name,
			len(providers),
		)
	}

	child := parent.Child(name)

	for providerIndex, provider := range providers {
		switch fn := provider.(type) {
		case func(do.Injector) (any, error):
			do.Provide(child.injector, fn)
		default:
			return nil, fmt.Errorf(
				"%w: scope=%q, providers=%d, provider index=%d, provider type=%T",
				ErrServiceRegistration,
				name,
				len(providers),
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
	names := []string{}

	current := s
	for current != nil {
		names = append(names, current.name)
		current = current.parent
	}

	slices.Reverse(names)

	return names
}

// Package creates a CLI bound to a pre-existing DI scope and registers the
// CLI itself as a service in that scope for self-injection.
//
// Usage:
//
//	scope := v3.NewScope("app")
//	cli, err := v3.Package(scope, "app", "My app", Config{})
//	if err != nil {
//	    log.Fatal(err)
//	}
func Package[T any](scope *Scope, name, short string, defaults T, opts ...CLIOption) (*CLI[T], error) {
	cliOpts := append([]CLIOption{WithCLIScope(scope)}, opts...)

	cli, err := NewCLI(name, short, defaults, cliOpts...)
	if err != nil {
		return nil, fmt.Errorf("Package(name=%q, short=%q): %w", name, short, err)
	}

	do.ProvideValue(cli.spec.scope.injector, cli)

	return cli, nil
}
