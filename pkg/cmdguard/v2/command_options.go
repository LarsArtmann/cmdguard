package v2

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// CommandOption is a functional option for configuring a Command.
type CommandOption[T any, F any] func(*Command[T, F])

// WithShort sets the short description.
func WithShort[T, F any](short string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.short = short
	}
}

// WithLong sets the long description.
func WithLong[T, F any](long string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.long = long
	}
}

// WithAliases sets the command aliases.
func WithAliases[T, F any](aliases ...string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.aliases = aliases
	}
}

// WithExample sets the example usage.
func WithExample[T, F any](example string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.example = example
	}
}

// WithFlags sets the command-specific flags struct.
func WithFlags[T, F any](flags F) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.flags = flags
	}
}

// WithRunE sets the command handler.
func WithRunE[T, F any](runE func(ctx context.Context, cfg *T, flags F) error) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.runE = runE
	}
}

// WithPreRunE sets the pre-run validation hook.
func WithPreRunE[T, F any](
	preRunE func(ctx context.Context, cfg *T, flags F) error,
) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.preRunE = preRunE
	}
}

// WithPostRunE sets the post-run cleanup hook.
func WithPostRunE[T, F any](
	postRunE func(ctx context.Context, cfg *T, flags F) error,
) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.postRunE = postRunE
	}
}

// WithSubcommands sets the subcommands.
func WithSubcommands[T, F any](cmds ...Command[T, F]) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.commands = cmds
	}
}

// WithHidden sets whether the command is hidden.
func WithHidden[T, F any](hidden bool) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.hidden = hidden
	}
}

// WithDeprecated marks the command as deprecated.
func WithDeprecated[T, F any](msg string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.deprecated = msg
	}
}

// WithGroupID assigns the command to a named group in help output.
func WithGroupID[T, F any](group string) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.group = group
	}
}

// nonNegativeErr returns an error if n is negative.
func nonNegativeErr(name string, n int) error {
	if n < 0 {
		return fmt.Errorf("%s: %w: n=%d", name, ErrNegativeArgCount, n)
	}

	return nil
}

// WithArgs sets a custom positional arguments validator.
func WithArgs[T, F any](args cobra.PositionalArgs) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.args = args
	}
}

// WithExactArgs requires exactly n positional arguments.
// Returns an option that sets an error if n is negative.
func WithExactArgs[T, F any](n int) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		if err := nonNegativeErr("WithExactArgs", n); err != nil {
			c.optionErr = err

			return
		}

		c.args = cobra.ExactArgs(n)
	}
}

// WithMinimumArgs requires at least n positional arguments.
// Returns an option that sets an error if n is negative.
func WithMinimumArgs[T, F any](n int) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		if err := nonNegativeErr("WithMinimumArgs", n); err != nil {
			c.optionErr = err

			return
		}

		c.args = cobra.MinimumNArgs(n)
	}
}

// WithMaximumArgs allows at most n positional arguments.
// Returns an option that sets an error if n is negative.
func WithMaximumArgs[T, F any](n int) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		if err := nonNegativeErr("WithMaximumArgs", n); err != nil {
			c.optionErr = err

			return
		}

		c.args = cobra.MaximumNArgs(n)
	}
}

// WithRangeArgs requires between minArgs and maxArgs positional arguments (inclusive).
// Returns an option that sets an error if minArgs is negative or minArgs > maxArgs.
func WithRangeArgs[T, F any](minArgs, maxArgs int) CommandOption[T, F] {
	return func(c *Command[T, F]) {
		if minArgs < 0 {
			c.optionErr = fmt.Errorf("WithRangeArgs: %w: min=%d", ErrNegativeArgCount, minArgs)

			return
		}

		if minArgs > maxArgs {
			c.optionErr = fmt.Errorf("WithRangeArgs: %w: min=%d max=%d", ErrInvalidArgRange, minArgs, maxArgs)

			return
		}

		c.args = cobra.RangeArgs(minArgs, maxArgs)
	}
}

// WithNoArgs rejects any positional arguments.
func WithNoArgs[T, F any]() CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.args = cobra.NoArgs
	}
}

// WithPromptOnMissing enables interactive prompting for flags that have a
// `prompt:"Question?"` struct tag and were not provided via CLI arguments or
// environment variables. When a matching flag is missing, the user is prompted
// interactively using the huh TUI library.
//
// Supported field types:
//   - string, int, uint, float: text input prompt
//   - bool: yes/no confirmation prompt
//   - enum (values tag): selection prompt
func WithPromptOnMissing[T, F any]() CommandOption[T, F] {
	return func(c *Command[T, F]) {
		c.promptOnMissing = true
	}
}
