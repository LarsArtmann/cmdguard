package v2

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// GuardedCommand provides type-safe CLI construction with DI.
// It never panics - all operations return errors.
type GuardedCommand[T any] struct {
	name     string
	short    string
	long     string
	defaults T
	config   *T
	scope    *Scope
	rootCmd  *cobra.Command
}

// New creates a new CLI application with typed config.
// Returns an error if initialization fails (never panics).
func New[T any](name, short string, defaults T) (*GuardedCommand[T], error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidCommand)
	}

	// Create the root cobra command
	rootCmd := &cobra.Command{
		Use:   name,
		Short: short,
	}

	// Create the DI scope
	scope := NewScope(name)

	// Register defaults in scope
	cfg := defaults
	if err := ProvideValue(scope, &cfg); err != nil {
		return nil, fmt.Errorf("failed to register config: %w", err)
	}

	// Register flag registry for root config
	registry, err := NewFlagRegistry(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create flag registry: %w", err)
	}
	if err := ProvideValue(scope, registry); err != nil {
		return nil, fmt.Errorf("failed to register flag registry: %w", err)
	}

	return &GuardedCommand[T]{
		name:     name,
		short:    short,
		defaults: defaults,
		config:   &cfg,
		scope:    scope,
		rootCmd:  rootCmd,
	}, nil
}

// NewWithLong creates a new CLI application with a long description.
func NewWithLong[T any](name, short, long string, defaults T) (*GuardedCommand[T], error) {
	g, err := New(name, short, defaults)
	if err != nil {
		return nil, err
	}
	g.long = long
	g.rootCmd.Long = long
	return g, nil
}

// AddCommand adds a subcommand to the CLI.
// Returns an error instead of panicking on invalid commands.
func (g *GuardedCommand[T]) AddCommand(cmd Command[T]) error {
	// Validate the command
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Convert to cobra command
	cobraCmd, err := g.toCobraCommand(cmd)
	if err != nil {
		return err
	}

	g.rootCmd.AddCommand(cobraCmd)
	return nil
}

// AddCommandFunc adds a command using a constructor function.
// Useful for lazy initialization.
func (g *GuardedCommand[T]) AddCommandFunc(fn func() Command[T]) error {
	return g.AddCommand(fn())
}

// toCobraCommand converts a Command[T] to a cobra.Command.
func (g *GuardedCommand[T]) toCobraCommand(cmd Command[T]) (*cobra.Command, error) {
	cobraCmd := &cobra.Command{
		Use:        cmd.Use,
		Short:      cmd.Short,
		Long:       cmd.Long,
		Aliases:    cmd.Aliases,
		Example:    cmd.Example,
		Hidden:     cmd.Hidden,
		Deprecated: cmd.Deprecated,
	}

	// Register command-specific flags
	var flagRegistry *FlagRegistry
	var flagsCopy any

	if cmd.Flags != nil {
		// Create a copy of flags for this command
		flagsCopy = cloneFlags(cmd.Flags)
		if flagsCopy == nil {
			flagsCopy = cmd.Flags
		}

		var err error
		flagRegistry, err = NewFlagRegistry(flagsCopy)
		if err != nil {
			return nil, NewCommandError(cmd.Use, fmt.Errorf("failed to create flag registry: %w", err))
		}

		if err := flagRegistry.RegisterFlags(cobraCmd); err != nil {
			return nil, NewCommandError(cmd.Use, fmt.Errorf("failed to register flags: %w", err))
		}
	}

	// Set up the command handler
	if cmd.RunE != nil {
		cobraCmd.RunE = func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Parse flags into the flags copy
			if flagRegistry != nil && flagsCopy != nil {
				if err := flagRegistry.ParseFlags(c, flagsCopy); err != nil {
					return fmt.Errorf("%w: %v", ErrFlagParseFailed, err)
				}
			}

			// Execute the command handler
			return cmd.RunE(ctx, g.config, flagsCopy)
		}
	}

	// Set up pre-run hook
	if cmd.PreRunE != nil {
		cobraCmd.PreRunE = func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Parse flags for pre-run
			var flagsCopy any
			if cmd.Flags != nil {
				flagsCopy = cloneFlags(cmd.Flags)
				if flagsCopy == nil {
					flagsCopy = cmd.Flags
				}
				if flagRegistry != nil {
					if err := flagRegistry.ParseFlags(c, flagsCopy); err != nil {
						return fmt.Errorf("%w: %v", ErrFlagParseFailed, err)
					}
				}
			}

			return cmd.PreRunE(ctx, g.config, flagsCopy)
		}
	}

	// Set up post-run hook
	if cmd.PostRunE != nil {
		cobraCmd.PostRunE = func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Parse flags for post-run
			var flagsCopy any
			if cmd.Flags != nil {
				flagsCopy = cloneFlags(cmd.Flags)
				if flagsCopy == nil {
					flagsCopy = cmd.Flags
				}
				if flagRegistry != nil {
					if err := flagRegistry.ParseFlags(c, flagsCopy); err != nil {
						return fmt.Errorf("%w: %v", ErrFlagParseFailed, err)
					}
				}
			}

			return cmd.PostRunE(ctx, g.config, flagsCopy)
		}
	}

	// Recursively add subcommands
	for _, subCmd := range cmd.Commands {
		cobraSubCmd, err := g.toCobraCommand(subCmd)
		if err != nil {
			return nil, fmt.Errorf("subcommand of %q: %w", cmd.Use, err)
		}
		cobraCmd.AddCommand(cobraSubCmd)
	}

	// Apply command options
	if cmd.SilenceErrors {
		cobraCmd.SilenceErrors = true
	}
	if cmd.SilenceUsage {
		cobraCmd.SilenceUsage = true
	}
	if cmd.Version != "" {
		cobraCmd.Version = cmd.Version
	}

	return cobraCmd, nil
}

