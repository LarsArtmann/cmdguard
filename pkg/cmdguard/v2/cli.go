package v2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"charm.land/fang/v2"
	output "github.com/larsartmann/go-output"
	auditlog "github.com/larsartmann/samber-do-auditlog"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// CLI provides type-safe CLI construction with a single type parameter.
// T is the application config type. Commands can have any flags type.
// This is the recommended API for new code (v2.1+).
type CLI[T any] struct {
	name             string
	short            string
	long             string
	version          string
	defaults         T
	config           *T
	scope            *Scope
	rootCmd          *cobra.Command
	registry         *FlagRegistry
	registeredCmds   map[string]struct{}
	flowCtx          *BranchingFlowContext
	useFang          bool
	fangOpts         []fang.Option
	middleware       []Middleware[T]
	envPrefix        string
	signalHandling   bool
	outputFormat     OutputFormat
	validationMode   ValidationMode
	configValidate   func(*T) error
	postFlagParse    []func(*cobra.Command, *T) error
	cleanupHooks     []func(*cobra.Command, *T, error) error
	cleanupWired     bool
	configFilePaths  []string
	configFileLoader ConfigFileLoader
	helpTransforms   []HelpTransformFunc
	noColorFlag      *bool
	gracefulShutdown bool
	diLogf           func(string, ...any)
	auditLog         *auditlog.Plugin
}

// NewCLI creates a new CLI application with typed config.
// Returns an error if initialization fails (never panics).
// T is the application config type.
func NewCLI[T any](name, short string, defaults T, opts ...CLIOption[T]) (*CLI[T], error) {
	err := validateName(name)
	if err != nil {
		return nil, fmt.Errorf("short=%q: creating CLI %q: %w", short, name, err)
	}

	cli := &CLI[T]{
		name:           name,
		short:          short,
		defaults:       defaults,
		scope:          nil,
		rootCmd:        &cobra.Command{Use: name, Short: short},
		registry:       nil,
		registeredCmds: make(map[string]struct{}),
		useFang:        true,
		noColorFlag:    new(bool),
	}

	for _, opt := range opts {
		opt(cli)
	}

	err = cli.initialize(defaults)
	if err != nil {
		return nil, fmt.Errorf("short=%q, initializing CLI %q: %w", short, name, err)
	}

	return cli, nil
}

func (cli *CLI[T]) initialize(defaults T) error {
	if cli.scope == nil {
		opts := cli.buildInjectorOpts()

		if opts != nil {
			cli.scope = NewScopeWithOpts(cli.name, opts)
		} else {
			cli.scope = NewScope(cli.name)
		}
	}

	cfg := defaults

	err := ProvideValue(cli.scope, &cfg)
	if err != nil {
		return fmt.Errorf("%w: registering config type=%T: %w", ErrServiceRegistration, cfg, err)
	}

	cli.config = &cfg

	registry, err := NewFlagRegistry(cli.config)
	if err != nil {
		return fmt.Errorf(
			"%w: creating flag registry for config=%T: %w",
			ErrServiceRegistration,
			cli.config,
			err,
		)
	}

	err = ProvideValue(cli.scope, registry)
	if err != nil {
		return fmt.Errorf(
			"%w: registering flag registry for %T: %w",
			ErrServiceRegistration,
			defaults,
			err,
		)
	}

	cli.registry = registry

	if setFields, err := cli.loadConfigFileOrSkip(); err != nil {
		return fmt.Errorf("%w: loading config file: %w", ErrConfigFileLoad, err)
	} else if len(setFields) > 0 {
		registry.updateTagDefaultsFromConfig(cli.config, setFields)
	}

	if cli.envPrefix != "" {
		registry.SetEnvPrefix(cli.envPrefix)
	}

	cli.initOutputFlag()
	cli.initNoColorFlag()

	// Silence usage-on-error by default. Raw Cobra prints the full usage block
	// after every command error — the single most reported Cobra footgun. A CLI
	// library that aims to make consumers "use Cobra correctly" must not expose
	// that behaviour by default. Fang already forces this true when it executes;
	// setting it here guarantees the same sane behaviour when fang is disabled.
	// --help is unaffected (SilenceUsage only suppresses error usage).
	cli.rootCmd.SilenceUsage = true

	err = registry.RegisterScopedFlags(cli.rootCmd)
	if err != nil {
		return fmt.Errorf(
			"%w: registering global flags for %T: %w",
			ErrFlagParseFailed,
			defaults,
			err,
		)
	}

	cli.rootCmd.PersistentPreRunE = func(c *cobra.Command, _ []string) error {
		if err := registry.ParseFlags(c, cli.config); err != nil {
			return fmt.Errorf("parsing global flags: %w", err)
		}

		// Store the resolved config in the command context so raw cobra
		// subcommands (added via cli.RootCommand().AddCommand) can access it
		// via ConfigFromContext[T](cmd.Context()) without a parallel
		// context-key system.
		c.SetContext(context.WithValue(c.Context(), configKey, cli.config))

		if cli.configValidate != nil {
			if err := cli.configValidate(cli.config); err != nil {
				return fmt.Errorf("%w: %w", ErrConfigValidation, err)
			}
		}

		for _, fn := range cli.postFlagParse {
			if err := fn(c, cli.config); err != nil {
				return fmt.Errorf("post-flag-parse hook: %w", err)
			}
		}

		return cli.parseOutputFlag(c)
	}

	return nil
}

