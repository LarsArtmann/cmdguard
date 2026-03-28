package v2

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the current version of the cmdguard v2 package.
// This follows semantic versioning (https://semver.org/).
const Version = "2.1.0"

// GuardedCommand provides type-safe CLI construction with DI.
// It never panics - all operations return errors.
// T is the application config type, F is the command-specific flags type.
//
// Deprecated: Use CLI[T] instead, which allows each command to have its own flags type.
// GuardedCommand[T, F] will be removed in v3.0. See MIGRATION.md for upgrade guide.
// Example migration:
//
//	// Old:
//	cli, err := v2.New[Config, v2.NoFlags]("myapp", "My App", Config{})
//
//	// New:
//	cli, err := v2.NewCLI[Config]("myapp", "My App", Config{})
type GuardedCommand[T any, F any] struct {
	name           string
	short          string
	long           string
	defaults       T
	config         *T
	scope          *Scope
	rootCmd        *cobra.Command
	registry       *FlagRegistry
	registeredCmds map[string]bool // tracks registered command paths for duplicate detection
}

// New creates a new CLI application with typed config.
// Returns an error if initialization fails (never panics).
// T is the application config type, F is the command-specific flags type.
// F must be a struct (like NoFlags) or pointer to struct for flag binding.
//
// Deprecated: Use NewCLI[T] instead, which allows each command to have its own flags type.
// This function will be removed in v3.0.
func New[T, F any](name, short string, defaults T) (*GuardedCommand[T, F], error) {
	err := validateName(name)
	if err != nil {
		return nil, err
	}

	err = FlagTypeConstraint[F]()
	if err != nil {
		return nil, err
	}

	g := createGuardedCommand[T, F](name, short, defaults)

	err = g.initialize(defaults)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// validateName checks that the command name is not empty.
func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required, name=%q", ErrInvalidCommand, name)
	}

	return nil
}

// createGuardedCommand creates a new GuardedCommand with basic fields set.
func createGuardedCommand[T, F any](name, short string, defaults T) *GuardedCommand[T, F] {
	return &GuardedCommand[T, F]{
		name:           name,
		short:          short,
		defaults:       defaults,
		scope:          nil, // initialized below
		rootCmd:        &cobra.Command{Use: name, Short: short},
		registry:       nil, // initialized below
		registeredCmds: make(map[string]bool),
	}
}

// initialize sets up the DI scope, flag registry, and global flags.
func (g *GuardedCommand[T, F]) initialize(defaults T) error {
	g.scope = NewScope(g.name)

	err := g.registerConfig(defaults)
	if err != nil {
		return err
	}

	err = g.setupFlagRegistry()
	if err != nil {
		return err
	}

	return nil
}

// registerConfig registers the config in the DI scope.
func (g *GuardedCommand[T, F]) registerConfig(defaults T) error {
	cfg := defaults

	err := ProvideValue(g.scope, &cfg)
	if err != nil {
		return fmt.Errorf("failed to register config type=%T: %w", cfg, err)
	}

	g.config = &cfg

	return nil
}

// setupFlagRegistry creates and configures the flag registry.
func (g *GuardedCommand[T, F]) setupFlagRegistry() error {
	registry, err := NewFlagRegistry(g.config)
	if err != nil {
		return fmt.Errorf("failed to create flag registry: config=%T: %w", g.config, err)
	}

	if err := ProvideValue(g.scope, registry); err != nil {
		return fmt.Errorf("failed to register flag registry: %w", err)
	}

	g.registry = registry

	if err := registry.RegisterFlags(g.rootCmd); err != nil {
		return fmt.Errorf("failed to register global flags: %w", err)
	}

	g.rootCmd.PersistentPreRunE = func(c *cobra.Command, _ []string) error {
		return registry.ParseFlags(c, g.config)
	}

	return nil
}

// NewWithLong creates a new CLI application with a long description.
//
// Deprecated: Use NewCLI[T] with WithCLILong[T] option instead.
func NewWithLong[T, F any](name, short, long string, defaults T) (*GuardedCommand[T, F], error) {
	g, err := New[T, F](name, short, defaults)
	if err != nil {
		return nil, err
	}

	g.long = long
	g.rootCmd.Long = long

	return g, nil
}

// SimpleCLI is a type alias for CLIs that don't need command-specific flags.
// Use this when your CLI only has global config (T) and all commands share NoFlags.
// This is an alias, not a new type, so it works seamlessly with all GuardedCommand methods.
//
// Deprecated: Use CLI[T] instead, which provides the same functionality with a cleaner API.
type SimpleCLI[T any] = GuardedCommand[T, NoFlags]

// NewSimple creates a new CLI application with typed config and no command-specific flags.
// This is a convenience wrapper around New[T, NoFlags] for the common case where
// commands don't need additional flags beyond the global config.
//
// Deprecated: Use NewCLI[T] instead.
func NewSimple[T any](name, short string, defaults T) (*SimpleCLI[T], error) {
	return New[T, NoFlags](name, short, defaults)
}

// NewSimpleWithLong creates a new CLI application with a long description and no command-specific flags.
// This is a convenience wrapper around NewWithLong[T, NoFlags].
//
// Deprecated: Use NewCLI[T] with WithCLILong[T] option instead.
func NewSimpleWithLong[T any](name, short, long string, defaults T) (*SimpleCLI[T], error) {
	return NewWithLong[T, NoFlags](name, short, long, defaults)
}

// toCobraCommand converts a Command[T, F] to a cobra.Command.
// This is a wrapper around toCobraCommandAny that uses the GuardedCommand's config.
func (g *GuardedCommand[T, F]) toCobraCommand(cmd Command[T, F]) (*cobra.Command, error) {
	return toCobraCommandAny[T, F](g.config, cmd)
}
