package v2

import (
	"context"
	"fmt"
)

// NoFlags is a convenience type for commands without command-specific flags.
// Use it as the F type parameter: Command[MyConfig, NoFlags]
type NoFlags = struct{}

// Command represents a type-safe CLI command with typed flags and config.
// The type parameter T is the application-level config type.
// The type parameter F is the command-specific flags type (use NoFlags if none).
type Command[T any, F any] struct {
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
	Flags F

	// RunE is the command handler. It receives:
	// - ctx: context for cancellation and deadlines
	// - cfg: typed application-level config
	// - flags: typed command-specific flags (same type as Flags field)
	// Returns an error if the command fails.
	RunE func(ctx context.Context, cfg *T, flags F) error

	// PreRunE is called before RunE for validation.
	// Use this to validate flag combinations or prerequisites.
	PreRunE func(ctx context.Context, cfg *T, flags F) error

	// PostRunE is called after RunE for cleanup.
	// Called even if RunE returns an error.
	PostRunE func(ctx context.Context, cfg *T, flags F) error

	// Commands are subcommands of this command.
	Commands []Command[T, F]

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
// Also checks for duplicate subcommand names within this command.
func (c Command[T, F]) Validate() error {
	if c.Use == "" {
		return fmt.Errorf("%w: command has no Use field", ErrInvalidCommand)
	}

	// A command must have either a RunE handler or subcommands
	if c.RunE == nil && len(c.Commands) == 0 {
		return fmt.Errorf("%w: %q has no RunE and no subcommands", ErrMissingHandler, c.Use)
	}

	// Check for duplicate subcommand names within this command
	seen := make(map[string]bool)
	for _, sub := range c.Commands {
		if seen[sub.Use] {
			return fmt.Errorf("%w: duplicate subcommand %q in command %q", ErrDuplicateCommand, sub.Use, c.Use)
		}
		seen[sub.Use] = true
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
func (c Command[T, F]) HasSubcommands() bool {
	return len(c.Commands) > 0
}

// HasHandler returns true if this command has a RunE handler.
func (c Command[T, F]) HasHandler() bool {
	return c.RunE != nil
}

// IsExecutable returns true if this command can be executed directly.
// A command is executable if it has a RunE handler.
func (c Command[T, F]) IsExecutable() bool {
	return c.RunE != nil
}

// CommandOption is a functional option for configuring a Command.
type CommandOption[T any, F any] func(*Command[T, F])

// WithShort sets the short description.
func WithShort[T, F any](short string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.Short = short
	}
}

// WithLong sets the long description.
func WithLong[T, F any](long string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.Long = long
	}
}

// WithAliases sets the command aliases.
func WithAliases[T, F any](aliases ...string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.Aliases = aliases
	}
}

// WithExample sets the example usage.
func WithExample[T, F any](example string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.Example = example
	}
}

// WithFlags sets the command-specific flags struct.
func WithFlags[T, F any](flags F) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.Flags = flags
	}
}

// WithRunE sets the command handler.
func WithRunE[T, F any](runE func(ctx context.Context, cfg *T, flags F) error) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.RunE = runE
	}
}

// WithPreRunE sets the pre-run validation hook.
func WithPreRunE[T, F any](preRunE func(ctx context.Context, cfg *T, flags F) error) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.PreRunE = preRunE
	}
}

// WithPostRunE sets the post-run cleanup hook.
func WithPostRunE[T, F any](postRunE func(ctx context.Context, cfg *T, flags F) error) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.PostRunE = postRunE
	}
}

// WithSubcommands sets the subcommands.
func WithSubcommands[T, F any](cmds ...Command[T, F]) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.Commands = cmds
	}
}

// WithHidden sets whether the command is hidden.
func WithHidden[T, F any](hidden bool) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.Hidden = hidden
	}
}

// WithDeprecated marks the command as deprecated.
func WithDeprecated[T, F any](msg string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.Deprecated = msg
	}
}

// NewCommand creates a new command with functional options.
func NewCommand[T, F any](use string, opts ...CommandOption[T, F]) (Command[T, F], error) {
	if use == "" {
		return Command[T, F]{}, fmt.Errorf("%w: use is required", ErrMissingName)
	}

	cmd := Command[T, F]{Use: use}
	for _, opt := range opts {
		opt(&cmd)
	}

	if err := cmd.Validate(); err != nil {
		return Command[T, F]{}, err
	}

	return cmd, nil
}

// MustNewCommand creates a new command and panics on error.
// Use only in static initialization where failure is fatal.
func MustNewCommand[T, F any](use string, opts ...CommandOption[T, F]) Command[T, F] {
	cmd, err := NewCommand[T, F](use, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to create command %q: %v", use, err))
	}
	return cmd
}
