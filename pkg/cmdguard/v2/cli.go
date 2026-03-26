package v2

import (
	"context"
	"fmt"

	"charm.land/fang/v2"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// CLI provides type-safe CLI construction with a single type parameter.
// T is the application config type. Commands can have any flags type.
// This is the recommended API for new code (v2.1+).
type CLI[T any] struct {
	name           string
	short          string
	long           string
	version        string
	defaults       T
	config         *T
	scope          *Scope
	rootCmd        *cobra.Command
	registry       *FlagRegistry
	registeredCmds map[string]bool
}

// CLIOption is a functional option for configuring a CLI.
type CLIOption[T any] func(*CLI[T])

// NewCLI creates a new CLI application with typed config.
// Returns an error if initialization fails (never panics).
// T is the application config type.
func NewCLI[T any](name, short string, defaults T, opts ...CLIOption[T]) (*CLI[T], error) {
	err := validateName(name)
	if err != nil {
		return nil, fmt.Errorf("creating CLI %q: %w", name, err)
	}

	cli := &CLI[T]{
		name:           name,
		short:          short,
		defaults:       defaults,
		scope:          nil,
		rootCmd:        &cobra.Command{Use: name, Short: short},
		registry:       nil,
		registeredCmds: make(map[string]bool),
	}

	for _, opt := range opts {
		opt(cli)
	}

	err = cli.initialize(defaults)
	if err != nil {
		return nil, fmt.Errorf("initializing CLI %q: %w", name, err)
	}

	return cli, nil
}

func (cli *CLI[T]) initialize(defaults T) error {
	cli.scope = NewScope(cli.name)

	cfg := defaults

	err := ProvideValue(cli.scope, &cfg)
	if err != nil {
		return fmt.Errorf("failed to register config type=%T: %w", cfg, err)
	}

	cli.config = &cfg

	registry, err := NewFlagRegistry(cli.config)
	if err != nil {
		return fmt.Errorf("failed to create flag registry: config=%T: %w", cli.config, err)
	}

	err = ProvideValue(cli.scope, registry)
	if err != nil {
		return fmt.Errorf("registering flag registry for %T: %w", defaults, err)
	}

	cli.registry = registry

	err = registry.RegisterFlags(cli.rootCmd)
	if err != nil {
		return fmt.Errorf("registering global flags for %T: %w", defaults, err)
	}

	cli.rootCmd.PersistentPreRunE = func(c *cobra.Command, _ []string) error {
		return registry.ParseFlags(c, cli.config)
	}

	return nil
}

// WithCLIVersion sets the version string.
func WithCLIVersion[T any](version string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.version = version
		cli.rootCmd.Version = version
	}
}

// WithCLILong sets the long description.
func WithCLILong[T any](long string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.long = long
		cli.rootCmd.Long = long
	}
}

// WithCLIScope sets a custom DI scope.
func WithCLIScope[T any](scope *Scope) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.scope = scope
	}
}

// AddCommand adds a subcommand to the CLI with any flags type.
func AddCommand[T, F any](cli *CLI[T], cmd Command[T, F]) error {
	if cli.registeredCmds[cmd.Use] {
		return fmt.Errorf("%w: command %q already exists", ErrDuplicateCommand, cmd.Use)
	}

	err := cmd.Validate()
	if err != nil {
		return fmt.Errorf("validating command %q on CLI %q: %w", cmd.Use, cli.name, err)
	}

	cli.registeredCmds[cmd.Use] = true

	cobraCmd, err := cliToCobraCommand(cli.config, cmd)
	if err != nil {
		return fmt.Errorf("converting command %q for CLI %q: %w", cmd.Use, cli.name, err)
	}

	cli.rootCmd.AddCommand(cobraCmd)

	return nil
}

