package v2

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// NoFlags is a convenience type for commands without command-specific flags.
// Use it as the F type parameter: Command[MyConfig, NoFlags].
type NoFlags struct{}

// commandSpec holds all command configuration. Non-generic fields are typed
// directly; generic fields (lifecycle hooks, subcommands) are stored as any
// and type-asserted once during cobra wiring. This is guaranteed safe because
// the generic constructors (NewCommand, NewParentCommand) ensure T and F
// match across storage and retrieval.
type commandSpec struct {
	use             string
	short           string
	long            string
	example         string
	aliases         []string
	hidden          bool
	deprecated      string
	version         string
	silenceErrors   bool
	silenceUsage    bool
	group           string
	completionFn    CompletionFunc
	validArgs       []string
	args            cobra.PositionalArgs
	promptOnMissing bool
	optionErr       error

	// Typed lifecycle hooks stored as any, set by generic option helpers
	// (WithPreRunE, WithPostRunE). Type-asserted during cobra wiring.
	preRunEAny  any
	postRunEAny any

	// Subcommands stored as any, set by WithSubcommands. Type-asserted
	// during cobra wiring.
	subcommandsAny []any
}

// CommandOption configures a command. All metadata options (descriptions,
// grouping, arg validators, completion, etc.) are non-generic — no type
// parameters needed. The few options that require type parameters
// (WithPreRunE, WithPostRunE, WithSubcommands) are generic functions that
// return a non-generic CommandOption.
type CommandOption func(*commandSpec)

// Command represents a type-safe CLI command with typed flags and config.
// Fields are unexported to enforce construction through NewCommand or
// NewParentCommand, making invalid states unrepresentable at compile time.
type Command[T any, F any] struct {
	spec     commandSpec
	flags    F
	runE     func(ctx context.Context, cfg *T, flags F) error
	commands []Command[T, F]
}

// --- Accessors ---

// Use returns the command name and usage string.
func (c Command[T, F]) Use() string { return c.spec.use }

// Short returns the short description shown in help.
func (c Command[T, F]) Short() string { return c.spec.short }

// Long returns the long description shown in help.
func (c Command[T, F]) Long() string { return c.spec.long }

// Aliases returns the alternative names for the command.
func (c Command[T, F]) Aliases() []string { return c.spec.aliases }

// Example returns the example usage string.
func (c Command[T, F]) Example() string { return c.spec.example }

// Flags returns the command-specific flags struct.
func (c Command[T, F]) Flags() F { return c.flags }

// RunE returns the command handler function.
func (c Command[T, F]) RunE() func(ctx context.Context, cfg *T, flags F) error {
	return c.runE
}

// PreRunE returns the pre-run validation hook.
func (c Command[T, F]) PreRunE() func(ctx context.Context, cfg *T, flags F) error {
	if c.spec.preRunEAny == nil {
		return nil
	}

	return c.spec.preRunEAny.(func(context.Context, *T, F) error)
}

// PostRunE returns the post-run cleanup hook.
func (c Command[T, F]) PostRunE() func(ctx context.Context, cfg *T, flags F) error {
	if c.spec.postRunEAny == nil {
		return nil
	}

	return c.spec.postRunEAny.(func(context.Context, *T, F) error)
}

// Commands returns the subcommands of this command.
func (c Command[T, F]) Commands() []Command[T, F] { return c.commands }

// Hidden returns whether the command is hidden from help output.
func (c Command[T, F]) Hidden() bool { return c.spec.hidden }

// Deprecated returns the deprecation message, if any.
func (c Command[T, F]) Deprecated() string { return c.spec.deprecated }

// Version returns the version string.
func (c Command[T, F]) Version() string { return c.spec.version }

// SilenceErrors returns whether error output is suppressed.
func (c Command[T, F]) SilenceErrors() bool { return c.spec.silenceErrors }

// SilenceUsage returns whether usage output on error is suppressed.
func (c Command[T, F]) SilenceUsage() bool { return c.spec.silenceUsage }

// Group returns the command group name.
func (c Command[T, F]) Group() string { return c.spec.group }

// HasSubcommands returns true if this command has subcommands.
func (c Command[T, F]) HasSubcommands() bool {
	return len(c.commands) > 0
}

// HasHandler returns true if this command has a RunE handler.
func (c Command[T, F]) HasHandler() bool {
	return c.runE != nil
}

