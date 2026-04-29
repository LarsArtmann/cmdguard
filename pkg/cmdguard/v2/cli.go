package v2

import (
	"context"
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
	name            string
	short           string
	long            string
	version         string
	defaults        T
	config          *T
	scope           *Scope
	rootCmd         *cobra.Command
	registry        *FlagRegistry
	registeredCmds  map[string]bool
	flowCtx         *BranchingFlowContext
	useFang         bool
	fangOpts        []fang.Option
	middleware       []Middleware[T]
	envPrefix       string
	signalHandling  bool
}

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
		useFang:        true,
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
	if cli.scope == nil {
		cli.scope = NewScope(cli.name)
	}

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

	if cli.envPrefix != "" {
		registry.SetEnvPrefix(cli.envPrefix)
	}

	err = registry.RegisterPersistentFlags(cli.rootCmd)
	if err != nil {
		return fmt.Errorf("registering global flags for %T: %w", defaults, err)
	}

	cli.rootCmd.PersistentPreRunE = func(c *cobra.Command, _ []string) error {
		return registry.ParseFlags(c, cli.config)
	}

	return nil
}

// AddCommand adds a subcommand to the CLI with any flags type.
func AddCommand[T, F any](cli *CLI[T], cmd Command[T, F]) error {
	if cli.registeredCmds[cmd.use] {
		return fmt.Errorf("%w: command %q already exists", ErrDuplicateCommand, cmd.use)
	}

	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("validating command %q on CLI %q: %w", cmd.use, cli.name, err)
	}

	cli.registeredCmds[cmd.use] = true

	cobraCmd, err := cliToCobraCommand(cli.config, cmd, cli.middleware)
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
func (cli *CLI[T]) ExecuteAndExit(ctx context.Context) {
	err := cli.Execute(ctx)
	if err != nil {
		os.Exit(1)
	}
}

// validateName checks that the command name is not empty.
func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required, name=%q", ErrMissingName, name)
	}

	return nil
}
