package v2

import (
	"context"

	auditlog "github.com/larsartmann/samber-do-auditlog"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// Scope returns the DI scope for service registration.
func (cli *CLI[T]) Scope() *Scope {
	return cli.scope
}

// Injector returns the underlying DI injector for direct samber/do/v2 operations.
func (cli *CLI[T]) Injector() do.Injector {
	return cli.scope.Injector()
}

// Config returns the resolved configuration.
func (cli *CLI[T]) Config() *T {
	return cli.config
}

// SetConfig updates the configuration.
//
// WARNING: This replaces the config pointer but does NOT reinitialize the
// FlagRegistry. Flags parsed after this call will still write to the OLD
// config struct. Only call this before Execute or when you don't use the
// typed flag system for this config.
func (cli *CLI[T]) SetConfig(cfg T) {
	cli.config = &cfg
}

// RootCommand returns the underlying cobra root command.
func (cli *CLI[T]) RootCommand() *cobra.Command {
	return cli.rootCmd
}

// Shutdown gracefully shuts down the CLI application.
func (cli *CLI[T]) Shutdown(ctx context.Context) error {
	return cli.scope.Shutdown(ctx)
}

// HealthCheck runs health checks on all registered services.
func (cli *CLI[T]) HealthCheck() error {
	return cli.scope.HealthCheck()
}

// HealthCheckResults runs health checks and returns per-service results.
func (cli *CLI[T]) HealthCheckResults() map[string]error {
	return cli.scope.HealthCheckResults()
}

// HealthCheckWithContext runs health checks with context.
func (cli *CLI[T]) HealthCheckWithContext(ctx context.Context) error {
	return cli.scope.HealthCheckWithContext(ctx)
}

// HealthCheckResultsWithContext runs health checks with context and returns per-service results.
func (cli *CLI[T]) HealthCheckResultsWithContext(ctx context.Context) map[string]error {
	return cli.scope.HealthCheckResultsWithContext(ctx)
}

// Name returns the CLI application name.
func (cli *CLI[T]) Name() string {
	return cli.name
}

// Short returns the short description.
func (cli *CLI[T]) Short() string {
	return cli.short
}

// Long returns the long description.
func (cli *CLI[T]) Long() string {
	return cli.long
}

// SetLong sets the long description.
func (cli *CLI[T]) SetLong(long string) {
	cli.setLong(long)
}

func (cli *CLI[T]) setLong(long string) {
	cli.long = long
	cli.rootCmd.Long = long
}

// SetVersion sets the version string.
func (cli *CLI[T]) SetVersion(version string) {
	cli.setVersion(version)
}

func (cli *CLI[T]) setVersion(version string) {
	cli.version = version
	cli.rootCmd.Version = version
}

// FlowContext returns the branching flow context for command path tracking.
// This is nil until Execute is called.
func (cli *CLI[T]) FlowContext() *BranchingFlowContext {
	return cli.flowCtx
}

// AddGlobalFlag adds a persistent flag available to all commands.
func (cli *CLI[T]) AddGlobalFlag(name, shorthand, defaultValue, help string) {
	cli.rootCmd.PersistentFlags().StringP(name, shorthand, defaultValue, help)
}

// AddGlobalBoolFlag adds a persistent boolean flag available to all commands.
func (cli *CLI[T]) AddGlobalBoolFlag(name, shorthand string, defaultValue bool, help string) {
	cli.rootCmd.PersistentFlags().BoolP(name, shorthand, defaultValue, help)
}

// NoColor returns true if --no-color was explicitly passed by the user.
func (cli *CLI[T]) NoColor() bool {
	if cli.noColorFlag == nil {
		return false
	}

	return *cli.noColorFlag
}

// AuditLog returns the audit log plugin, or nil if audit logging is not enabled.
// Use this to access reports, exports, and HTML visualization after commands execute.
func (cli *CLI[T]) AuditLog() *auditlog.Plugin {
	return cli.auditLog
}

// AuditLogReport returns a consolidated audit report snapshot.
// Returns nil if audit logging is not enabled.
func (cli *CLI[T]) AuditLogReport() *auditlog.Report {
	if cli.auditLog == nil {
		return nil
	}

	report := cli.auditLog.Report()

	return &report
}

// RecordAuditHealthCheck runs health checks on all DI services and records
// the results as audit events. Returns the per-service results map.
// Returns nil if audit logging is not enabled.
func (cli *CLI[T]) RecordAuditHealthCheck(ctx context.Context) map[string]error {
	if cli.auditLog == nil {
		return nil
	}

	return cli.auditLog.RecordHealthCheckWithContext(ctx, cli.scope.Injector())
}
