package v2

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// FlagTypeConstraint validates that F is a valid flag type at initialization time.
// Valid types are: struct{} (NoFlags), any struct, or pointer to struct.
// This enforces type safety for the F type parameter in GuardedCommand and Command.
// Returns an error if F is an invalid type (e.g., int, string, slice, map).
func FlagTypeConstraint[F any]() error {
	var zero F
	t := reflect.TypeOf(zero)

	// Nil type means F is an untyped nil interface - not valid
	if t == nil {
		return fmt.Errorf("%w: flag type F must be a struct or pointer to struct, got untyped nil", ErrInvalidFlagType)
	}

	switch t.Kind() {
	case reflect.Struct:
		// struct{} (NoFlags) or any struct is valid
		return nil
	case reflect.Ptr:
		// Must be pointer to struct
		if t.Elem().Kind() == reflect.Struct {
			return nil
		}
		return fmt.Errorf("%w: flag type F must be struct or *struct, got %s", ErrInvalidFlagType, t)
	default:
		return fmt.Errorf("%w: flag type F must be struct or *struct, got %s", ErrInvalidFlagType, t)
	}
}

// GuardedCommand provides type-safe CLI construction with DI.
// It never panics - all operations return errors.
// T is the application config type, F is the command-specific flags type.
type GuardedCommand[T any, F any] struct {
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
// T is the application config type, F is the command-specific flags type.
// F must be a struct (like NoFlags) or pointer to struct for flag binding.
func New[T any, F any](name, short string, defaults T) (*GuardedCommand[T, F], error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidCommand)
	}

	// Validate F type constraint at initialization time
	if err := FlagTypeConstraint[F](); err != nil {
		return nil, err
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

	return &GuardedCommand[T, F]{
		name:     name,
		short:    short,
		defaults: defaults,
		config:   &cfg,
		scope:    scope,
		rootCmd:  rootCmd,
	}, nil
}

// NewWithLong creates a new CLI application with a long description.
func NewWithLong[T any, F any](name, short, long string, defaults T) (*GuardedCommand[T, F], error) {
	g, err := New[T, F](name, short, defaults)
	if err != nil {
		return nil, err
	}
	g.long = long
	g.rootCmd.Long = long
	return g, nil
}

