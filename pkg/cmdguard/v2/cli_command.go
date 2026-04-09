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

	var (
		flagRegistry *FlagRegistry
		err          error
	)

	if !isNoFlags(cmd.Flags) {
		prototype := createFlagPrototype(cmd.Flags)
		if !isNilPointer(prototype) {
			flagRegistry, err = NewFlagRegistry(prototype)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to create flag registry for command %q: %w",
					cmd.Use,
					err,
				)
			}

			err = flagRegistry.RegisterFlags(cobraCmd)
			if err != nil {
				return nil, fmt.Errorf("failed to register flags for command %q: %w", cmd.Use, err)
			}
		}
	}

	if cmd.RunE != nil {
		handler := cmd.RunE
		cobraCmd.RunE = func(c *cobra.Command, _ []string) error {
			ctx, flags, err := prepareRunContext(c, cmd.Flags, flagRegistry, "command "+cmd.Use)
			if err != nil {
				return err
			}

			return handler(ctx, config, flags)
		}
	}

	if cmd.PreRunE != nil {
		handler := cmd.PreRunE
		cobraCmd.PreRunE = func(c *cobra.Command, _ []string) error {
			ctx, flags, err := prepareRunContext(
				c, cmd.Flags, flagRegistry, "pre-run of command "+cmd.Use,
			)
			if err != nil {
				return err
			}

			return handler(ctx, config, flags)
		}
	}

	if cmd.PostRunE != nil {
		handler := cmd.PostRunE
		cobraCmd.PostRunE = func(c *cobra.Command, _ []string) error {
			ctx, flags, err := prepareRunContext(
				c, cmd.Flags, flagRegistry, "post-run of command "+cmd.Use,
			)
			if err != nil {
				return err
			}

			return handler(ctx, config, flags)
		}
	}

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
