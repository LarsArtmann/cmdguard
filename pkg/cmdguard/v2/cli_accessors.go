package v2

import (
	"context"

	auditlog "github.com/larsartmann/samber-do-auditlog"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// Scope returns the DI scope for service registration.
func (cli *CLI[T]) Scope() *Scope {
	return cli.spec.scope
}

// Injector returns the underlying DI injector for direct samber/do/v2 operations.
func (cli *CLI[T]) Injector() do.Injector {
	return cli.spec.scope.Injector()
}

// Config returns the resolved configuration.
func (cli *CLI[T]) Config() *T {
	return cli.config
}

// SetConfig updates the configuration.
func (cli *CLI[T]) SetConfig(cfg T) {
	cli.config = &cfg
}

// RootCommand returns the underlying cobra root command.
func (cli *CLI[T]) RootCommand() *cobra.Command {
	return cli.rootCmd
}

// Shutdown gracefully shuts down the CLI application.
func (cli *CLI[T]) Shutdown(ctx context.Context) error {
	return cli.spec.scope.Shutdown(ctx)
}

// HealthCheck runs health checks on all registered services.
func (cli *CLI[T]) HealthCheck() error {
	return cli.spec.scope.HealthCheck()
}

// HealthCheckResults runs health checks and returns per-service results.
func (cli *CLI[T]) HealthCheckResults() map[string]error {
	return cli.spec.scope.HealthCheckResults()
}

// HealthCheckWithContext runs health checks with context.
func (cli *CLI[T]) HealthCheckWithContext(ctx context.Context) error {
	return cli.spec.scope.HealthCheckWithContext(ctx)
}

// HealthCheckResultsWithContext runs health checks with context and returns per-service results.
func (cli *CLI[T]) HealthCheckResultsWithContext(ctx context.Context) map[string]error {
	return cli.spec.scope.HealthCheckResultsWithContext(ctx)
}

// Name returns the CLI application name.
func (cli *CLI[T]) Name() string {
	return cli.spec.name
}

// Short returns the short description.
func (cli *CLI[T]) Short() string {
	return cli.spec.short
}

// Long returns the long description.
func (cli *CLI[T]) Long() string {
	return cli.spec.long
}

// SetLong sets the long description.
func (cli *CLI[T]) SetLong(long string) {
	cli.spec.long = long
	cli.rootCmd.Long = long
}

// SetVersion sets the version string.
func (cli *CLI[T]) SetVersion(version string) {
	cli.spec.version = version
	cli.rootCmd.Version = version
}

// FlowContext returns the branching flow context for command path tracking.
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

// RegisterLocalCommandFlags registers the CLI's local-scoped flags on the given subcommand.
func (cli *CLI[T]) RegisterLocalCommandFlags(cmd *cobra.Command) error {
	if cli.registry == nil {
		return nil
	}

	return cli.registry.RegisterLocalFlags(cmd)
}

// NoColor returns true if --no-color was explicitly passed by the user.
func (cli *CLI[T]) NoColor() bool {
	if cli.noColorFlag == nil {
		return false
	}

	return *cli.noColorFlag
}

// AuditLog returns the audit log plugin, or nil if audit logging is not enabled.
func (cli *CLI[T]) AuditLog() *auditlog.Plugin {
	return cli.spec.auditLog
}

// AuditLogReport returns a consolidated audit report snapshot.
func (cli *CLI[T]) AuditLogReport() *auditlog.Report {
	if cli.spec.auditLog == nil {
		return nil
	}

	report := cli.spec.auditLog.Report()

	return &report
}

// RecordAuditHealthCheck runs health checks on all DI services and records
// the results as audit events.
func (cli *CLI[T]) RecordAuditHealthCheck(ctx context.Context) map[string]error {
	if cli.spec.auditLog == nil {
		return nil
	}

	return cli.spec.auditLog.RecordHealthCheckWithContext(ctx, cli.spec.scope.Injector())
}
