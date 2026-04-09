package v2

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// prepareRunContext extracts context from cobra, parses flags, and returns both.
func prepareRunContext[F any](
	c *cobra.Command, cmdFlags F, registry *FlagRegistry, phase string,
) (context.Context, F, error) {
	ctx := c.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	flags, parseErr := cloneAndParseFlags(c, cmdFlags, registry)
	if parseErr != nil {
		var zero F
		return ctx, zero, fmt.Errorf("parsing flags for %s: %w", phase, parseErr)
	}

	return ctx, flags, nil
}

func cliToCobraCommand[T, F any](config *T, cmd Command[T, F]) (*cobra.Command, error) {
	cobraCmd := &cobra.Command{
		Use:           cmd.Use,
		Short:         cmd.Short,
		Long:          cmd.Long,
		Example:       cmd.Example,
		Aliases:       cmd.Aliases,
		Hidden:        cmd.Hidden,
		Deprecated:    cmd.Deprecated,
		Version:       cmd.Version,
		SilenceErrors: cmd.SilenceErrors,
		SilenceUsage:  cmd.SilenceUsage,
	}

	flagRegistry, err := initCommandFlags(cobraCmd, cmd.Use, cmd.Flags)
	if err != nil {
		return nil, err
	}

	wireHandler(&cobraCmd.RunE, cmd.RunE, config, cmd.Flags, flagRegistry, "command "+cmd.Use)

	wireHandler(
		&cobraCmd.PreRunE, cmd.PreRunE, config, cmd.Flags, flagRegistry,
		"pre-run of command "+cmd.Use,
	)

	wireHandler(
		&cobraCmd.PostRunE, cmd.PostRunE, config, cmd.Flags, flagRegistry,
		"post-run of command "+cmd.Use,
	)

	for _, subCmd := range cmd.Commands {
		subCobraCmd, err := cliToCobraCommand(config, subCmd)
		if err != nil {
			return nil, fmt.Errorf("subcommand of %q: %w", cmd.Use, err)
		}

		cobraCmd.AddCommand(subCobraCmd)
	}

	return cobraCmd, nil
}

func isNoFlags[F any](flags F) bool {
	switch any(flags).(type) {
	case NoFlags, *NoFlags:
		return true
	default:
		return false
	}
}

func initCommandFlags[F any](
	cobraCmd *cobra.Command, use string, flags F,
) (*FlagRegistry, error) {
	if isNoFlags(flags) {
		return (*FlagRegistry)(nil), nil
	}

	prototype := createFlagPrototype(flags)
	if isNilPointer(prototype) {
		return (*FlagRegistry)(nil), nil
	}

	registry, err := NewFlagRegistry(prototype)
	if err != nil {
		return nil, fmt.Errorf("failed to create flag registry for command %q: %w", use, err)
	}

	err = registry.RegisterFlags(cobraCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to register flags for command %q: %w", use, err)
	}

	return registry, nil
}

func wireHandler[T, F any](
	target *func(*cobra.Command, []string) error,
	handler func(context.Context, *T, F) error,
	config *T, flags F, registry *FlagRegistry, phase string,
) {
	if handler == nil {
		return
	}

	h := handler
	*target = func(c *cobra.Command, _ []string) error {
		ctx, parsed, err := prepareRunContext(c, flags, registry, phase)
		if err != nil {
			return err
		}

		return h(ctx, config, parsed)
	}
}
