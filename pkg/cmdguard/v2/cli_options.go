package v2

import (
	"io"

	"charm.land/fang/v2"
	"charm.land/lipgloss/v2"
	auditlog "github.com/larsartmann/samber-do-auditlog"
	"github.com/spf13/cobra"
)

// CLIOption is a functional option for configuring a CLI.
type CLIOption[T any] func(*CLI[T])

// WithCLIVersion sets the version string.
// When fang is enabled, the version is automatically passed to fang.WithVersion
// for styled version output alongside cmdguard's own version subcommand.
func WithCLIVersion[T any](version string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.setVersion(version)
		cli.fangOpts = append(cli.fangOpts, fang.WithVersion(version))
	}
}

// WithCLICommit sets the git commit hash appended to the version string.
// When fang is enabled, the commit is automatically passed to fang.WithCommit.
func WithCLICommit[T any](commit string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.fangOpts = append(cli.fangOpts, fang.WithCommit(commit))
	}
}

// WithCLILong sets the long description.
func WithCLILong[T any](long string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.setLong(long)
	}
}

// WithCLIScope sets a custom DI scope.
func WithCLIScope[T any](scope *Scope) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.scope = scope
	}
}

// WithSilenceErrors suppresses automatic error printing from cobra.
func WithSilenceErrors[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.rootCmd.SilenceErrors = true
	}
}

// WithSilenceUsage suppresses automatic usage printing on error.
func WithSilenceUsage[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.rootCmd.SilenceUsage = true
	}
}

// WithFang enables or disables fang-based styled output.
// When disabled, falls back to cobra's default plain text output.
func WithFang[T any](enabled bool) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.useFang = enabled
	}
}

// WithFangOptions sets fang options for the CLI's Execute method.
func WithFangOptions[T any](opts ...fang.Option) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.fangOpts = append(cli.fangOpts, opts...)
	}
}

// WithFangErrorHandler sets a custom error display function for fang's styled output.
// The function receives the writer, fang styles, and the error to display.
// Only effective when fang is enabled (default).
func WithFangErrorHandler[T any](handler func(w io.Writer, styles fang.Styles, err error)) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.fangOpts = append(cli.fangOpts, fang.WithErrorHandler(handler))
	}
}

// WithFangColorScheme sets a custom color scheme for fang's styled help and error output.
// The function receives a lipgloss.LightDarkFunc and returns a fang.ColorScheme.
// Only effective when fang is enabled (default).
func WithFangColorScheme[T any](cs func(lightDark lipgloss.LightDarkFunc) fang.ColorScheme) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.fangOpts = append(cli.fangOpts, fang.WithColorSchemeFunc(cs))
	}
}

// WithMiddleware adds middleware that wraps every command handler.
// Middleware are applied in order: first wraps the second, etc.
func WithMiddleware[T any](mw ...Middleware[T]) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.middleware = append(cli.middleware, mw...)
	}
}

// WithGroup registers a command group on the root command.
// Groups organize commands in help output under titled sections.
// Use the Group field on Command to assign a command to a registered group.
func WithGroup[T any](groupID, title string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.rootCmd.AddGroup(&cobra.Group{ID: groupID, Title: title})
	}
}

// WithEnvPrefix sets a prefix for environment variable lookups.
// When set, env tags are prefixed: prefix "APP_" + env tag "PORT" → "APP_PORT".
func WithEnvPrefix[T any](prefix string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.envPrefix = prefix
	}
}

// WithSignalHandling adds automatic context cancellation on SIGINT/SIGTERM.
// When a signal is received, the context passed to handlers is cancelled.
// This does NOT trigger DI service shutdown — use WithGracefulShutdown for that.
func WithSignalHandling[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.signalHandling = true
	}
}

// WithConfigValidation adds a validation function that runs after root flag parsing
// but before any command handler. Use this to validate the full config struct
// (e.g., cross-field validation, business rules).
//
// The validator receives a pointer to the resolved config. Return an error to
// stop execution before any command runs.
func WithConfigValidation[T any](validate func(*T) error) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.configValidate = validate
	}
}

// WithStrictValidation enables strict command validation:
//   - All commands must have a short description
func WithStrictValidation[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.validationMode = Strict
	}
}

// WithGracefulShutdown enables graceful DI shutdown on SIGINT/SIGTERM.
// When a signal is received, all services implementing do.ShutdownerWithError
// or do.ShutdownerWithContextAndError are shut down in reverse invocation order.
// This also enables signal-based context cancellation (implies WithSignalHandling).
func WithGracefulShutdown[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.signalHandling = true
		cli.gracefulShutdown = true
	}
}

// WithDILogging enables internal logging for the DI container.
// The provided function receives formatted log messages from samber/do
// for service registration, invocation, and lifecycle events.
func WithDILogging[T any](logf func(format string, args ...any)) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.diLogf = logf
	}
}

// WithDraconianValidation enables draconian command validation:
//   - All commands must have a short description
//   - All leaf commands must have an example
func WithDraconianValidation[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.validationMode = Draconian
	}
}

// WithAuditLog enables DI audit logging via samber-do-auditlog.
// The plugin captures service registration, invocation, shutdown, and health check events.
// Use cli.AuditLog() to access the plugin for reports, exports, and HTML visualization.
//
// When Config.Enabled is false (the zero value), the plugin checks DO_AUDITLOG_ENABLED.
// Set it to "true", "1", or "yes" to enable audit logging without changing code.
func WithAuditLog[T any](plugin *auditlog.Plugin) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.auditLog = plugin
	}
}
