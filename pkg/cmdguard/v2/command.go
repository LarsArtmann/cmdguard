package v2

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// NoFlags is a convenience type for commands without command-specific flags.
// Use it as the F type parameter: Command[MyConfig, NoFlags].
type NoFlags = struct{}

// Command represents a type-safe CLI command with typed flags and config.
// Fields are unexported to enforce construction through NewCommand or NewParentCommand,
// making invalid states unrepresentable at compile time.
type Command[T any, F any] struct {
	use           string
	short         string
	long          string
	aliases       []string
	example       string
	flags         F
	runE          func(ctx context.Context, cfg *T, flags F) error
	preRunE       func(ctx context.Context, cfg *T, flags F) error
	postRunE      func(ctx context.Context, cfg *T, flags F) error
	commands      []Command[T, F]
	hidden        bool
	deprecated    string
	version       string
	silenceErrors bool
	silenceUsage  bool
	group         string
	completionFn  CompletionFunc
	validArgs     []string
	args          cobra.PositionalArgs
}

// Use returns the command name and usage string.
func (c Command[T, F]) Use() string { return c.use }

// Short returns the short description shown in help.
func (c Command[T, F]) Short() string { return c.short }

// Long returns the long description shown in help.
func (c Command[T, F]) Long() string { return c.long }

// Aliases returns the alternative names for the command.
func (c Command[T, F]) Aliases() []string { return c.aliases }

// Example returns the example usage string.
func (c Command[T, F]) Example() string { return c.example }

// Flags returns the command-specific flags struct.
func (c Command[T, F]) Flags() F { return c.flags }

// RunE returns the command handler function.
func (c Command[T, F]) RunE() func(ctx context.Context, cfg *T, flags F) error {
	return c.runE
}

// PreRunE returns the pre-run validation hook.
func (c Command[T, F]) PreRunE() func(ctx context.Context, cfg *T, flags F) error {
	return c.preRunE
}

// PostRunE returns the post-run cleanup hook.
func (c Command[T, F]) PostRunE() func(ctx context.Context, cfg *T, flags F) error {
	return c.postRunE
}

// Commands returns the subcommands of this command.
func (c Command[T, F]) Commands() []Command[T, F] { return c.commands }

// Hidden returns whether the command is hidden from help output.
func (c Command[T, F]) Hidden() bool { return c.hidden }

// Deprecated returns the deprecation message, if any.
func (c Command[T, F]) Deprecated() string { return c.deprecated }

// Version returns the version string.
func (c Command[T, F]) Version() string { return c.version }

// SilenceErrors returns whether error output is suppressed.
func (c Command[T, F]) SilenceErrors() bool { return c.silenceErrors }

// SilenceUsage returns whether usage output on error is suppressed.
func (c Command[T, F]) SilenceUsage() bool { return c.silenceUsage }

// Group returns the command group name.
func (c Command[T, F]) Group() string { return c.group }

// HasSubcommands returns true if this command has subcommands.
func (c Command[T, F]) HasSubcommands() bool {
	return len(c.commands) > 0
}

// HasHandler returns true if this command has a RunE handler.
func (c Command[T, F]) HasHandler() bool {
	return c.runE != nil
}

// IsExecutable returns true if this command can be executed directly.
//
// Deprecated: Use HasHandler() instead. Will be removed in v3.
func (c Command[T, F]) IsExecutable() bool {
	return c.HasHandler()
}

// ValidationMode controls how strictly commands are validated.
type ValidationMode int

const (
	// Lenient requires only name and handler.
	Lenient ValidationMode = iota
	// Strict additionally requires short descriptions on all commands.
	Strict
	// Draconian additionally requires examples on leaf commands.
	Draconian
)

// Validate checks that the command is properly configured.
func (c Command[T, F]) Validate() error {
	return c.validate(Lenient)
}

// ValidateStrict checks that the command is properly configured with strict rules.
// In strict mode, all commands must have a short description.
func (c Command[T, F]) ValidateStrict() error {
	return c.validate(Strict)
}

