package v3

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

// cliSpec holds all non-generic CLI configuration. Typed hooks (config
// validation, middleware, post-flag-parse, cleanup) are stored behind sealed
// interfaces, preserving type safety without genericizing the option type.
type cliSpec struct {
	name             string
	short            string
	long             string
	version          string
	scope            *Scope
	useFang          bool
	fangOpts         []fang.Option
	envPrefix        string
	signalHandling   bool
	silenceErrors    bool
	silenceUsage     bool
	outputFormat     OutputFormat
	validationMode   ValidationMode
	configValidate   configValidator
	middleware       middlewareList
	postFlagParse    postFlagParseList
	cleanupHooks     cleanupHookList
	configFilePaths  []string
	configLoader     ConfigFileLoader
	helpTransforms   []HelpTransformFunc
	groups           []cobraGroup
	gracefulShutdown bool
	diLogf           func(string, ...any)
	auditLog         *auditlog.Plugin
	pluginErr        error
}

type cobraGroup struct {
	id    string
	title string
}

// Sealed interfaces for typed CLI hooks — the unexported methods prevent
// external implementations, ensuring type safety through the non-generic
// CLIOption.

type configValidator interface {
	isConfigValidator()
}

type typedConfigValidator[T any] struct {
	fn func(*T) error
}

func (*typedConfigValidator[T]) isConfigValidator() {}

type middlewareList interface {
	isMiddlewareList()
}

type typedMiddlewareList[T any] struct {
	mws []Middleware[T]
}

func (*typedMiddlewareList[T]) isMiddlewareList() {}

type postFlagParseList interface {
	isPostFlagParseList()
}

type typedPostFlagParseList[T any] struct {
	fns []func(*cobra.Command, *T) error
}

func (*typedPostFlagParseList[T]) isPostFlagParseList() {}

type cleanupHookList interface {
	isCleanupHookList()
}

type typedCleanupHookList[T any] struct {
	fns []func(*cobra.Command, *T, error) error
}

func (*typedCleanupHookList[T]) isCleanupHookList() {}

// CLIOption configures a CLI. Non-generic — no type parameters needed on
// most options. Generic-returning options (WithConfigValidation, WithMiddleware,
// WithPostFlagParse, WithCleanup) use sealed interfaces internally.
type CLIOption func(*cliSpec)

// CLI provides type-safe CLI construction with a single type parameter.
// T is the application config type. Commands can have any flags type.
type CLI[T any] struct {
	spec           cliSpec
	defaults       T
	config         *T
	rootCmd        *cobra.Command
	registry       *FlagRegistry
	registeredCmds map[string]struct{}
	flowCtx        *BranchingFlowContext
	noColorFlag    *bool
	cleanupWired   bool
}

// NewCLI creates a new CLI application with typed config.
// Returns an error if initialization fails (never panics).
// T is the application config type.
//
// Type parameter T must be specified explicitly; CLIOptions do not need
// type parameters:
//
//	cli, err := cmdguard.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
//	    cmdguard.WithCLIVersion("1.0.0"),
//	    cmdguard.WithEnvPrefix("MYAPP_"),
//	    cmdguard.WithStrictValidation(),
//	)
func NewCLI[T any](name, short string, defaults T, opts ...CLIOption) (*CLI[T], error) {
	err := validateName(name)
	if err != nil {
		return nil, fmt.Errorf("short=%q: creating CLI %q: %w", short, name, err)
	}

	spec := cliSpec{
		name:         name,
		short:        short,
		useFang:      true,
		silenceUsage: true,
	}

	for _, opt := range opts {
		opt(&spec)
	}

	if spec.pluginErr != nil {
		return nil, fmt.Errorf("short=%q, creating CLI %q: %w", short, name, spec.pluginErr)
	}

	cli := &CLI[T]{
		spec:           spec,
		defaults:       defaults,
		rootCmd:        &cobra.Command{Use: name, Short: short},
		registeredCmds: make(map[string]struct{}),
		noColorFlag:    new(bool),
	}

	// Apply command groups
	for _, g := range spec.groups {
		cli.rootCmd.AddGroup(&cobra.Group{ID: g.id, Title: g.title})
	}

	// Apply version and long description
	if spec.version != "" {
		cli.rootCmd.Version = spec.version
	}

	if spec.long != "" {
		cli.rootCmd.Long = spec.long
	}

	// Apply silence settings
	if spec.silenceErrors {
		cli.rootCmd.SilenceErrors = true
	}

	cli.rootCmd.SilenceUsage = spec.silenceUsage

	err = cli.initialize(defaults)
	if err != nil {
		return nil, fmt.Errorf("short=%q, initializing CLI %q: %w", short, name, err)
	}

	return cli, nil
}