// cloneFlags creates a copy of a flags struct using reflection.
// This ensures each command execution gets its own flag instance.
// Returns nil if cloning fails.
func cloneFlags(flags any) any {
	if flags == nil {
		return nil
	}

	// Use reflection to create a new instance
	v := reflect.ValueOf(flags)

	// Handle pointer to struct
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		// Create new pointer to same type
		newPtr := reflect.New(v.Elem().Type())
		// Copy the value
		newPtr.Elem().Set(v.Elem())
		return newPtr.Interface()
	}

	// Handle struct directly
	if v.Kind() == reflect.Struct {
		newStruct := reflect.New(v.Type()).Elem()
		newStruct.Set(v)
		return newStruct.Interface()
	}

	// For other types, return as-is (can't clone safely)
	return flags
}

// Execute runs the CLI application.
// Returns an error if execution fails (never panics).
func (g *GuardedCommand[T]) Execute(ctx context.Context) error {
	// Set context on root command
	g.rootCmd.SetContext(ctx)

	// Execute the cobra command
	err := g.rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

// ExecuteWithArgs runs the CLI application with specific arguments.
// Useful for testing.
func (g *GuardedCommand[T]) ExecuteWithArgs(ctx context.Context, args []string) error {
	g.rootCmd.SetArgs(args)
	return g.Execute(ctx)
}

// ExecuteAndExit runs the CLI and exits with the appropriate exit code.
// This is the simplest way to run a CLI application.
func (g *GuardedCommand[T]) ExecuteAndExit(ctx context.Context) {
	if err := g.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Scope returns the DI scope for service registration.
// Use this to register services that commands can access.
func (g *GuardedCommand[T]) Scope() do.Injector {
	return g.scope.Injector()
}

// ScopeStruct returns the wrapped Scope struct for advanced operations.
func (g *GuardedCommand[T]) ScopeStruct() *Scope {
	return g.scope
}

// Config returns the resolved configuration.
// This is populated after flag parsing.
func (g *GuardedCommand[T]) Config() *T {
	return g.config
}

// SetConfig updates the configuration.
// Useful for setting config programmatically before execution.
func (g *GuardedCommand[T]) SetConfig(cfg T) {
	g.config = &cfg
}

// RootCommand returns the underlying cobra root command.
// Use this for advanced cobra configuration.
func (g *GuardedCommand[T]) RootCommand() *cobra.Command {
	return g.rootCmd
}

// Shutdown gracefully shuts down the CLI application.
// Call this after Execute returns for cleanup.
func (g *GuardedCommand[T]) Shutdown(ctx context.Context) error {
	return g.scope.Shutdown(ctx)
}

// HealthCheck runs health checks on all registered services.
func (g *GuardedCommand[T]) HealthCheck() error {
	return g.scope.HealthCheck()
}

// Name returns the CLI application name.
func (g *GuardedCommand[T]) Name() string {
	return g.name
}

// Short returns the short description.
func (g *GuardedCommand[T]) Short() string {
	return g.short
}

// Long returns the long description.
func (g *GuardedCommand[T]) Long() string {
	return g.long
}

// SetLong sets the long description.
func (g *GuardedCommand[T]) SetLong(long string) {
	g.long = long
	g.rootCmd.Long = long
}

// SetVersion sets the version string.
func (g *GuardedCommand[T]) SetVersion(version string) {
	g.rootCmd.Version = version
}

// AddGlobalFlag adds a persistent flag available to all commands.
func (g *GuardedCommand[T]) AddGlobalFlag(name, shorthand, defaultValue, help string) {
	g.rootCmd.PersistentFlags().StringP(name, shorthand, defaultValue, help)
}

// AddGlobalBoolFlag adds a persistent boolean flag available to all commands.
func (g *GuardedCommand[T]) AddGlobalBoolFlag(name, shorthand string, defaultValue bool, help string) {
	g.rootCmd.PersistentFlags().BoolP(name, shorthand, defaultValue, help)
}