func (c Command[T, F]) validate(mode ValidationMode) error {
	if c.use == "" {
		return fmt.Errorf("%w: command has no Use field", ErrInvalidCommand)
	}

	if mode >= Strict && c.short == "" {
		return fmt.Errorf("%w: %q has no short description", ErrMissingShort, c.use)
	}

	if mode >= Draconian && c.runE != nil && c.example == "" {
		return fmt.Errorf(
			"%w: %q has no example (required in draconian mode)",
			ErrMissingExample,
			c.use,
		)
	}

	if c.runE == nil && len(c.commands) == 0 {
		return fmt.Errorf(
			"%w: mode=%v, %q has no RunE and no subcommands",
			ErrMissingHandler,
			mode,
			c.use,
		)
	}

	if len(c.commands) > 0 && c.long == "" {
		return fmt.Errorf(
			"%w: mode=%v, %q has subcommands but no Long description",
			ErrMissingLong,
			mode,
			c.use,
		)
	}

	seen := make(map[string]bool)
	for _, sub := range c.commands {
		if seen[sub.use] {
			return fmt.Errorf(
				"%w: duplicate subcommand %q in command %q",
				ErrDuplicateCommand,
				sub.use,
				c.use,
			)
		}

		seen[sub.use] = true
	}

	for i, sub := range c.commands {
		err := sub.validate(mode)
		if err != nil {
			return fmt.Errorf("mode=%v, subcommand %d of %q: %w", mode, i, c.use, err)
		}
	}

	return nil
}

// NewCommand creates a new executable command with the given options.
// The runE parameter is required and cannot be nil.
// Use NewParentCommand for commands with subcommands.
func NewCommand[T, F any](
	use string,
	runE func(ctx context.Context, cfg *T, flags F) error,
	opts ...CommandOption[T, F],
) (Command[T, F], error) {
	if use == "" {
		return Command[T, F]{}, fmt.Errorf("%w: use is required", ErrMissingName)
	}

	if runE == nil {
		return Command[T, F]{}, fmt.Errorf(
			"%w: runE is required for command %q",
			ErrMissingHandler,
			use,
		)
	}

	cmd := Command[T, F]{use: use, runE: runE}
	for _, opt := range opts {
		opt(&cmd)
	}

	err := cmd.Validate()
	if err != nil {
		return Command[T, F]{}, err
	}

	return cmd, nil
}

// NewParentCommand creates a parent command with subcommands.
// The long description and at least one subcommand are required.
// This makes it impossible to forget the Long description when adding subcommands.
func NewParentCommand[T, F any](
	use string,
	long string,
	subcommands []Command[T, F],
	opts ...CommandOption[T, F],
) (Command[T, F], error) {
	if use == "" {
		return Command[T, F]{}, fmt.Errorf("%w: use is required", ErrMissingName)
	}

	if long == "" {
		return Command[T, F]{}, fmt.Errorf(
			"%w: long=%q for parent command %q",
			ErrMissingLong,
			long,
			use,
		)
	}

	if len(subcommands) == 0 {
		return Command[T, F]{}, fmt.Errorf(
			"%w: long=%q, parent command %q requires at least one subcommand",
			ErrMissingHandler,
			long,
			use,
		)
	}

	cmd := Command[T, F]{use: use, long: long, commands: subcommands}
	for _, opt := range opts {
		opt(&cmd)
	}

	err := cmd.Validate()
	if err != nil {
		return Command[T, F]{}, fmt.Errorf("long=%q: %w", long, err)
	}

	return cmd, nil
}

// MustNewCommand creates a leaf command or panics.
// Use this when the command configuration is known at compile time.
func MustNewCommand[T, F any](
	use string,
	runE func(ctx context.Context, cfg *T, flags F) error,
	opts ...CommandOption[T, F],
) Command[T, F] {
	cmd, err := NewCommand(use, runE, opts...)
	if err != nil {
		panic(fmt.Sprintf("MustNewCommand(%q): %v", use, err))
	}

	return cmd
}

// MustNewParentCommand creates a parent command or panics.
// Use this when the command configuration is known at compile time.
func MustNewParentCommand[T, F any](
	use string,
	long string,
	subcommands []Command[T, F],
	opts ...CommandOption[T, F],
) Command[T, F] {
	cmd, err := NewParentCommand(use, long, subcommands, opts...)
	if err != nil {
		panic(fmt.Sprintf("MustNewParentCommand(%q): %v", use, err))
	}

	return cmd
}