func (cli *CLI[T]) initialize(defaults T) error {
	cli.ensureScope()

	cfg := defaults

	err := ProvideValue(cli.spec.scope, &cfg)
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

	err = ProvideValue(cli.spec.scope, registry)
	if err != nil {
		return fmt.Errorf(
			"%w: registering flag registry for %T: %w",
			ErrServiceRegistration,
			defaults,
			err,
		)
	}

	cli.registry = registry

	setFields, err := cli.loadConfigFileOrSkip()
	if err != nil {
		return fmt.Errorf("%w: loading config file: %w", ErrConfigFileLoad, err)
	} else if len(setFields) > 0 {
		registry.updateTagDefaultsFromConfig(cli.config, setFields)
	}

	if cli.spec.envPrefix != "" {
		registry.SetEnvPrefix(cli.spec.envPrefix)
	}

	cli.initOutputFlag()
	cli.initNoColorFlag()

	cli.rootCmd.SilenceUsage = cli.spec.silenceUsage

	err = registry.RegisterScopedFlags(cli.rootCmd)
	if err != nil {
		return fmt.Errorf(
			"%w: registering global flags for %T: %w",
			ErrFlagParseFailed,
			defaults,
			err,
		)
	}

	cli.setupPersistentPreRun()

	return nil
}

func (cli *CLI[T]) ensureScope() {
	if cli.spec.scope != nil {
		return
	}

	opts := cli.buildInjectorOpts()

	if opts != nil {
		cli.spec.scope = NewScopeWithOpts(cli.spec.name, opts)
	} else {
		cli.spec.scope = NewScope(cli.spec.name)
	}
}

func (cli *CLI[T]) setupPersistentPreRun() {
	registry := cli.registry
	s := &cli.spec

	cli.rootCmd.PersistentPreRunE = func(c *cobra.Command, _ []string) error {
		err := registry.ParseFlags(c, cli.config)
		if err != nil {
			return fmt.Errorf("parsing global flags: %w", err)
		}

		c.SetContext(context.WithValue(c.Context(), configKey, cli.config))

		if cv, ok := s.configValidate.(*typedConfigValidator[T]); ok && cv != nil {
			err := cv.fn(cli.config)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrConfigValidation, err)
			}
		}

		if pfpl, ok := s.postFlagParse.(*typedPostFlagParseList[T]); ok && pfpl != nil {
			for _, fn := range pfpl.fns {
				err := fn(c, cli.config)
				if err != nil {
					return fmt.Errorf("post-flag-parse hook: %w", err)
				}
			}
		}

		return cli.parseOutputFlag(c)
	}
}

