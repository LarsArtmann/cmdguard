package v2

import (
	"context"
	"fmt"
	"os"

	"charm.land/fang/v2"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// Execute runs the CLI application.
// Returns an error if execution fails (never panics).
// Uses fang for beautiful error styling.
func (g *GuardedCommand[T, F]) Execute(ctx context.Context) error {
	if err := fang.Execute(ctx, g.rootCmd); err != nil {
		return fmt.Errorf("failed to execute CLI: %w", err)
	}

	return nil
}

// ExecuteWithArgs runs the CLI application with specific arguments.
// Useful for testing.
func (g *GuardedCommand[T, F]) ExecuteWithArgs(ctx context.Context, args []string) error {
	g.rootCmd.SetArgs(args)

	return g.Execute(ctx)
}

// ExecuteAndExit runs the CLI and exits with the appropriate exit code.
// This is the simplest way to run a CLI application.
// Uses fang for beautiful error styling.
func (g *GuardedCommand[T, F]) ExecuteAndExit(ctx context.Context) {
	err := g.Execute(ctx)
	if err != nil {
		// fang handles error styling
		os.Exit(1)
	}
}

// Scope returns the DI scope for service registration.
// Use this to register services that commands can access.
func (g *GuardedCommand[T, F]) Scope() do.Injector {
	return g.scope.Injector()
}

// ScopeStruct returns the wrapped Scope struct for advanced operations.
func (g *GuardedCommand[T, F]) ScopeStruct() *Scope {
	return g.scope
}

// Config returns the resolved configuration.
// This is populated after flag parsing.
func (g *GuardedCommand[T, F]) Config() *T {
	return g.config
}

// SetConfig updates the configuration.
// Useful for setting config programmatically before execution.
func (g *GuardedCommand[T, F]) SetConfig(cfg T) {
	g.config = &cfg
}

// RootCommand returns the underlying cobra root command.
// Use this for advanced cobra configuration.
func (g *GuardedCommand[T, F]) RootCommand() *cobra.Command {
	return g.rootCmd
}

// Shutdown gracefully shuts down the CLI application.
// Call this after Execute returns for cleanup.
func (g *GuardedCommand[T, F]) Shutdown(ctx context.Context) error {
	return g.scope.Shutdown(ctx)
}

// HealthCheck runs health checks on all registered services.
func (g *GuardedCommand[T, F]) HealthCheck() error {
	return g.scope.HealthCheck()
}

// Name returns the CLI application name.
func (g *GuardedCommand[T, F]) Name() string {
	return g.name
}

// Short returns the short description.
func (g *GuardedCommand[T, F]) Short() string {
	return g.short
}

// Long returns the long description.
func (g *GuardedCommand[T, F]) Long() string {
	return g.long
}

// SetLong sets the long description.
func (g *GuardedCommand[T, F]) SetLong(long string) {
	g.long = long
	g.rootCmd.Long = long
}

// SetVersion sets the version string.
func (g *GuardedCommand[T, F]) SetVersion(version string) {
	g.rootCmd.Version = version
}

// AddGlobalFlag adds a persistent flag available to all commands.
func (g *GuardedCommand[T, F]) AddGlobalFlag(name, shorthand, defaultValue, help string) {
	g.rootCmd.PersistentFlags().StringP(name, shorthand, defaultValue, help)
}

// AddGlobalBoolFlag adds a persistent boolean flag available to all commands.
func (g *GuardedCommand[T, F]) AddGlobalBoolFlag(
	name, shorthand string,
	defaultValue bool,
	help string,
) {
	g.rootCmd.PersistentFlags().BoolP(name, shorthand, defaultValue, help)
}
