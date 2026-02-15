package v2

import (
	"context"
	"fmt"
)

// Command represents a type-safe CLI command with typed flags and config.
// The type parameter T is the application-level config type.
type Command[T any] struct {
	// Use is the command name and usage (e.g., "greet [name]")
	Use string

	// Short is the short description shown in help
	Short string

	// Long is the long description shown in help
	Long string

	// Aliases are alternative names for the command
	Aliases []string

	// Example shows example usage in help
	Example string

	// Flags is a struct with flag tags for command-specific flags.
	// Must be a pointer to a struct for proper flag binding.
	// Example:
	//   struct {
	//       Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
	//       Shout bool   `flag:"shout" default:"false" help:"Shout the greeting"`
	//   }
	Flags any

	// RunE is the command handler. It receives:
	// - ctx: context for cancellation and deadlines
	// - cfg: typed application-level config
	// - flags: typed command-specific flags (same type as Flags field)
	// Returns an error if the command fails.
	RunE func(ctx context.Context, cfg *T, flags any) error

	// PreRunE is called before RunE for validation.
	// Use this to validate flag combinations or prerequisites.
	PreRunE func(ctx context.Context, cfg *T, flags any) error

	// PostRunE is called after RunE for cleanup.
	// Called even if RunE returns an error.
	PostRunE func(ctx context.Context, cfg *T, flags any) error

	// Commands are subcommands of this command.
	Commands []Command[T]

	// Hidden hides the command from help output
	Hidden bool

	// Deprecated marks the command as deprecated with the given message
	Deprecated string

	// Version shows version info when --version is passed (for root command)
	Version string

	// SilenceErrors suppresses error output
	SilenceErrors bool

	// SilenceUsage suppresses usage output on error
	SilenceUsage bool
}

// Validate checks that the command is properly configured.
// Returns an error if the command is invalid.
func (c Command[T]) Validate() error {
	if c.Use == "" {
		return fmt.Errorf("%w: command has no Use field", ErrInvalidCommand)
	}

	// A command must have either a RunE handler or subcommands
	if c.RunE == nil && len(c.Commands) == 0 {
		return fmt.Errorf("%w: %q has no RunE and no subcommands", ErrMissingHandler, c.Use)
	}

	// Validate subcommands recursively
	for i, sub := range c.Commands {
		if err := sub.Validate(); err != nil {
			return fmt.Errorf("subcommand %d of %q: %w", i, c.Use, err)
		}
	}

	return nil
}

// HasSubcommands returns true if this command has subcommands.
func (c Command[T]) HasSubcommands() bool {
	return len(c.Commands) > 0
}

// HasHandler returns true if this command has a RunE handler.
func (c Command[T]) HasHandler() bool {
	return c.RunE != nil
}

// IsExecutable returns true if this command can be executed directly.
// A command is executable if it has a RunE handler.
func (c Command[T]) IsExecutable() bool {
	return c.RunE != nil
}

// CommandOption is a functional option for configuring a Command.
type CommandOption[T any] func(*Command[T])

// WithShort sets the short description.
func WithShort[T any](short string) CommandOption[T] {
	return func(c *Command[T]) {
		c.Short = short
	}
}

// WithLong sets the long description.
func WithLong[T any](long string) CommandOption[T] {
	return func(c *Command[T]) {
		c.Long = long
	}
}

// WithAliases sets the command aliases.
func WithAliases[T any](aliases ...string) CommandOption[T] {
	return func(c *Command[T]) {
		c.Aliases = aliases
	}
}

// WithExample sets the example usage.
func WithExample[T any](example string) CommandOption[T] {
	return func(c *Command[T]) {
		c.Example = example
	}
}

// WithFlags sets the command-specific flags struct.
func WithFlags[T any](flags any) CommandOption[T] {
	return func(c *Command[T]) {
		c.Flags = flags
	}
}

// WithRunE sets the command handler.
func WithRunE[T any](runE func(ctx context.Context, cfg *T, flags any) error) CommandOption[T] {
	return func(c *Command[T]) {
		c.RunE = runE
	}
}

// WithPreRunE sets the pre-run validation hook.
func WithPreRunE[T any](preRunE func(ctx context.Context, cfg *T, flags any) error) CommandOption[T] {
	return func(c *Command[T]) {
		c.PreRunE = preRunE
	}
}

// WithPostRunE sets the post-run cleanup hook.
func WithPostRunE[T any](postRunE func(ctx context.Context, cfg *T, flags any) error) CommandOption[T] {
	return func(c *Command[T]) {
		c.PostRunE = postRunE
	}
}

// WithSubcommands sets the subcommands.
func WithSubcommands[T any](cmds ...Command[T]) CommandOption[T] {
	return func(c *Command[T]) {
		c.Commands = cmds
	}
}

// WithHidden sets whether the command is hidden.
func WithHidden[T any](hidden bool) CommandOption[T] {
	return func(c *Command[T]) {
		c.Hidden = hidden
	}
}

// WithDeprecated marks the command as deprecated.
func WithDeprecated[T any](msg string) CommandOption[T] {
	return func(c *Command[T]) {
		c.Deprecated = msg
	}
}

// NewCommand creates a new command with functional options.
func NewCommand[T any](use string, opts ...CommandOption[T]) (Command[T], error) {
	if use == "" {
		return Command[T]{}, fmt.Errorf("%w: use is required", ErrMissingName)
	}

	cmd := Command[T]{Use: use}
	for _, opt := range opts {
		opt(&cmd)
	}

	if err := cmd.Validate(); err != nil {
		return Command[T]{}, err
	}

	return cmd, nil
}

// MustNewCommand creates a new command and panics on error.
// Use only in static initialization where failure is fatal.
func MustNewCommand[T any](use string, opts ...CommandOption[T]) Command[T] {
	cmd, err := NewCommand[T](use, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to create command %q: %v", use, err))
	}
	return cmd
}