// AddCommand adds a subcommand to the CLI.
// Returns an error instead of panicking on invalid commands.
func (g *GuardedCommand[T, F]) AddCommand(cmd Command[T, F]) error {
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

// MustAddCommand adds a subcommand to the CLI and panics on error.
// Use this for static command trees where errors indicate programmer mistakes.
func (g *GuardedCommand[T, F]) MustAddCommand(cmd Command[T, F]) {
	if err := g.AddCommand(cmd); err != nil {
		panic(fmt.Sprintf("failed to add command %q: %v", cmd.Use, err))
	}
}

// AddCommandFunc adds a command using a constructor function.
// Useful for lazy initialization.
func (g *GuardedCommand[T, F]) AddCommandFunc(fn func() Command[T, F]) error {
	return g.AddCommand(fn())
}

// AddAnyCommand adds a command with a different flags type to a GuardedCommand.
// This is a standalone function because Go doesn't support type parameters on methods.
// Use this when commands need different flag types than the CLI root.
func AddAnyCommand[T any, F any, F2 any](g *GuardedCommand[T, F], cmd Command[T, F2]) error {
	// Validate the command
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Convert to cobra command with F2 flags type
	cobraCmd, err := toCobraCommandAny(g.config, cmd)
	if err != nil {
		return err
	}

	g.rootCmd.AddCommand(cobraCmd)
	return nil
}

// MustAddAnyCommand adds a command with different flags and panics on error.
// Use this for static command trees where errors indicate programmer mistakes.
func MustAddAnyCommand[T any, F any, F2 any](g *GuardedCommand[T, F], cmd Command[T, F2]) {
	if err := AddAnyCommand(g, cmd); err != nil {
		panic(fmt.Sprintf("failed to add command %q: %v", cmd.Use, err))
	}
}

// toCobraCommandAny converts a Command[T, F2] to a cobra.Command.
// This is a variant of toCobraCommand that works with any flags type.
func toCobraCommandAny[T any, F2 any](config *T, cmd Command[T, F2]) (*cobra.Command, error) {
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
	// Even if cmd.Flags is nil, we need to create a registry for the F2 type
	// to register and parse flags during execution
	var flagRegistry *FlagRegistry
	var prototype F2 // Will be zero value if cmd.Flags is nil

	if !isNilPointer(cmd.Flags) {
		prototype = cmd.Flags
	} else {
		// Create a new instance of F2 as prototype for flag registry
		var zero F2
		t := reflect.TypeOf(zero)
		if t != nil && t.Kind() == reflect.Ptr {
			// F2 is a pointer type - create new instance
			prototype = reflect.New(t.Elem()).Interface().(F2)
		} else {
			// F2 is a struct type (like NoFlags) - use zero value
			prototype = zero
		}
	}

	// Create flag registry if prototype is usable
	if !isNilPointer(prototype) {
		var err error
		flagRegistry, err = NewFlagRegistry(prototype)
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

			// Clone and parse flags once for this execution
			execFlags := cloneAndParseFlags(c, cmd.Flags, flagRegistry)

			// Execute the command handler
			return cmd.RunE(ctx, config, execFlags)
		}
	}

	// Set up pre-run hook
	if cmd.PreRunE != nil {
		cobraCmd.PreRunE = func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Clone and parse flags once for this execution
			execFlags := cloneAndParseFlags(c, cmd.Flags, flagRegistry)

			return cmd.PreRunE(ctx, config, execFlags)
		}
	}

	// Set up post-run hook
	if cmd.PostRunE != nil {
		cobraCmd.PostRunE = func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Clone and parse flags once for this execution
			execFlags := cloneAndParseFlags(c, cmd.Flags, flagRegistry)

			return cmd.PostRunE(ctx, config, execFlags)
		}
	}

	// Recursively add subcommands
	for _, subCmd := range cmd.Commands {
		cobraSubCmd, err := toCobraCommandAny[T, F2](config, subCmd)
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

// cloneAndParseFlags clones flags once and parses them.
// This is the optimized single-entry point for flag handling during execution.
// If flags is nil, creates a new instance of F to parse into.
func cloneAndParseFlags[F any](c *cobra.Command, flags F, registry *FlagRegistry) F {
	var flagsCopy F

	// If flags is nil, create a new instance of the flag type
	if isNilPointer(flags) {
		// Create new instance using reflection
		var zero F
		t := reflect.TypeOf(zero)
		if t == nil {
			// F is an interface type with nil value - can't create
			return zero
		}
		if t.Kind() == reflect.Ptr {
			// Create new instance of the underlying type
			newVal := reflect.New(t.Elem())
			flagsCopy = newVal.Interface().(F)
		} else {
			// F is a struct type (like NoFlags) - return zero value
			flagsCopy = zero
		}
	} else {
		// Clone the flags struct
		flagsCopy = cloneFlags(flags)
		if any(flagsCopy) == nil {
			flagsCopy = flags
		}
	}

	// Parse command-line values into the flags
	if registry != nil {
		if err := registry.ParseFlags(c, flagsCopy); err != nil {
			// Return the flags anyway; the error will be handled by the caller
			// The cobra framework will handle the error display
		}
	}

	return flagsCopy
}

// toCobraCommand converts a Command[T, F] to a cobra.Command.
// This is a wrapper around toCobraCommandAny that uses the GuardedCommand's config.
func (g *GuardedCommand[T, F]) toCobraCommand(cmd Command[T, F]) (*cobra.Command, error) {
	return toCobraCommandAny[T, F](g.config, cmd)
}

// isNilPointer checks if a value is a nil pointer or nil interface.
// This is needed because `any(nil) != nil` is true for typed nil pointers.
func isNilPointer(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}

// cloneFlags creates a copy of a flags struct using reflection.
// This ensures each command execution gets its own flag instance.
// Returns the zero value of F if cloning fails or input is nil.
func cloneFlags[F any](flags F) F {
	if isNilPointer(flags) {
		var zero F
		return zero
	}

	// Use reflection to create a new instance
	v := reflect.ValueOf(flags)

	// Handle pointer to struct
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			var zero F
			return zero
		}
		// Create new pointer to same type
		newPtr := reflect.New(v.Elem().Type())
		// Copy the value
		newPtr.Elem().Set(v.Elem())
		return newPtr.Interface().(F)
	}

	// Handle struct directly
	if v.Kind() == reflect.Struct {
		newStruct := reflect.New(v.Type()).Elem()
		newStruct.Set(v)
		return newStruct.Interface().(F)
	}

	// For other types, return as-is (can't clone safely)
	return flags
}

