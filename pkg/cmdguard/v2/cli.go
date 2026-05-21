package v2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"charm.land/fang/v2"
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
	registeredCmds map[string]struct{}
	flowCtx        *BranchingFlowContext
	useFang        bool
	fangOpts       []fang.Option
	middleware     []Middleware[T]
	envPrefix      string
	signalHandling bool
	outputEnabled  bool
	outputFormat   OutputFormat
	validationMode ValidationMode
	configValidate func(*T) error
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
		cli.scope = NewScope(cli.name)
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

	if cli.envPrefix != "" {
		registry.SetEnvPrefix(cli.envPrefix)
	}

	cli.initOutputFlag()

	err = registry.RegisterPersistentFlags(cli.rootCmd)
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
			return err
		}

		if cli.configValidate != nil {
			if err := cli.configValidate(cli.config); err != nil {
				return fmt.Errorf("%w: %w", ErrConfigValidation, err)
			}
		}

		return cli.parseOutputFlag(c)
	}

	return nil
}

// AddCommand adds a subcommand to the CLI with any flags type.
func AddCommand[T, F any](cli *CLI[T], cmd Command[T, F]) error {
	if _, exists := cli.registeredCmds[cmd.use]; exists {
		return fmt.Errorf("%w: command %q already exists", ErrDuplicateCommand, cmd.use)
	}

	if err := cmd.validate(cli.validationMode); err != nil {
		return fmt.Errorf("validating command %q on CLI %q: %w", cmd.use, cli.name, err)
	}

	cli.registeredCmds[cmd.use] = struct{}{}

	cobraCmd, err := cliToCobraCommand(cli.config, cmd, cli.middleware, cli.envPrefix)
	if err != nil {
		return fmt.Errorf("converting command %q for CLI %q: %w", cmd.use, cli.name, err)
	}

	cli.rootCmd.AddCommand(cobraCmd)

	return nil
}

// MustAddCommand adds a subcommand to the CLI or panics.
// Use this when the command configuration is known at compile time.
func MustAddCommand[T, F any](cli *CLI[T], cmd Command[T, F]) {
	err := AddCommand(cli, cmd)
	if err != nil {
		panic(fmt.Sprintf("MustAddCommand(%q): %v", cmd.use, err))
	}
}

// MustNewCLI creates a new CLI application or panics.
func MustNewCLI[T any](name, short string, defaults T, opts ...CLIOption[T]) *CLI[T] {
	cli, err := NewCLI(name, short, defaults, opts...)
	if err != nil {
		panic(fmt.Sprintf("MustNewCLI(%q): %v", name, err))
	}

	return cli
}

// Execute runs the CLI application.
// The context is wrapped with a BranchingFlowContext for command path tracking.
// If WithSignalHandling was set, the context is cancelled on SIGINT/SIGTERM.
func (cli *CLI[T]) Execute(ctx context.Context) error {
	if cli.signalHandling {
		var cancel context.CancelFunc

		ctx, cancel = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer cancel()
	}

	if cli.flowCtx == nil {
		cli.flowCtx = NewBranchingFlowContext(ctx)
	}

	flowCtx := WithBranchingFlowContext(ctx, cli.flowCtx)

	var execErr error

	if cli.useFang {
		execErr = fang.Execute(flowCtx, cli.rootCmd, cli.fangOpts...)
	} else {
		execErr = cli.rootCmd.ExecuteContext(flowCtx)
	}

	if execErr != nil {
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
		code := 1

		var exitCoder ExitCoder

		if errors.As(err, &exitCoder) {
			code = exitCoder.ExitCode()
		}

		os.Exit(code)
	}
}

// validateName checks that the command name is not empty.
func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required, name=%q", ErrMissingName, name)
	}

	return nil
}