func cliToCobraCommand[T, F any](config *T, cmd Command[T, F]) (*cobra.Command, error) {
	cobraCmd := &cobra.Command{
		Use:        cmd.Use,
		Short:      cmd.Short,
		Long:       cmd.Long,
		Aliases:    cmd.Aliases,
		Hidden:     cmd.Hidden,
		Deprecated: cmd.Deprecated,
		Version:    cmd.Version,
	}

	var (
		flagRegistry *FlagRegistry
		err          error
	)

	if !isNoFlags(cmd.Flags) {
		flagRegistry, err = NewFlagRegistry(cmd.Flags)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to create flag registry for command %q: %w",
				cmd.Use,
				err,
			)
		}

		err = flagRegistry.RegisterFlags(cobraCmd)
		if err != nil {
			return nil, fmt.Errorf("failed to register flags for command %q: %w", cmd.Use, err)
		}
	}

	if cmd.RunE != nil {
		cobraCmd.RunE = func(c *cobra.Command, _ []string) error {
			ctx := c.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			flags, parseErr := cloneAndParseFlags(c, cmd.Flags, flagRegistry)
			if parseErr != nil {
				return fmt.Errorf("parsing flags for command %q: %w", cmd.Use, parseErr)
			}

			return cmd.RunE(ctx, config, flags)
		}
	}

	if cmd.PreRunE != nil {
		cobraCmd.PreRunE = func(c *cobra.Command, _ []string) error {
			ctx := c.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			flags, parseErr := cloneAndParseFlags(c, cmd.Flags, flagRegistry)
			if parseErr != nil {
				return fmt.Errorf("parsing flags for pre-run of command %q: %w", cmd.Use, parseErr)
			}

			return cmd.PreRunE(ctx, config, flags)
		}
	}

	if cmd.PostRunE != nil {
		cobraCmd.PostRunE = func(c *cobra.Command, _ []string) error {
			ctx := c.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			flags, parseErr := cloneAndParseFlags(c, cmd.Flags, flagRegistry)
			if parseErr != nil {
				return fmt.Errorf("parsing flags for post-run of command %q: %w", cmd.Use, parseErr)
			}

			return cmd.PostRunE(ctx, config, flags)
		}
	}

	for _, subCmd := range cmd.Commands {
		subCobraCmd, err := cliToCobraCommand(config, subCmd)
		if err != nil {
			return nil, fmt.Errorf("subcommand of %q: %w", cmd.Use, err)
		}

		cobraCmd.AddCommand(subCobraCmd)
	}

	return cobraCmd, nil
}

func isNoFlags[F any](flags F) bool {
	switch any(flags).(type) {
	case NoFlags, *NoFlags:
		return true
	default:
		return false
	}
}

// Execute runs the CLI application.
func (cli *CLI[T]) Execute(ctx context.Context) error {
	if err := fang.Execute(ctx, cli.rootCmd); err != nil {
		return fmt.Errorf("failed to execute CLI: %w", err)
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
		fmt.Println(err)
	}
}

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

// HealthCheckWithContext runs health checks with context.
func (cli *CLI[T]) HealthCheckWithContext(ctx context.Context) error {
	return cli.scope.HealthCheckWithContext(ctx)
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
	cli.long = long
	cli.rootCmd.Long = long
}

// SetVersion sets the version string.
func (cli *CLI[T]) SetVersion(version string) {
	cli.version = version
	cli.rootCmd.Version = version
}

// AddGlobalFlag adds a persistent flag available to all commands.
func (cli *CLI[T]) AddGlobalFlag(name, shorthand, defaultValue, help string) {
	cli.rootCmd.PersistentFlags().StringP(name, shorthand, defaultValue, help)
}

// AddGlobalBoolFlag adds a persistent boolean flag available to all commands.
func (cli *CLI[T]) AddGlobalBoolFlag(name, shorthand string, defaultValue bool, help string) {
	cli.rootCmd.PersistentFlags().BoolP(name, shorthand, defaultValue, help)
}