// buildInjectorOpts merges DI logging and audit log hooks into a single InjectorOpts.
func (cli *CLI[T]) buildInjectorOpts() *do.InjectorOpts {
	s := &cli.spec
	if s.diLogf == nil && s.auditLog == nil {
		return nil
	}

	opts := &do.InjectorOpts{}

	if s.diLogf != nil {
		opts.Logf = s.diLogf
	}

	if s.auditLog != nil {
		auditOpts := s.auditLog.Opts()

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

	err := cmd.validate(cli.spec.validationMode)
	if err != nil {
		return fmt.Errorf("validating command %q on CLI %q: %w", cmd.spec.use, cli.spec.name, err)
	}

	cli.registeredCmds[cmd.spec.use] = struct{}{}

	cobraCmd, err := cliToCobraCommand(cli.config, cmd, cli.extractMiddleware(), cli.spec.envPrefix)
	if err != nil {
		return fmt.Errorf("converting command %q for CLI %q: %w", cmd.spec.use, cli.spec.name, err)
	}

	// Propagate CLI-level silence-usage to subcommands.
	cobraCmd.SilenceUsage = cli.spec.silenceUsage

	cli.rootCmd.AddCommand(cobraCmd)

	return nil
}

// extractMiddleware safely extracts the typed middleware from the sealed interface.
func (cli *CLI[T]) extractMiddleware() []Middleware[T] {
	if ml, ok := cli.spec.middleware.(*typedMiddlewareList[T]); ok {
		return ml.mws
	}

	return nil
}

// extractCleanupHooks safely extracts typed cleanup hooks from the sealed interface.
func (cli *CLI[T]) extractCleanupHooks() []func(*cobra.Command, *T, error) error {
	if cl, ok := cli.spec.cleanupHooks.(*typedCleanupHookList[T]); ok {
		return cl.fns
	}

	return nil
}

func (cli *CLI[T]) applyHelpTransforms() {
	for _, fn := range cli.spec.helpTransforms {
		fn(cli.rootCmd)
	}
}

func (cli *CLI[T]) initNoColorFlag() {
	cli.rootCmd.PersistentFlags().BoolVar(cli.noColorFlag, "no-color", false,
		"Disable colored output (also respected via NO_COLOR env var)")
}

func (cli *CLI[T]) applyNoColorIfSet() func() {
	if !cli.spec.useFang || cli.noColorFlag == nil || !*cli.noColorFlag {
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

func (cli *CLI[T]) applyCleanupHooks() {
	hooks := cli.extractCleanupHooks()
	if len(hooks) == 0 || cli.cleanupWired {
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

				for _, fn := range hooks {
					cerr := fn(c, cfg, runErr)
					if cerr != nil {
						cleanupErrs = append(cleanupErrs, fmt.Errorf("cleanup hook: %w", cerr))
					}
				}

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
func (cli *CLI[T]) Execute(ctx context.Context) error {
	cli.applyHelpTransforms()

	restoreNoColor := cli.applyNoColorIfSet()
	defer restoreNoColor()

	jsonErrors := cli.spec.outputFormat != "" &&
		(cli.spec.outputFormat == output.FormatJSON || cli.spec.outputFormat == output.FormatJSONL ||
			cli.spec.outputFormat == output.FormatYAML || cli.spec.outputFormat == output.FormatTOML)

	if jsonErrors {
		cli.rootCmd.SilenceErrors = true
		cli.rootCmd.SilenceUsage = true

		if cli.spec.useFang {
			cli.spec.fangOpts = append(
				cli.spec.fangOpts,
				fang.WithErrorHandler(func(_ io.Writer, _ fang.Styles, _ error) {}),
			)
		}
	}

	if cli.spec.signalHandling {
		var cancel context.CancelFunc

		ctx, cancel = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		if cli.spec.gracefulShutdown {
			shutdownCtx := context.WithoutCancel(ctx)
			defer func() { _ = cli.spec.scope.Shutdown(shutdownCtx) }()
		}
	}

	if cli.flowCtx == nil {
		cli.flowCtx = NewBranchingFlowContext(ctx)
	}

	// Make flow context available to commands via context
	flowCtx := WithBranchingFlowContext(ctx, cli.flowCtx)

	if cli.rootCmd.Long == "" && cli.spec.long != "" {
		cli.rootCmd.Long = cli.spec.long
	}

	cli.applyCleanupHooks()

	return cli.executeWithCobra(flowCtx)
}

// executeWithCobra runs the cobra command, dispatching through fang when enabled.
func (cli *CLI[T]) executeWithCobra(flowCtx context.Context) error {
	var execErr error

	if cli.spec.useFang {
		execErr = fang.Execute(flowCtx, cli.rootCmd, cli.spec.fangOpts...)
	} else {
		execErr = cli.rootCmd.ExecuteContext(flowCtx)
	}

	if execErr != nil {
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
