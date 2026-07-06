package v2

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type contextKeyType struct{}

var argsKey = contextKeyType{}

var configKey = contextKeyType{}

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

// ConfigFromContext retrieves the typed config pointer stored by cmdguard's
// PersistentPreRunE. This is the bridge for consumers that register raw
// *cobra.Command subcommands (via cli.RootCommand().AddCommand) and need
// access to the resolved config without building a parallel context-key system.
//
// Returns (*T, true) when the config was stored; (nil, false) when it was not
// (e.g., in unit tests calling RunE directly without going through Execute).
func ConfigFromContext[T any](ctx context.Context) (*T, bool) {
	if ctx == nil {
		return nil, false
	}

	cfg, ok := ctx.Value(configKey).(*T)

	return cfg, ok
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
	s := cmd.spec

	cobraCmd := &cobra.Command{
		Use:           s.use,
		Short:         s.short,
		Long:          s.long,
		Example:       s.example,
		Aliases:       s.aliases,
		Hidden:        s.hidden,
		Deprecated:    s.deprecated,
		Version:       s.version,
		SilenceErrors: s.silenceErrors,
		SilenceUsage:  s.silenceUsage,
	}

	if s.group != "" {
		cobraCmd.GroupID = s.group
	}

	flagRegistry, err := initCommandFlags(cobraCmd, s.use, cmd.flags, envPrefix)
	if err != nil {
		return nil, fmt.Errorf("envPrefix=%s: %w", envPrefix, err)
	}

	wireAllHandlers(cobraCmd, config, cmd, flagRegistry, middlewares)

	for _, subCmd := range cmd.commands {
		subCobraCmd, err := cliToCobraCommand(config, subCmd, middlewares, envPrefix)
		if err != nil {
			return nil, fmt.Errorf("envPrefix=%s, subcommand of %q: %w", envPrefix, s.use, err)
		}

		cobraCmd.AddCommand(subCobraCmd)
	}

	if s.completionFn != nil {
		cobraCmd.ValidArgsFunction = s.completionFn
	}

	if len(s.validArgs) > 0 {
		cobraCmd.ValidArgs = s.validArgs
	}

	if s.args != nil {
		cobraCmd.Args = s.args
	}

	return cobraCmd, nil
}

func wireAllHandlers[T, F any](
	cobraCmd *cobra.Command, config *T, cmd Command[T, F],
	flagRegistry *FlagRegistry, middlewares []Middleware[T],
) {
	s := cmd.spec
	info := CommandInfo{Name: s.use, Phase: PhaseRun, HasRunE: cmd.runE != nil}

	wireHandlerWithMiddleware(handlerConfig[T, F]{
		target: &cobraCmd.RunE, handler: cmd.runE, config: config,
		flags: cmd.flags, registry: flagRegistry,
		phase: "command " + s.use, info: info, middlewares: middlewares,
		promptOnMissing: s.promptOnMissing,
	})

	// Type-assert stored lifecycle hooks — safe because generic constructors
	// ensure T and F match across storage (WithPreRunE/WithPostRunE) and
	// retrieval (here).
	var preRunE func(context.Context, *T, F) error
	if s.preRunEAny != nil {
		preRunE = s.preRunEAny.(func(context.Context, *T, F) error)
	}

	var postRunE func(context.Context, *T, F) error
	if s.postRunEAny != nil {
		postRunE = s.postRunEAny.(func(context.Context, *T, F) error)
	}

	preInfo := info
	preInfo.Phase = PhasePreRun

	wireHandlerWithMiddleware(handlerConfig[T, F]{
		target: &cobraCmd.PreRunE, handler: preRunE, config: config,
		flags: cmd.flags, registry: flagRegistry,
		phase: "pre-run of command " + s.use, info: preInfo, middlewares: middlewares,
		promptOnMissing: s.promptOnMissing,
	})

	postInfo := info
	postInfo.Phase = PhasePostRun

	wireHandlerWithMiddleware(handlerConfig[T, F]{
		target: &cobraCmd.PostRunE, handler: postRunE, config: config,
		flags: cmd.flags, registry: flagRegistry,
		phase: "post-run of command " + s.use, info: postInfo, middlewares: middlewares,
		promptOnMissing: s.promptOnMissing,
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