// buildInjectorOpts merges DI logging and audit log hooks into a single InjectorOpts.
// Returns nil when neither is configured, so the default injector is used.
func (cli *CLI[T]) buildInjectorOpts() *do.InjectorOpts {
	if cli.diLogf == nil && cli.auditLog == nil {
		return nil
	}

	opts := &do.InjectorOpts{}

	if cli.diLogf != nil {
		opts.Logf = cli.diLogf
	}

	if cli.auditLog != nil {
		auditOpts := cli.auditLog.Opts()

		opts.HookBeforeRegistration = append(opts.HookBeforeRegistration, auditOpts.HookBeforeRegistration...)
		opts.HookAfterRegistration = append(opts.HookAfterRegistration, auditOpts.HookAfterRegistration...)
		opts.HookBeforeInvocation = append(opts.HookBeforeInvocation, auditOpts.HookBeforeInvocation...)
		opts.HookAfterInvocation = append(opts.HookAfterInvocation, auditOpts.HookAfterInvocation...)
		opts.HookBeforeShutdown = append(opts.HookBeforeShutdown, auditOpts.HookBeforeShutdown...)
		opts.HookAfterShutdown = append(opts.HookAfterShutdown, auditOpts.HookAfterShutdown...)
	}

	return opts
}

// AddCommand adds a subcommand to the CLI with any flags type.
func AddCommand[T, F any](cli *CLI[T], cmd Command[T, F]) error {
	if _, exists := cli.registeredCmds[cmd.spec.use]; exists {
		return fmt.Errorf("%w: command %q already exists", ErrDuplicateCommand, cmd.spec.use)
	}

	if err := cmd.validate(cli.validationMode); err != nil {
		return fmt.Errorf("validating command %q on CLI %q: %w", cmd.spec.use, cli.name, err)
	}

	cli.registeredCmds[cmd.spec.use] = struct{}{}

	cobraCmd, err := cliToCobraCommand(cli.config, cmd, cli.middleware, cli.envPrefix)
	if err != nil {
		return fmt.Errorf("converting command %q for CLI %q: %w", cmd.spec.use, cli.name, err)
	}

	cli.rootCmd.AddCommand(cobraCmd)

	return nil
}

func (cli *CLI[T]) applyHelpTransforms() {
	for _, fn := range cli.helpTransforms {
		fn(cli.rootCmd)
	}
}

func (cli *CLI[T]) initNoColorFlag() {
	cli.rootCmd.PersistentFlags().BoolVar(cli.noColorFlag, "no-color", false,
		"Disable colored output (also respected via NO_COLOR env var)")
}

// applyNoColorIfSet temporarily sets NO_COLOR=1 around fang execution
// if --no-color was passed. The original value is restored after execution
// to avoid process-wide env mutation.
func (cli *CLI[T]) applyNoColorIfSet() func() {
	if !cli.useFang || cli.noColorFlag == nil || !*cli.noColorFlag {
		return func() {}
	}

	previous := os.Getenv("NO_COLOR")
	_ = os.Setenv("NO_COLOR", "1")

	return func() {
		if previous == "" {
			_ = os.Unsetenv("NO_COLOR")
		} else {
			_ = os.Setenv("NO_COLOR", previous)
		}
	}
}

