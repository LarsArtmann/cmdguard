package v2

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// FlagRegistry manages flag registration and parsing.
type FlagRegistry struct {
	tags []FlagTag
}

// Tags returns all parsed flag tags.
func (r *FlagRegistry) Tags() []FlagTag {
	return r.tags
}

// NewFlagRegistry creates a new FlagRegistry from a config struct.
func NewFlagRegistry(cfg any) (*FlagRegistry, error) {
	tags, err := ParseFlagTags(cfg)
	if err != nil {
		return nil, err
	}

	return &FlagRegistry{tags: tags}, nil
}

// RegisterFlags adds flags to a cobra command based on the config struct.
func (r *FlagRegistry) RegisterFlags(cmd *cobra.Command) error {
	for _, tag := range r.tags {
		err := r.registerFlag(cmd, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

// registerFlag adds a single flag to the command.
func (r *FlagRegistry) registerFlag(cmd *cobra.Command, tag FlagTag) error {
	flags := cmd.PersistentFlags()

	switch tag.Type.Kind() {
	case reflect.String:
		r.addStringFlag(flags, tag)
	case reflect.Bool:
		r.addBoolFlag(flags, tag)
	case reflect.Int, reflect.Int64:
		r.addIntFlag(flags, tag)
	case reflect.Uint:
		r.addUintFlag(flags, tag)
	case reflect.Uint64:
		r.addUint64Flag(flags, tag)
	case reflect.Float64:
		r.addFloat64Flag(flags, tag)
	case reflect.Float32:
		r.addFloat32Flag(flags, tag)
	case reflect.Slice:
		r.addStringSliceFlag(flags, tag)
	default:
		// Handle custom types
		switch tag.Type {
		case reflect.TypeFor[Duration]():
			r.addDurationFlag(flags, tag)
		case reflect.TypeFor[Enum](), reflect.TypeFor[LogLevel](), reflect.TypeFor[LogFormat]():
			r.addEnumFlag(flags, tag)
		default:
			// Default to string for unknown types
			r.addStringFlag(flags, tag)
		}
	}

	return nil
}

// registerStringFlag registers a string flag with optional shorthand.
func registerStringFlag(flags *pflag.FlagSet, name, short, value, usage string) {
	if short != "" {
		_ = flags.StringP(name, short, value, usage)
	} else {
		_ = flags.String(name, value, usage)
	}
}

func (r *FlagRegistry) addStringFlag(flags *pflag.FlagSet, tag FlagTag) {
	registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)
}

func (r *FlagRegistry) addBoolFlag(flags *pflag.FlagSet, tag FlagTag) {
	def, _ := strconv.ParseBool(tag.Default)
	if tag.Short != "" {
		flags.BoolP(tag.Name, tag.Short, def, tag.Help)
	} else {
		flags.Bool(tag.Name, def, tag.Help)
	}
}

func (r *FlagRegistry) addIntFlag(flags *pflag.FlagSet, tag FlagTag) {
	def, _ := strconv.ParseInt(tag.Default, 10, 64)
	if tag.Short != "" {
		flags.IntP(tag.Name, tag.Short, int(def), tag.Help)
	} else {
		flags.Int(tag.Name, int(def), tag.Help)
	}
}

func (r *FlagRegistry) addUintFlag(flags *pflag.FlagSet, tag FlagTag) {
	def, _ := strconv.ParseUint(tag.Default, 10, 64)
	if tag.Short != "" {
		flags.UintP(tag.Name, tag.Short, uint(def), tag.Help)
	} else {
		flags.Uint(tag.Name, uint(def), tag.Help)
	}
}

func (r *FlagRegistry) addUint64Flag(flags *pflag.FlagSet, tag FlagTag) {
	def, _ := strconv.ParseUint(tag.Default, 10, 64)
	if tag.Short != "" {
		flags.Uint64P(tag.Name, tag.Short, def, tag.Help)
	} else {
		flags.Uint64(tag.Name, def, tag.Help)
	}
}

func (r *FlagRegistry) addFloat64Flag(flags *pflag.FlagSet, tag FlagTag) {
	def, _ := strconv.ParseFloat(tag.Default, 64)
	if tag.Short != "" {
		flags.Float64P(tag.Name, tag.Short, def, tag.Help)
	} else {
		flags.Float64(tag.Name, def, tag.Help)
	}
}

func (r *FlagRegistry) addFloat32Flag(flags *pflag.FlagSet, tag FlagTag) {
	def, _ := strconv.ParseFloat(tag.Default, 32)
	if tag.Short != "" {
		flags.Float32P(tag.Name, tag.Short, float32(def), tag.Help)
	} else {
		flags.Float32(tag.Name, float32(def), tag.Help)
	}
}

func (r *FlagRegistry) addStringSliceFlag(flags *pflag.FlagSet, tag FlagTag) {
	var def []string
	if tag.Default != "" {
		def = strings.Split(tag.Default, ",")
	}

	if tag.Short != "" {
		flags.StringSliceP(tag.Name, tag.Short, def, tag.Help)
	} else {
		flags.StringSlice(tag.Name, def, tag.Help)
	}
}

func (r *FlagRegistry) addDurationFlag(flags *pflag.FlagSet, tag FlagTag) {
	r.addStringFlag(flags, tag)
}

func (r *FlagRegistry) addEnumFlag(flags *pflag.FlagSet, tag FlagTag) {
	help := tag.Help
	if len(tag.Values) > 0 {
		help = fmt.Sprintf("%s (one of: %s)", tag.Help, strings.Join(tag.Values, ", "))
	}

	registerStringFlag(flags, tag.Name, tag.Short, tag.Default, help)
}

// ValidateFlags validates flag values against allowed values and checks required flags.
func (r *FlagRegistry) ValidateFlags(cmd *cobra.Command) error {
	for _, tag := range r.tags {
		err := r.validateTag(cmd, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateTag validates a single flag tag.
func (r *FlagRegistry) validateTag(cmd *cobra.Command, tag FlagTag) error {
	err := r.validateRequiredFlag(cmd, tag)
	if err != nil {
		return err
	}

	return r.validateEnumValue(cmd, tag)
}

// validateRequiredFlag checks if a required flag was set.
func (r *FlagRegistry) validateRequiredFlag(cmd *cobra.Command, tag FlagTag) error {
	if !tag.Required {
		return nil
	}

	flag := r.lookupFlagForValidation(cmd, tag.Name)
	if flag == nil || !flag.Changed {
		return NewFlagError(tag.Name, ErrRequiredFlag)
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
		return NewFlagError(tag.Name, fmt.Errorf("invalid value, must be one of: %v", tag.Values))
	}

	return nil
}

// lookupFlagForValidation finds a flag by name for validation purposes.
func (r *FlagRegistry) lookupFlagForValidation(cmd *cobra.Command, name string) *pflag.Flag {
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
