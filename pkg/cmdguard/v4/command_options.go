package v4

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// --- Non-generic metadata options ---

// WithShort sets the short description.
func WithShort(short string) CommandOption {
	return func(s *commandSpec) {
		s.short = short
	}
}

// WithLong sets the long description.
func WithLong(long string) CommandOption {
	return func(s *commandSpec) {
		s.long = long
	}
}

// WithAliases sets the command aliases.
func WithAliases(aliases ...string) CommandOption {
	return func(s *commandSpec) {
		s.aliases = aliases
	}
}

// WithExample sets the example usage.
func WithExample(example string) CommandOption {
	return func(s *commandSpec) {
		s.example = example
	}
}

// WithHidden sets whether the command is hidden.
func WithHidden(hidden bool) CommandOption {
	return func(s *commandSpec) {
		s.hidden = hidden
	}
}

// WithDeprecated marks the command as deprecated.
func WithDeprecated(msg string) CommandOption {
	return func(s *commandSpec) {
		s.deprecated = msg
	}
}

// WithGroupID assigns the command to a named group in help output.
func WithGroupID(group string) CommandOption {
	return func(s *commandSpec) {
		s.group = group
	}
}

// WithPromptOnMissing enables interactive prompting for flags that have a
// `prompt:"Question?"` struct tag and were not provided via CLI arguments or
// environment variables.
func WithPromptOnMissing() CommandOption {
	return func(s *commandSpec) {
		s.promptOnMissing = true
	}
}

// --- Arg validators (all non-generic) ---

// nonNegativeErr returns an error if n is negative.
func nonNegativeErr(name string, n int) error {
	if n < 0 {
		return fmt.Errorf("%s: %w: n=%d", name, ErrNegativeArgCount, n)
	}

	return nil
}

// WithArgs sets a custom positional arguments validator.
func WithArgs(args cobra.PositionalArgs) CommandOption {
	return func(s *commandSpec) {
		s.args = args
	}
}

// WithExactArgs requires exactly n positional arguments.
func WithExactArgs(n int) CommandOption {
	return func(s *commandSpec) {
		err := nonNegativeErr("WithExactArgs", n)
		if err != nil {
			s.optionErr = err

			return
		}

		s.args = cobra.ExactArgs(n)
	}
}

// WithMinimumArgs requires at least n positional arguments.
func WithMinimumArgs(n int) CommandOption {
	return func(s *commandSpec) {
		err := nonNegativeErr("WithMinimumArgs", n)
		if err != nil {
			s.optionErr = err

			return
		}

		s.args = cobra.MinimumNArgs(n)
	}
}

// WithMaximumArgs allows at most n positional arguments.
func WithMaximumArgs(n int) CommandOption {
	return func(s *commandSpec) {
		err := nonNegativeErr("WithMaximumArgs", n)
		if err != nil {
			s.optionErr = err

			return
		}

		s.args = cobra.MaximumNArgs(n)
	}
}

// WithRangeArgs requires between minArgs and maxArgs positional arguments (inclusive).
func WithRangeArgs(minArgs, maxArgs int) CommandOption {
	return func(s *commandSpec) {
		if minArgs < 0 {
			s.optionErr = fmt.Errorf("WithRangeArgs: %w: min=%d", ErrNegativeArgCount, minArgs)

			return
		}

		if minArgs > maxArgs {
			s.optionErr = fmt.Errorf("WithRangeArgs: %w: min=%d max=%d", ErrInvalidArgRange, minArgs, maxArgs)

			return
		}

		s.args = cobra.RangeArgs(minArgs, maxArgs)
	}
}

// WithNoArgs rejects any positional arguments.
func WithNoArgs() CommandOption {
	return func(s *commandSpec) {
		s.args = cobra.NoArgs
	}
}

// --- Completion (non-generic) ---

// WithCompletion sets the shell completion function for a command.
func WithCompletion(fn CompletionFunc) CommandOption {
	return func(s *commandSpec) {
		s.completionFn = fn
	}
}

// WithValidArgs sets static valid arguments for a command.
func WithValidArgs(args ...string) CommandOption {
	return func(s *commandSpec) {
		s.validArgs = args
	}
}

// --- Generic lifecycle options (return non-generic CommandOption) ---

// WithPreRunE sets the pre-run validation hook.
// This is one of the few options that requires type parameters — it returns
// a non-generic CommandOption that stores the typed function behind a sealed
// interface, preserving compile-time type safety.
func WithPreRunE[T, F any](
	preRunE func(ctx context.Context, cfg *T, flags F) error,
) CommandOption {
	return func(s *commandSpec) {
		s.preRunE = &typedHook[T, F]{fn: preRunE}
	}
}

// WithPostRunE sets the post-run cleanup hook.
// This is one of the few options that requires type parameters.
func WithPostRunE[T, F any](
	postRunE func(ctx context.Context, cfg *T, flags F) error,
) CommandOption {
	return func(s *commandSpec) {
		s.postRunE = &typedHook[T, F]{fn: postRunE}
	}
}

// WithSubcommands sets the subcommands.
// Type parameters are inferred from the provided commands.
func WithSubcommands[T, F any](cmds ...Command[T, F]) CommandOption {
	return func(s *commandSpec) {
		s.subcommands = &typedSubcommands[T, F]{cmds: cmds}
	}
}
