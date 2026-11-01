package v4

import (
	"fmt"
	"iter"
	"os"
	"reflect"
	"slices"
	"strings"

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

		err := r.registerFlag(flagSet, tag)
		if err != nil {
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

		err := r.registerFlag(cmd.Flags(), tag)
		if err != nil {
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
	err := dispatchRegister(r.types, flags, tag)
	if err != nil {
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

// validateTag validates a single flag tag. Each check returns an error that
// already names the flag and wraps its sentinel, so this function adds no
// further wrapping: nested prefixes make user-facing errors unreadable.
func (r *FlagRegistry) validateTag(cmd *cobra.Command, tag FlagTag) error {
	err := r.validateRequiredFlag(cmd, tag)
	if err != nil {
		return err
	}

	err = r.validateEnumValue(cmd, tag)
	if err != nil {
		return err
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

	if flag != nil && flag.Changed {
		return nil
	}

	// An environment variable is an explicit, supported source for a flag
	// (flag > env > config > default), so it satisfies "required".
	if _, ok := r.envValue(tag); ok {
		return nil
	}

	return fmt.Errorf("%w: %q is missing", ErrRequiredFlag, tag.Name)
}

// validateEnumValue validates that an enum flag has an allowed value.
func (r *FlagRegistry) validateEnumValue(cmd *cobra.Command, tag FlagTag) error {
	if len(tag.Values) == 0 {
		return nil
	}

	value, present := r.resolvedValue(cmd, tag)
	if !present {
		return nil
	}

	if !r.isAllowedValue(value, tag.Values) {
		return fmt.Errorf("%w: %q = %q, allowed: %v", ErrInvalidEnum, tag.Name, value, tag.Values)
	}

	return nil
}

// resolvedValue returns the value a flag resolves to at run time, using the
// same precedence as parsing: explicit flag > environment > tag default.
// Validation operates on this value, so env- and default-sourced values are
// checked too, not only values typed on the command line.
func (r *FlagRegistry) resolvedValue(cmd *cobra.Command, tag FlagTag) (string, bool) {
	flag := r.lookupFlagForValidation(cmd, tag.Name)
	if flag != nil && flag.Changed {
		return flag.Value.String(), true
	}

	if value, ok := r.envValue(tag); ok {
		return value, true
	}

	if tag.Default != "" {
		return tag.Default, true
	}

	return "", false
}

// envValue returns the environment value for tag, applying the registry's
// prefix exactly as flag parsing does.
func (r *FlagRegistry) envValue(tag FlagTag) (string, bool) {
	if tag.Env == "" {
		return "", false
	}

	name := tag.Env
	if r.envPrefix != "" {
		name = r.envPrefix + name
	}

	return os.LookupEnv(name)
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

	expr := tag.Validate
	if tag.Type != nil && tag.Type.Kind() == reflect.String {
		// min/max on a string field mean length; using the numeric validators
		// would try to parse the text as a number and fail confusingly.
		expr = stringLengthRules(expr)
	}

	rules, err := parseValidateRulesWithRegistry(expr, r.validators)
	if err != nil {
		return fmt.Errorf("flag %q: %w", tag.Name, err)
	}

	value, present := r.resolvedValue(cmd, tag)
	if !present {
		return nil
	}

	for _, rule := range rules {
		err := rule.Validate(value)
		if err != nil {
			return fmt.Errorf("flag %q: %w", tag.Name, err)
		}
	}

	return nil
}

// unreachable placeholder
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

// skipsFlagValidation reports whether a command is an introspection or
// maintenance command that must stay runnable even when business flags are
// required (agents run describe/version/doctor before they know the flags).
func skipsFlagValidation(name string) bool {
	switch name {
	case describeCommandName, "version", "doctor", "completion", helpCommandName:
		return true
	default:
		return false
	}
}

// stringLengthRules rewrites numeric bounds to their length equivalents for a
// string field, so validate:"min=3" means "at least 3 characters".
func stringLengthRules(expr string) string {
	parts := strings.Split(expr, ",")

	for i, part := range parts {
		name, param, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}

		switch strings.TrimSpace(name) {
		case "min":
			parts[i] = "minlen=" + param
		case "max":
			parts[i] = "maxlen=" + param
		}
	}

	return strings.Join(parts, ",")
}
