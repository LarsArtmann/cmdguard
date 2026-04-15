package v2

import (
	"context"
	"fmt"
	"os"

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
	registeredCmds map[string]bool
	flowCtx        *BranchingFlowContext
	useFang        bool
	fangOpts       []fang.Option
	middleware     []Middleware[T]
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
	if cli.registeredCmds[cmd.Use] {
		return fmt.Errorf("%w: command %q already exists", ErrDuplicateCommand, cmd.Use)
	}

	err := cmd.Validate()
	if err != nil {
		return fmt.Errorf("validating command %q on CLI %q: %w", cmd.Use, cli.name, err)
	}

	cli.registeredCmds[cmd.Use] = true

	cobraCmd, err := cliToCobraCommand(cli.config, cmd, cli.middleware)
	if err != nil {
		return fmt.Errorf("converting command %q for CLI %q: %w", cmd.Use, cli.name, err)
	}

	cli.rootCmd.AddCommand(cobraCmd)

	return nil
}

// Execute runs the CLI application.
// The context is wrapped with a BranchingFlowContext for command path tracking.
func (cli *CLI[T]) Execute(ctx context.Context) error {
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
		return fmt.Errorf("%w: name is required, name=%q", ErrInvalidCommand, name)
	}

	return nil
}
