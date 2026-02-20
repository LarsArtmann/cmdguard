package v2

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GuardedCommand provides type-safe CLI construction with DI.
// It never panics - all operations return errors.
// T is the application config type, F is the command-specific flags type.
type GuardedCommand[T any, F any] struct {
	name              string
	short             string
	long              string
	defaults          T
	config            *T
	scope             *Scope
	rootCmd           *cobra.Command
	registry          *FlagRegistry
	registeredCmds    map[string]bool // tracks registered command paths for duplicate detection
}

// New creates a new CLI application with typed config.
// Returns an error if initialization fails (never panics).
// T is the application config type, F is the command-specific flags type.
// F must be a struct (like NoFlags) or pointer to struct for flag binding.
func New[T, F any](name, short string, defaults T) (*GuardedCommand[T, F], error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	if err := FlagTypeConstraint[F](); err != nil {
		return nil, err
	}

	g := createGuardedCommand[T, F](name, short, defaults)
	if err := g.initialize(defaults); err != nil {
		return nil, err
	}

	return g, nil
}

// validateName checks that the command name is not empty.
func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidCommand)
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

	if err := g.registerConfig(defaults); err != nil {
		return err
	}

	if err := g.setupFlagRegistry(); err != nil {
		return err
	}

	return nil
}

// registerConfig registers the config in the DI scope.
func (g *GuardedCommand[T, F]) registerConfig(defaults T) error {
	cfg := defaults
	if err := ProvideValue(g.scope, &cfg); err != nil {
		return fmt.Errorf("failed to register config: %w", err)
	}
	g.config = &cfg
	return nil
}

// setupFlagRegistry creates and configures the flag registry.
func (g *GuardedCommand[T, F]) setupFlagRegistry() error {
	registry, err := NewFlagRegistry(g.config)
	if err != nil {
		return fmt.Errorf("failed to create flag registry: %w", err)
	}
	if err := ProvideValue(g.scope, registry); err != nil {
		return fmt.Errorf("failed to register flag registry: %w", err)
	}
	g.registry = registry

	if err := registry.RegisterFlags(g.rootCmd); err != nil {
		return fmt.Errorf("failed to register global flags: %w", err)
	}

	g.rootCmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		return registry.ParseFlags(c, g.config)
	}

	return nil
}

// NewWithLong creates a new CLI application with a long description.
func NewWithLong[T, F any](name, short, long string, defaults T) (*GuardedCommand[T, F], error) {
	g, err := New[T, F](name, short, defaults)
	if err != nil {
		return nil, err
	}
	g.long = long
	g.rootCmd.Long = long
	return g, nil
}

// toCobraCommand converts a Command[T, F] to a cobra.Command.
// This is a wrapper around toCobraCommandAny that uses the GuardedCommand's config.
func (g *GuardedCommand[T, F]) toCobraCommand(cmd Command[T, F]) (*cobra.Command, error) {
	return toCobraCommandAny[T, F](g.config, cmd)
}
