package v4

import (
	"io"

	"charm.land/fang/v2"
	"charm.land/lipgloss/v2"
	auditlog "github.com/larsartmann/samber-do-auditlog"
	"github.com/spf13/cobra"
)

// HelpTransformFunc transforms a cobra command's help text before execution.
// Optional modules (e.g. cmdguard/glamour) use this hook to render markdown,
// without the core module importing the rendering library.
type HelpTransformFunc func(cmd *cobra.Command)

// --- Non-generic CLI options (operate on cliSpec) ---

// WithHelpTransform registers a function that transforms command help text
// before the CLI executes. Multiple transforms run in registration order.
func WithHelpTransform(fn HelpTransformFunc) CLIOption {
	return func(s *cliSpec) {
		s.helpTransforms = append(s.helpTransforms, fn)
	}
}

// WithCLIVersion sets the version string.
func WithCLIVersion(version string) CLIOption {
	return func(s *cliSpec) {
		s.version = version
		s.fangOpts = append(s.fangOpts, fang.WithVersion(version))
	}
}

// WithCLICommit sets the git commit hash.
func WithCLICommit(commit string) CLIOption {
	return func(s *cliSpec) {
		s.fangOpts = append(s.fangOpts, fang.WithCommit(commit))
	}
}

// WithCLILong sets the long description.
func WithCLILong(long string) CLIOption {
	return func(s *cliSpec) {
		s.long = long
	}
}

// WithCLIScope sets a custom DI scope.
func WithCLIScope(scope *Scope) CLIOption {
	return func(s *cliSpec) {
		s.scope = scope
	}
}

// WithSilenceErrors suppresses automatic error printing from cobra.
func WithSilenceErrors() CLIOption {
	return func(s *cliSpec) {
		s.silenceErrors = true
	}
}

// WithSilenceUsage suppresses automatic usage printing on error for the root
// command and all subcommands. This is enabled by default to avoid cobra's
// notorious footgun of dumping full usage text on every error. This option
// exists to make the intent explicit.
func WithSilenceUsage() CLIOption {
	return func(s *cliSpec) {
		s.silenceUsage = true
	}
}

// WithoutSilenceUsage re-enables cobra's default behavior of printing usage
// text on error. Use this when you want usage-on-error for debugging or when
// your users expect to see usage hints when commands fail.
func WithoutSilenceUsage() CLIOption {
	return func(s *cliSpec) {
		s.silenceUsage = false
	}
}

// WithFang enables or disables fang-based styled output.
func WithFang(enabled bool) CLIOption {
	return func(s *cliSpec) {
		s.useFang = enabled
	}
}

// WithFangOptions sets fang options for the CLI's Execute method.
func WithFangOptions(opts ...fang.Option) CLIOption {
	return func(s *cliSpec) {
		s.fangOpts = append(s.fangOpts, opts...)
	}
}

// WithFangErrorHandler sets a custom error display function for fang.
func WithFangErrorHandler(handler func(w io.Writer, styles fang.Styles, err error)) CLIOption {
	return func(s *cliSpec) {
		s.fangOpts = append(s.fangOpts, fang.WithErrorHandler(handler))
	}
}

// WithFangColorScheme sets a custom color scheme for fang.
func WithFangColorScheme(cs func(lightDark lipgloss.LightDarkFunc) fang.ColorScheme) CLIOption {
	return func(s *cliSpec) {
		s.fangOpts = append(s.fangOpts, fang.WithColorSchemeFunc(cs))
	}
}

// WithGroup registers a command group on the root command.
// Groups organize commands in help output under titled sections.
func WithGroup(groupID, title string) CLIOption {
	return func(s *cliSpec) {
		s.groups = append(s.groups, cobraGroup{id: groupID, title: title})
	}
}

// WithEnvPrefix sets a prefix for environment variable lookups.
func WithEnvPrefix(prefix string) CLIOption {
	return func(s *cliSpec) {
		s.envPrefix = prefix
	}
}

// WithSignalHandling adds automatic context cancellation on SIGINT/SIGTERM.
func WithSignalHandling() CLIOption {
	return func(s *cliSpec) {
		s.signalHandling = true
	}
}