// applyCleanupHooks wraps every command's RunE — both cmdguard-managed
// commands and raw *cobra.Command subcommands added via
// RootCommand().AddCommand — so that cleanup hooks fire after RunE completes,
// including when RunE returns an error.
//
// Cobra's PostRunE and PersistentPostRunE are NOT called when RunE errors, so
// cleanup that must run on failure (flushing buffers, releasing resources,
// emitting failure telemetry) has nowhere to live except the RunE wrapper
// itself. This wires that wrapper once, uniformly, for the whole tree.
//
// It is a no-op when no cleanup hooks are registered, so CLIs that do not use
// WithCleanup pay zero overhead. The wiring is idempotent so calling Execute
// more than once does not double-wrap.
func (cli *CLI[T]) applyCleanupHooks() {
	if len(cli.cleanupHooks) == 0 || cli.cleanupWired {
		return
	}

	cli.cleanupWired = true

	var wrap func(cmd *cobra.Command)

	wrap = func(cmd *cobra.Command) {
		if cmd.RunE != nil {
			original := cmd.RunE
			cmd.RunE = func(c *cobra.Command, args []string) error {
				runErr := original(c, args)

				cfg, _ := ConfigFromContext[T](c.Context())

				var cleanupErrs []error

				for _, fn := range cli.cleanupHooks {
					if cerr := fn(c, cfg, runErr); cerr != nil {
						cleanupErrs = append(cleanupErrs, fmt.Errorf("cleanup hook: %w", cerr))
					}
				}

				// The original RunE error is never swallowed. When the handler
				// failed it stays the primary error; cleanup failures are joined
				// so they remain reachable via errors.Is without hiding runErr.
				if runErr != nil {
					if len(cleanupErrs) == 0 {
						return runErr
					}

					return errors.Join(append([]error{runErr}, cleanupErrs...)...)
				}

				return errors.Join(cleanupErrs...)
			}
		}

		for _, sub := range cmd.Commands() {
			wrap(sub)
		}
	}

	wrap(cli.rootCmd)
}

// Execute runs the CLI application.
// The context is wrapped with a BranchingFlowContext for command path tracking.
// If WithSignalHandling was set, the context is cancelled on SIGINT/SIGTERM.
// If WithGracefulShutdown was set, DI services are shut down on signal after command completes.
func (cli *CLI[T]) Execute(ctx context.Context) error {
	cli.applyHelpTransforms()

	restoreNoColor := cli.applyNoColorIfSet()
	defer restoreNoColor()

	jsonErrors := cli.outputFormat != "" &&
		(cli.outputFormat == output.FormatJSON || cli.outputFormat == output.FormatJSONL ||
			cli.outputFormat == output.FormatYAML || cli.outputFormat == output.FormatTOML)

	if jsonErrors {
		cli.rootCmd.SilenceErrors = true
		cli.rootCmd.SilenceUsage = true

		if cli.useFang {
			cli.fangOpts = append(cli.fangOpts, fang.WithErrorHandler(func(_ io.Writer, _ fang.Styles, _ error) {}))
		}
	}

	if cli.signalHandling {
		var cancel context.CancelFunc

		ctx, cancel = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		if cli.gracefulShutdown {
			shutdownCtx := context.WithoutCancel(ctx)
			defer func() { _ = cli.scope.Shutdown(shutdownCtx) }()
		}
	}

	if cli.flowCtx == nil {
		cli.flowCtx = NewBranchingFlowContext(ctx)
	}

	flowCtx := WithBranchingFlowContext(ctx, cli.flowCtx)

	cli.applyCleanupHooks()

	var execErr error

	if cli.useFang {
		execErr = fang.Execute(flowCtx, cli.rootCmd, cli.fangOpts...)
	} else {
		execErr = cli.rootCmd.ExecuteContext(flowCtx)
	}

	if execErr != nil {
		// writeFormattedError emits a structured JSON/YAML error to stderr when a
		// machine-readable output format is selected (and fang is silenced above).
		// The error is always returned so ExecuteAndExit can map it to an exit code.
		cli.writeFormattedError(execErr)

		return fmt.Errorf("failed to execute CLI: %w", execErr)
	}

	return nil
}

// ExecuteWithArgs runs the CLI application with specific arguments.
func (cli *CLI[T]) ExecuteWithArgs(ctx context.Context, args []string) error {
	cli.rootCmd.SetArgs(args)

	return cli.Execute(ctx)
}

// ExecuteAndExit runs the CLI and exits with the appropriate exit code.
// If the error implements ExitCoder, its exit code is used; otherwise defaults to 1.
func (cli *CLI[T]) ExecuteAndExit(ctx context.Context) {
	err := cli.Execute(ctx)
	if err != nil {
		os.Exit(ExitCode(err))
	}
}

// validateName checks that the command name is not empty.
func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required, name=%q", ErrMissingName, name)
	}

	return nil
}