// --- Validation ---

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
	if c.spec.use == "" {
		return fmt.Errorf("%w: command has no Use field", ErrInvalidCommand)
	}

	if mode >= Strict && c.spec.short == "" {
		return fmt.Errorf("%w: %q has no short description", ErrMissingShort, c.spec.use)
	}

	if mode >= Draconian && c.runE != nil && c.spec.example == "" {
		return fmt.Errorf(
			"%w: %q has no example (required in draconian mode)",
			ErrMissingExample,
			c.spec.use,
		)
	}

	if c.runE == nil && len(c.commands) == 0 {
		return fmt.Errorf(
			"%w: mode=%v, %q has no RunE and no subcommands",
			ErrMissingHandler,
			mode,
			c.spec.use,
		)
	}

	if len(c.commands) > 0 && c.spec.long == "" {
		return fmt.Errorf(
			"%w: mode=%v, %q has subcommands but no Long description",
			ErrMissingLong,
			mode,
			c.spec.use,
		)
	}

	seen := make(map[string]bool)
	for _, sub := range c.commands {
		if seen[sub.spec.use] {
			return fmt.Errorf(
				"%w: duplicate subcommand %q in command %q",
				ErrDuplicateCommand,
				sub.spec.use,
				c.spec.use,
			)
		}

		seen[sub.spec.use] = true
	}

	for i, sub := range c.commands {
		err := sub.validate(mode)
		if err != nil {
			return fmt.Errorf("mode=%v, subcommand %d of %q: %w", mode, i, c.spec.use, err)
		}
	}

	return nil
}

// requireUse returns an error if use is empty.
func requireUse(use string) error {
	if use == "" {
		return fmt.Errorf("use=%q: %w: use is required", use, ErrMissingName)
	}

	return nil
}

// --- Constructors ---

// NewCommand creates a new executable command with the given options.
// The runE parameter is required and cannot be nil.
// Use NewParentCommand for commands with subcommands.
//
// Type parameters T (config) and F (flags) are inferred from the flags and
// runE arguments — no explicit type parameters needed:
//
//	cmd, err := cmdguard.NewCommand("greet", &GreetFlags{}, greetHandler,
//	    cmdguard.WithShort("Say hello"),
//	    cmdguard.WithNoArgs(),
//	)
func NewCommand[T, F any](
	use string,
	flags F,
	runE func(ctx context.Context, cfg *T, flags F) error,
	opts ...CommandOption,
) (Command[T, F], error) {
	if err := requireUse(use); err != nil {
		return Command[T, F]{}, err
	}

	if runE == nil {
		return Command[T, F]{}, fmt.Errorf(
			"%w: runE is required for command %q",
			ErrMissingHandler,
			use,
		)
	}

	spec := commandSpec{use: use}

	for _, opt := range opts {
		opt(&spec)
	}

	cmd := Command[T, F]{spec: spec, flags: flags, runE: runE}

	err := cmd.Validate()
	if err != nil {
		return Command[T, F]{}, err
	}

	if cmd.spec.optionErr != nil {
		return Command[T, F]{}, cmd.spec.optionErr
	}

	return cmd, nil
}

// NewParentCommand creates a parent command with subcommands.
// The long description and at least one subcommand are required.
// This makes it impossible to forget the Long description when adding subcommands.
//
// Type parameter T (config) must be specified explicitly; F defaults to the
// flags type. Subcommands must share the same T and F.
//
//	cmd, err := cmdguard.NewParentCommand[AppConfig]("db", "Database ops", NoFlags{},
//	    cmdguard.WithSubcommands(migrateCmd, seedCmd),
//	)
func NewParentCommand[T, F any](
	use string,
	long string,
	flags F,
	opts ...CommandOption,
) (Command[T, F], error) {
	if err := requireUse(use); err != nil {
		return Command[T, F]{}, err
	}

	if long == "" {
		return Command[T, F]{}, fmt.Errorf(
			"%w: long=%q for parent command %q",
			ErrMissingLong,
			long,
			use,
		)
	}

	spec := commandSpec{use: use, long: long}

	for _, opt := range opts {
		opt(&spec)
	}

	cmd := Command[T, F]{spec: spec, flags: flags}

	// Extract subcommands from spec
	if len(spec.subcommandsAny) > 0 {
		cmd.commands = make([]Command[T, F], len(spec.subcommandsAny))
		for i, sub := range spec.subcommandsAny {
			cmd.commands[i] = sub.(Command[T, F])
		}
	}

	if len(cmd.commands) == 0 {
		return Command[T, F]{}, fmt.Errorf(
			"%w: long=%q, parent command %q requires at least one subcommand",
			ErrMissingHandler,
			long,
			use,
		)
	}

	err := cmd.Validate()
	if err != nil {
		return Command[T, F]{}, fmt.Errorf("long=%q: %w", long, err)
	}

	if cmd.spec.optionErr != nil {
		return Command[T, F]{}, cmd.spec.optionErr
	}

	return cmd, nil
}
