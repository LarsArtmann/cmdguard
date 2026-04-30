package v2

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type contextKeyType struct{}

var argsKey = contextKeyType{}

// ArgsFromContext retrieves the positional arguments passed to the current command.
// Use this in RunE handlers to access positional args, since cmdguard's RunE
// signature does not include a []string args parameter.
func ArgsFromContext(ctx context.Context) []string {
	if ctx == nil {
		return nil
	}

	if args, ok := ctx.Value(argsKey).([]string); ok {
		return args
	}

	return nil
}

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

func cliToCobraCommand[T, F any](
	config *T, cmd Command[T, F], middlewares []Middleware[T], envPrefix string,
) (*cobra.Command, error) {
	cobraCmd := &cobra.Command{
		Use:           cmd.use,
		Short:         cmd.short,
		Long:          cmd.long,
		Example:       cmd.example,
		Aliases:       cmd.aliases,
		Hidden:        cmd.hidden,
		Deprecated:    cmd.deprecated,
		Version:       cmd.version,
		SilenceErrors: cmd.silenceErrors,
		SilenceUsage:  cmd.silenceUsage,
	}

	if cmd.group != "" {
		cobraCmd.GroupID = cmd.group
	}

	flagRegistry, err := initCommandFlags(cobraCmd, cmd.use, cmd.flags, envPrefix)
	if err != nil {
		return nil, err
	}

	info := CommandInfo{
		Name:    cmd.use,
		HasRunE: cmd.runE != nil,
	}

	wireHandlerWithMiddleware(
		&cobraCmd.RunE, cmd.runE, config, cmd.flags, flagRegistry,
		"command "+cmd.use, info, middlewares,
	)

	preInfo := info
	preInfo.Phase = "pre-run"

	wireHandlerWithMiddleware(
		&cobraCmd.PreRunE, cmd.preRunE, config, cmd.flags, flagRegistry,
		"pre-run of command "+cmd.use, preInfo, middlewares,
	)

	postInfo := info
	postInfo.Phase = "post-run"

	wireHandlerWithMiddleware(
		&cobraCmd.PostRunE, cmd.postRunE, config, cmd.flags, flagRegistry,
		"post-run of command "+cmd.use, postInfo, middlewares,
	)

	for _, subCmd := range cmd.commands {
		subCobraCmd, err := cliToCobraCommand(config, subCmd, middlewares, envPrefix)
		if err != nil {
			return nil, fmt.Errorf("subcommand of %q: %w", cmd.use, err)
		}

		cobraCmd.AddCommand(subCobraCmd)
	}

	if len(cmd.commands) > 0 {
		wireSubcommandSuggestions(cobraCmd)
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
	cobraCmd *cobra.Command, use string, flags F, envPrefix string,
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

	if envPrefix != "" {
		registry.SetEnvPrefix(envPrefix)
	}

	err = registry.RegisterFlags(cobraCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to register flags for command %q: %w", use, err)
	}

	return registry, nil
}

func wireHandlerWithMiddleware[T, F any](
	target *func(*cobra.Command, []string) error,
	handler func(context.Context, *T, F) error,
	config *T, flags F, registry *FlagRegistry, phase string,
	info CommandInfo, middlewares []Middleware[T],
) {
	if handler == nil {
		return
	}

	h := handler
	*target = func(c *cobra.Command, args []string) error {
		ctx, parsed, err := prepareRunContext(c, flags, registry, phase)
		if err != nil {
			return err
		}

		ctx = context.WithValue(ctx, argsKey, args)

		if len(middlewares) == 0 {
			return h(ctx, config, parsed)
		}

		chain := buildChain(ctx, config, info, middlewares, func() error {
			return h(ctx, config, parsed)
		})

		return chain()
	}
}

// wireSubcommandSuggestions enhances parent commands with "did you mean?" suggestions
// when an unknown subcommand is provided.
func wireSubcommandSuggestions(cmd *cobra.Command) {
	root := cmd.Root()
	root.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return err
	})
}