// WithStrictValidation enables strict command validation.
func WithStrictValidation() CLIOption {
	return func(s *cliSpec) {
		s.validationMode = Strict
	}
}

// WithGracefulShutdown enables graceful DI shutdown on SIGINT/SIGTERM.
func WithGracefulShutdown() CLIOption {
	return func(s *cliSpec) {
		s.signalHandling = true
		s.gracefulShutdown = true
	}
}

// WithDILogging enables internal logging for the DI container.
func WithDILogging(logf func(format string, args ...any)) CLIOption {
	return func(s *cliSpec) {
		s.diLogf = logf
	}
}

// WithOnError registers a callback invoked when Execute returns a non-nil error,
// before the error is returned to the caller. Use this for structured logging
// (e.g. slog for journald/Loki), metrics, or audit trails at the CLI boundary.
//
// The callback receives the raw execution error (before the "failed to execute
// CLI" wrapper). It is called at most once per Execute call.
//
//	cli, _ := v4.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
//	    v4.WithOnError(func(err error) {
//	        slog.Error("command failed", "err", err)
//	    }),
//	)
func WithOnError(fn func(error)) CLIOption {
	return func(s *cliSpec) {
		s.onError = fn
	}
}

// WithDraconianValidation enables draconian command validation.
func WithDraconianValidation() CLIOption {
	return func(s *cliSpec) {
		s.validationMode = Draconian
	}
}

// WithAuditLog enables DI audit logging via samber-do-auditlog.
func WithAuditLog(plugin *auditlog.Plugin) CLIOption {
	return func(s *cliSpec) {
		s.auditLog = plugin
	}
}

// WithOutputFormat sets the default output format for structured output.
func WithOutputFormat(defaultFormat OutputFormat) CLIOption {
	return func(s *cliSpec) {
		s.outputFormat = defaultFormat
	}
}

// WithConfigFile adds config file loading with the given search paths.
// Paths are tried in order; the first existing file wins.
// Environment variables and ~ are expanded in paths.
// Supports JSON, YAML, and TOML formats with automatic detection by file extension.
func WithConfigFile(paths ...string) CLIOption {
	return func(s *cliSpec) {
		s.configFilePaths = paths
		s.configLoader = NewKoanfLoader(paths...)
	}
}

// WithConfigFileLoader sets a custom config file loader and paths.
func WithConfigFileLoader(loader ConfigFileLoader, paths ...string) CLIOption {
	return func(s *cliSpec) {
		s.configLoader = loader
		s.configFilePaths = paths
	}
}

// --- Generic CLI options (return non-generic CLIOption via sealed interface) ---

// WithConfigValidation adds a validation function that runs after root flag
// parsing but before any command handler.
func WithConfigValidation[T any](validate func(*T) error) CLIOption {
	return func(s *cliSpec) {
		s.configValidate = &typedConfigValidator[T]{fn: validate}
	}
}

// WithMiddleware adds middleware that wraps every command handler.
func WithMiddleware[T any](mw ...Middleware[T]) CLIOption {
	return func(s *cliSpec) {
		if existing, ok := s.middleware.(*typedMiddlewareList[T]); ok {
			existing.mws = append(existing.mws, mw...)
		} else {
			s.middleware = &typedMiddlewareList[T]{mws: mw}
		}
	}
}

// WithPostFlagParse adds hooks that run after flag parsing and config
// validation but before any command handler.
func WithPostFlagParse[T any](fns ...func(cmd *cobra.Command, cfg *T) error) CLIOption {
	return func(s *cliSpec) {
		if existing, ok := s.postFlagParse.(*typedPostFlagParseList[T]); ok {
			existing.fns = append(existing.fns, fns...)
		} else {
			s.postFlagParse = &typedPostFlagParseList[T]{fns: fns}
		}
	}
}

// WithCleanup registers hooks that run after a command's RunE completes,
// including when RunE returns an error.
func WithCleanup[T any](fns ...func(cmd *cobra.Command, cfg *T, runErr error) error) CLIOption {
	return func(s *cliSpec) {
		if existing, ok := s.cleanupHooks.(*typedCleanupHookList[T]); ok {
			existing.fns = append(existing.fns, fns...)
		} else {
			s.cleanupHooks = &typedCleanupHookList[T]{fns: fns}
		}
	}
}