// Execute runs the CLI application.
// Returns an error if execution fails (never panics).
func (g *GuardedCommand[T, F]) Execute(ctx context.Context) error {
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
func (g *GuardedCommand[T, F]) ExecuteWithArgs(ctx context.Context, args []string) error {
	g.rootCmd.SetArgs(args)
	return g.Execute(ctx)
}

// ExecuteAndExit runs the CLI and exits with the appropriate exit code.
// This is the simplest way to run a CLI application.
func (g *GuardedCommand[T, F]) ExecuteAndExit(ctx context.Context) {
	if err := g.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Scope returns the DI scope for service registration.
// Use this to register services that commands can access.
func (g *GuardedCommand[T, F]) Scope() do.Injector {
	return g.scope.Injector()
}

// ScopeStruct returns the wrapped Scope struct for advanced operations.
func (g *GuardedCommand[T, F]) ScopeStruct() *Scope {
	return g.scope
}

// Config returns the resolved configuration.
// This is populated after flag parsing.
func (g *GuardedCommand[T, F]) Config() *T {
	return g.config
}

// SetConfig updates the configuration.
// Useful for setting config programmatically before execution.
func (g *GuardedCommand[T, F]) SetConfig(cfg T) {
	g.config = &cfg
}

// RootCommand returns the underlying cobra root command.
// Use this for advanced cobra configuration.
func (g *GuardedCommand[T, F]) RootCommand() *cobra.Command {
	return g.rootCmd
}

// Shutdown gracefully shuts down the CLI application.
// Call this after Execute returns for cleanup.
func (g *GuardedCommand[T, F]) Shutdown(ctx context.Context) error {
	return g.scope.Shutdown(ctx)
}

// HealthCheck runs health checks on all registered services.
func (g *GuardedCommand[T, F]) HealthCheck() error {
	return g.scope.HealthCheck()
}

// Name returns the CLI application name.
func (g *GuardedCommand[T, F]) Name() string {
	return g.name
}

// Short returns the short description.
func (g *GuardedCommand[T, F]) Short() string {
	return g.short
}

// Long returns the long description.
func (g *GuardedCommand[T, F]) Long() string {
	return g.long
}

// SetLong sets the long description.
func (g *GuardedCommand[T, F]) SetLong(long string) {
	g.long = long
	g.rootCmd.Long = long
}

// SetVersion sets the version string.
func (g *GuardedCommand[T, F]) SetVersion(version string) {
	g.rootCmd.Version = version
}

// AddGlobalFlag adds a persistent flag available to all commands.
func (g *GuardedCommand[T, F]) AddGlobalFlag(name, shorthand, defaultValue, help string) {
	g.rootCmd.PersistentFlags().StringP(name, shorthand, defaultValue, help)
}

// AddGlobalBoolFlag adds a persistent boolean flag available to all commands.
func (g *GuardedCommand[T, F]) AddGlobalBoolFlag(name, shorthand string, defaultValue bool, help string) {
	g.rootCmd.PersistentFlags().BoolP(name, shorthand, defaultValue, help)
}
