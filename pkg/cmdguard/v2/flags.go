package v2

import (
	"fmt"
	"iter"
	"reflect"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// FlagRegistry manages flag registration and parsing.
// Each instance shares type and validator registries from package-level defaults
// via copy-on-write: reads use the shared global maps, writes trigger a lazy clone.
// This avoids the clone cost when no per-instance customization is used.
type FlagRegistry struct {
	tags       []FlagTag
	envPrefix  string
	validators *validatorRegistry
	types      *typeRegistry
}

// Tags returns a defensive copy of all parsed flag tags.
// For zero-allocation iteration, prefer TagsSeq().
func (r *FlagRegistry) Tags() []FlagTag {
	return slices.Clone(r.tags)
}

// TagsSeq returns an iterator over all parsed flag tags without allocating a slice.
// Prefer this over Tags() when you only need to range over the tags.
func (r *FlagRegistry) TagsSeq() iter.Seq[FlagTag] {
	return slices.Values(r.tags)
}

// SetEnvPrefix sets the environment variable prefix for this registry.
// When set, env tag lookups prepend this prefix (e.g., prefix "APP_" + tag "PORT" → "APP_PORT").
func (r *FlagRegistry) SetEnvPrefix(prefix string) {
	r.envPrefix = prefix
}

// NewFlagRegistry creates a new FlagRegistry from a config struct.
func NewFlagRegistry(cfg any) (*FlagRegistry, error) {
	tags, err := ParseFlagTags(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing flag tags for %T: %w", ErrFlagParseFailed, cfg, err)
	}

	return &FlagRegistry{
		tags:       tags,
		validators: globalValidators.share(),
		types:      globalTypeRegistry.share(),
	}, nil
}

// RegisterFlagValidator adds a named validator to this registry's instance-scoped set.
func (r *FlagRegistry) RegisterFlagValidator(name string, validator FlagValidator) {
	r.validators.register(name, validator)
}

// RegisterTypeHandler registers a custom TypeHandler for a specific reflect.Type
// on this registry instance. New flags registered after this call will use the handler.
func (r *FlagRegistry) RegisterTypeHandler(typ reflect.Type, handler TypeHandler) {
	r.types.register(typ, handler)
}

// RegisterFlags adds flags to a cobra command based on the config struct.
// Flags are registered as local (non-persistent) — they only apply to this command.
func (r *FlagRegistry) RegisterFlags(cmd *cobra.Command) error {
	return r.registerAllFlags(cmd.Flags(), cmd)
}

// RegisterPersistentFlags adds flags as persistent — they propagate to all subcommands.
// Use this for root-level config flags that must be available on every command.
func (r *FlagRegistry) RegisterPersistentFlags(cmd *cobra.Command) error {
	return r.registerAllFlags(cmd.PersistentFlags(), cmd)
}

// RegisterScopedFlags registers each flag according to its declared scope.
// Persistent flags (the default) propagate to subcommands; local flags
// (fields tagged local:"true") are registered on this command only and are NOT
// inherited by subcommands. Use this on the root command when some root flags
// are only meaningful for the root run and would otherwise pollute every
// subcommand's --help.
func (r *FlagRegistry) RegisterScopedFlags(cmd *cobra.Command) error {
	for _, tag := range r.tags {
		flagSet := cmd.PersistentFlags()
		if tag.Local {
			flagSet = cmd.Flags()
		}

		if err := r.registerFlag(flagSet, tag); err != nil {
			return fmt.Errorf(
				"%w: registering scoped flags on command %q: %w",
				ErrFlagParseFailed,
				cmd.Use,
				err,
			)
		}
	}

	return nil
}

// RegisterLocalFlags registers only the local-scoped flags (local:"true") on the
// command's own flag set. Use this on a subcommand that shares a group of
// root-only flags (e.g. an execution subcommand) so those flags are accepted
// and parsed without making them persistent across the whole command tree.
//
// Flags the subcommand already defines itself are skipped — the subcommand's
// own definition wins (for example, a subcommand with its own --build-mode keeps
// its distinct default and help). Has no effect when the config has no
// local-scoped flags.
func (r *FlagRegistry) RegisterLocalFlags(cmd *cobra.Command) error {
	for _, tag := range r.tags {
		if !tag.Local {
			continue
		}

		if cmd.Flags().Lookup(tag.Name) != nil {
			continue
		}

		if err := r.registerFlag(cmd.Flags(), tag); err != nil {
			return fmt.Errorf(
				"%w: registering local flags on command %q: %w",
				ErrFlagParseFailed,
				cmd.Use,
				err,
			)
		}
	}

	return nil
}

// registerAllFlags registers all flags to the given flag set.
func (r *FlagRegistry) registerAllFlags(flagSet *pflag.FlagSet, cmd *cobra.Command) error {
	for _, tag := range r.tags {
		err := r.registerFlag(flagSet, tag)
		if err != nil {
			return fmt.Errorf(
				"%w: registering flags on command %q: %w",
				ErrFlagParseFailed,
				cmd.Use,
				err,
			)
		}
	}

	return nil
}

// registerFlag adds a single flag to the given flag set via the TypeHandler registry.
func (r *FlagRegistry) registerFlag(flags *pflag.FlagSet, tag FlagTag) error {
	if err := dispatchRegister(r.types, flags, tag); err != nil {
		return fmt.Errorf("%w: registering flag %q: %w", ErrFlagParseFailed, tag.Name, err)
	}

	if tag.Hidden {
		if f := flags.Lookup(tag.Name); f != nil {
			f.Hidden = true
		}
	}

	return nil
}

// ValidateFlags validates flag values against allowed values and checks required flags.
func (r *FlagRegistry) ValidateFlags(cmd *cobra.Command) error {
	for _, tag := range r.tags {
		err := r.validateTag(cmd, tag)
		if err != nil {
			return fmt.Errorf(
				"%w: validating flag %q on command %q: %w",
				ErrFlagParseFailed,
				tag.Name,
				cmd.Use,
				err,
			)
		}
	}

	return nil
}

// validateTag validates a single flag tag.
func (r *FlagRegistry) validateTag(cmd *cobra.Command, tag FlagTag) error {
	err := r.validateRequiredFlag(cmd, tag)
	if err != nil {
		return fmt.Errorf(
			"%w: validating required flag %q on command %q: %w",
			ErrRequiredFlag,
			tag.Name,
			cmd.Use,
			err,
		)
	}

	err = r.validateEnumValue(cmd, tag)
	if err != nil {
		return fmt.Errorf("validating enum value for flag %q: %w", tag.Name, err)
	}

	return r.validateTagRules(cmd, tag)
}

// validateRequiredFlag checks if a required flag was set.
func (r *FlagRegistry) validateRequiredFlag(cmd *cobra.Command, tag FlagTag) error {
	if !tag.Required {
		return nil
	}

	flag := r.lookupFlagForValidation(cmd, tag.Name)

	// A local flag is only present on the command it was registered on; it
	// cannot be required on a command where it does not exist.
	if flag == nil && tag.Local {
		return nil
	}

	if flag == nil || !flag.Changed {
		return fmt.Errorf(
			"required flag %q not set on command %q: %w",
			tag.Name,
			cmd.Use,
			ErrRequiredFlag,
		)
	}

	return nil
}

// validateEnumValue validates that an enum flag has an allowed value.
func (r *FlagRegistry) validateEnumValue(cmd *cobra.Command, tag FlagTag) error {
	if len(tag.Values) == 0 {
		return nil
	}

	flag := r.lookupFlagForValidation(cmd, tag.Name)
	if flag == nil || !flag.Changed {
		return nil
	}

	if !r.isAllowedValue(flag.Value.String(), tag.Values) {
		return fmt.Errorf("enum flag %q on command %q has invalid value %q: must be one of: %v: %w",
			tag.Name, cmd.Use, flag.Value.String(), tag.Values, ErrInvalidEnum)
	}

	return nil
}

// lookupFlagForValidation finds a flag by name for validation purposes.
func (r *FlagRegistry) lookupFlagForValidation(cmd *cobra.Command, name string) *pflag.Flag {
	return lookupFlagInCommand(cmd, name)
}

// lookupFlagInCommand finds a flag by name, checking local then persistent flags.
func lookupFlagInCommand(cmd *cobra.Command, name string) *pflag.Flag {
	flag := cmd.Flags().Lookup(name)
	if flag == nil {
		flag = cmd.PersistentFlags().Lookup(name)
	}

	return flag
}

// isAllowedValue checks if a value is in the allowed list.
func (r *FlagRegistry) isAllowedValue(value string, allowed []string) bool {
	return slices.Contains(allowed, value)
}

// validateTagRules runs validate tag rules against the flag's current value.
// Uses instance-scoped validators first, falling back to the global registry.
func (r *FlagRegistry) validateTagRules(cmd *cobra.Command, tag FlagTag) error {
	if tag.Validate == "" {
		return nil
	}

	flag := r.lookupFlagForValidation(cmd, tag.Name)
	if flag == nil || !flag.Changed {
		return nil
	}

	rules, err := parseValidateRulesWithRegistry(tag.Validate, r.validators)
	if err != nil {
		return fmt.Errorf(
			"%w: validating flag %q on command %q: %w",
			ErrFlagParseFailed,
			tag.Name,
			cmd.Use,
			err,
		)
	}

	for _, rule := range rules {
		if err := rule.Validate(flag.Value.String()); err != nil {
			return fmt.Errorf(
				"%w: validating flag %q on command %q: %w",
				ErrFlagParseFailed,
				tag.Name,
				cmd.Use,
				err,
			)
		}
	}

	return nil
}
