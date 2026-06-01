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
		return nil, fmt.Errorf("envPrefix=%s: %w", envPrefix, err)
	}

	wireAllHandlers(cobraCmd, config, cmd, flagRegistry, middlewares)

	for _, subCmd := range cmd.commands {
		subCobraCmd, err := cliToCobraCommand(config, subCmd, middlewares, envPrefix)
		if err != nil {
			return nil, fmt.Errorf("envPrefix=%s, subcommand of %q: %w", envPrefix, cmd.use, err)
		}

		cobraCmd.AddCommand(subCobraCmd)
	}

	if cmd.completionFn != nil {
		cobraCmd.ValidArgsFunction = cmd.completionFn
	}

	if len(cmd.validArgs) > 0 {
		cobraCmd.ValidArgs = cmd.validArgs
	}

	if cmd.args != nil {
		cobraCmd.Args = cmd.args
	}

	return cobraCmd, nil
}

func wireAllHandlers[T, F any](
	cobraCmd *cobra.Command, config *T, cmd Command[T, F],
	flagRegistry *FlagRegistry, middlewares []Middleware[T],
) {
	info := CommandInfo{Name: cmd.use, Phase: PhaseRun, HasRunE: cmd.runE != nil}

	wireHandlerWithMiddleware(handlerConfig[T, F]{
		target: &cobraCmd.RunE, handler: cmd.runE, config: config,
		flags: cmd.flags, registry: flagRegistry,
		phase: "command " + cmd.use, info: info, middlewares: middlewares,
		promptOnMissing: cmd.promptOnMissing,
	})

	preInfo := info
	preInfo.Phase = PhasePreRun

	wireHandlerWithMiddleware(handlerConfig[T, F]{
		target: &cobraCmd.PreRunE, handler: cmd.preRunE, config: config,
		flags: cmd.flags, registry: flagRegistry,
		phase: "pre-run of command " + cmd.use, info: preInfo, middlewares: middlewares,
		promptOnMissing: cmd.promptOnMissing,
	})

	postInfo := info
	postInfo.Phase = PhasePostRun

	wireHandlerWithMiddleware(handlerConfig[T, F]{
		target: &cobraCmd.PostRunE, handler: cmd.postRunE, config: config,
		flags: cmd.flags, registry: flagRegistry,
		phase: "post-run of command " + cmd.use, info: postInfo, middlewares: middlewares,
		promptOnMissing: cmd.promptOnMissing,
	})
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
		return nil, fmt.Errorf("envPrefix=%s, command %q: %w", envPrefix, use, err)
	}

	if envPrefix != "" {
		registry.SetEnvPrefix(envPrefix)
	}

	err = registry.RegisterFlags(cobraCmd)
	if err != nil {
		return nil, fmt.Errorf("envPrefix=%s, command %q: %w", envPrefix, use, err)
	}

	return registry, nil
}

type handlerConfig[T, F any] struct {
	target          *func(*cobra.Command, []string) error
	handler         func(context.Context, *T, F) error
	config          *T
	flags           F
	registry        *FlagRegistry
	phase           string
	info            CommandInfo
	middlewares     []Middleware[T]
	promptOnMissing bool
}

func wireHandlerWithMiddleware[T, F any](cfg handlerConfig[T, F]) {
	if cfg.handler == nil {
		return
	}

	h := cfg.handler

	*cfg.target = func(c *cobra.Command, args []string) error {
		info := cfg.info
		info.FullPath = c.CommandPath()

		if cfg.promptOnMissing && cfg.info.Phase == PhaseRun {
			if err := promptMissingCommandFlags(c, cfg.registry); err != nil {
				return fmt.Errorf("prompting for missing command flags: %w", err)
			}
		}

		ctx, parsed, err := prepareRunContext(c, cfg.flags, cfg.registry, cfg.phase)
		if err != nil {
			return fmt.Errorf("preparing run context: %w", err)
		}

		ctx = context.WithValue(ctx, argsKey, args)

		if len(cfg.middlewares) == 0 {
			return h(ctx, cfg.config, parsed)
		}

		chain := buildChain(ctx, cfg.config, info, cfg.middlewares, func() error {
			return h(ctx, cfg.config, parsed)
		})

		return chain()
	}
}
