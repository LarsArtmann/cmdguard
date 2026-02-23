package v2

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// AddCommand adds a subcommand to the CLI.
// Returns an error instead of panicking on invalid commands.
// Returns ErrDuplicateCommand if a command with the same name already exists.
func (g *GuardedCommand[T, F]) AddCommand(cmd Command[T, F]) error {
	// Check for duplicate command name
	if g.registeredCmds[cmd.Use] {
		return fmt.Errorf("%w: command %q already exists", ErrDuplicateCommand, cmd.Use)
	}

	// Validate the command
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Register this command before processing to detect duplicates in subcommands
	g.registeredCmds[cmd.Use] = true

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
func (g *GuardedCommand[T, F]) AddCommandFunc(fn func() Command[T, F]) error {
	return g.AddCommand(fn())
}

// AddAnyCommand adds a command with a different flags type to a GuardedCommand.
// This is a standalone function because Go doesn't support type parameters on methods.
// Use this when commands need different flag types than the CLI root.
// Returns ErrDuplicateCommand if a command with the same name already exists.
func AddAnyCommand[T, F, F2 any](g *GuardedCommand[T, F], cmd Command[T, F2]) error {
	// Check for duplicate command name
	if g.registeredCmds[cmd.Use] {
		return fmt.Errorf("%w: command %q already exists", ErrDuplicateCommand, cmd.Use)
	}

	// Validate the command
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Register this command before processing
	g.registeredCmds[cmd.Use] = true

	// Convert to cobra command with F2 flags type
	cobraCmd, err := toCobraCommandAny(g.config, cmd)
	if err != nil {
		return err
	}

	g.rootCmd.AddCommand(cobraCmd)
	return nil
}

// toCobraCommandAny converts a Command[T, F2] to a cobra.Command.
// This is a variant of toCobraCommand that works with any flags type.
func toCobraCommandAny[T, F2 any](config *T, cmd Command[T, F2]) (*cobra.Command, error) {
	cobraCmd := createCobraCommand(cmd)
	flagRegistry, err := setupFlagRegistry(cobraCmd, cmd)
	if err != nil {
		return nil, err
	}

	setupRunHandler(cobraCmd, cmd, config, flagRegistry)
	setupPreRunHandler(cobraCmd, cmd, config, flagRegistry)
	setupPostRunHandler(cobraCmd, cmd, config, flagRegistry)

	if err := addSubcommands(cobraCmd, cmd, config); err != nil {
		return nil, err
	}

	applyCommandOptions(cobraCmd, cmd)
	return cobraCmd, nil
}

// createCobraCommand creates the base cobra.Command from Command metadata.
func createCobraCommand[T, F any](cmd Command[T, F]) *cobra.Command {
	return &cobra.Command{
		Use:        cmd.Use,
		Short:      cmd.Short,
		Long:       cmd.Long,
		Aliases:    cmd.Aliases,
		Example:    cmd.Example,
		Hidden:     cmd.Hidden,
		Deprecated: cmd.Deprecated,
	}
}

// setupFlagRegistry creates and registers flags for the command.
func setupFlagRegistry[T, F any](cobraCmd *cobra.Command, cmd Command[T, F]) (*FlagRegistry, error) {
	prototype := createFlagPrototype(cmd.Flags)
	if isNilPointer(prototype) {
		return nil, nil
	}

	registry, err := NewFlagRegistry(prototype)
	if err != nil {
		return nil, NewCommandError(cmd.Use, fmt.Errorf("failed to create flag registry: %w", err))
	}

	if err := registry.RegisterFlags(cobraCmd); err != nil {
		return nil, NewCommandError(cmd.Use, fmt.Errorf("failed to register flags: %w", err))
	}

	return registry, nil
}

// setupHandler configures a cobra handler (RunE, PreRunE, or PostRunE) for the command.
func setupHandler[T, F any](
	cobraCmd *cobra.Command,
	cmd Command[T, F],
	config *T,
	registry *FlagRegistry,
	handler func(context.Context, *T, F) error,
	setter func(*cobra.Command, func(*cobra.Command, []string) error),
) {
	if handler == nil {
		return
	}
	setter(cobraCmd, func(c *cobra.Command, args []string) error {
		return executeHandler(c, handler, cmd.Flags, config, registry)
	})
}

// setupRunHandler configures the RunE handler for the command.
func setupRunHandler[T, F any](cobraCmd *cobra.Command, cmd Command[T, F], config *T, registry *FlagRegistry) {
	setupHandler(cobraCmd, cmd, config, registry, cmd.RunE, func(c *cobra.Command, fn func(*cobra.Command, []string) error) {
		c.RunE = fn
	})
}

// setupPreRunHandler configures the PreRunE handler for the command.
func setupPreRunHandler[T, F any](cobraCmd *cobra.Command, cmd Command[T, F], config *T, registry *FlagRegistry) {
	setupHandler(cobraCmd, cmd, config, registry, cmd.PreRunE, func(c *cobra.Command, fn func(*cobra.Command, []string) error) {
		c.PreRunE = fn
	})
}

// setupPostRunHandler configures the PostRunE handler for the command.
func setupPostRunHandler[T, F any](cobraCmd *cobra.Command, cmd Command[T, F], config *T, registry *FlagRegistry) {
	setupHandler(cobraCmd, cmd, config, registry, cmd.PostRunE, func(c *cobra.Command, fn func(*cobra.Command, []string) error) {
		c.PostRunE = fn
	})
}

// executeHandler is a generic handler executor for RunE, PreRunE, PostRunE.
func executeHandler[T, F any](c *cobra.Command, handler func(context.Context, *T, F) error, flags F, config *T, registry *FlagRegistry) error {
	ctx := c.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	execFlags, err := cloneAndParseFlags(c, flags, registry)
	if err != nil {
		return err
	}

	return handler(ctx, config, execFlags)
}

// addSubcommands recursively adds subcommands to the parent command.
func addSubcommands[T, F any](parent *cobra.Command, cmd Command[T, F], config *T) error {
	for _, subCmd := range cmd.Commands {
		cobraSubCmd, err := toCobraCommandAny[T, F](config, subCmd)
		if err != nil {
			return fmt.Errorf("subcommand of %q: %w", cmd.Use, err)
		}
		parent.AddCommand(cobraSubCmd)
	}
	return nil
}

// applyCommandOptions applies command options like SilenceErrors, SilenceUsage, Version.
func applyCommandOptions[T, F any](cobraCmd *cobra.Command, cmd Command[T, F]) {
	cobraCmd.SilenceErrors = cmd.SilenceErrors
	cobraCmd.SilenceUsage = cmd.SilenceUsage
	if cmd.Version != "" {
		cobraCmd.Version = cmd.Version
	}
}
