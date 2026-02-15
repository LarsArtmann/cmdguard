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
		if err := r.registerFlag(cmd, tag); err != nil {
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
	case reflect.Float64:
		r.addFloat64Flag(flags, tag)
	case reflect.Slice:
		r.addStringSliceFlag(flags, tag)
	default:
		// Handle custom types
		switch tag.Type {
		case reflect.TypeOf(Duration{}):
			r.addDurationFlag(flags, tag)
		case reflect.TypeOf(Enum{}), reflect.TypeOf(LogLevel{}), reflect.TypeOf(LogFormat{}):
			r.addEnumFlag(flags, tag)
		default:
			// Default to string for unknown types
			r.addStringFlag(flags, tag)
		}
	}

	return nil
}

func (r *FlagRegistry) addStringFlag(flags *pflag.FlagSet, tag FlagTag) {
	if tag.Short != "" {
		flags.StringP(tag.Name, tag.Short, tag.Default, tag.Help)
	} else {
		flags.String(tag.Name, tag.Default, tag.Help)
	}
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

func (r *FlagRegistry) addFloat64Flag(flags *pflag.FlagSet, tag FlagTag) {
	def, _ := strconv.ParseFloat(tag.Default, 64)
	if tag.Short != "" {
		flags.Float64P(tag.Name, tag.Short, def, tag.Help)
	} else {
		flags.Float64(tag.Name, def, tag.Help)
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
	if tag.Short != "" {
		flags.StringP(tag.Name, tag.Short, tag.Default, tag.Help)
	} else {
		flags.String(tag.Name, tag.Default, tag.Help)
	}
}

func (r *FlagRegistry) addEnumFlag(flags *pflag.FlagSet, tag FlagTag) {
	help := tag.Help
	if len(tag.Values) > 0 {
		help = fmt.Sprintf("%s (one of: %s)", tag.Help, strings.Join(tag.Values, ", "))
	}
	if tag.Short != "" {
		flags.StringP(tag.Name, tag.Short, tag.Default, help)
	} else {
		flags.String(tag.Name, tag.Default, help)
	}
}

// ParseFlags populates a config struct from parsed flags.
func (r *FlagRegistry) ParseFlags(cmd *cobra.Command, cfg any) error {
	for _, tag := range r.tags {
		if err := r.parseFlag(cmd, cfg, tag); err != nil {
			return err
		}
	}
	return nil
}

// parseFlag reads a flag value and sets it on the config struct.
func (r *FlagRegistry) parseFlag(cmd *cobra.Command, cfg any, tag FlagTag) error {
	flags := cmd.Flags()

	// Check if flag was changed by user
	flag := flags.Lookup(tag.Name)
	if flag == nil {
		// Try persistent flags
		flag = cmd.PersistentFlags().Lookup(tag.Name)
	}
	if flag == nil {
		return NewFlagError(tag.Name, fmt.Errorf("flag not found"))
	}

	// Skip if flag wasn't changed and we're not using defaults
	if !flag.Changed && tag.Default == "" {
		return nil
	}

	value := flag.Value.String()

	// Parse and set the value based on type
	switch tag.Type.Kind() {
	case reflect.String:
		return SetField(cfg, tag.Field, value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return NewFlagError(tag.Name, err)
		}
		return SetField(cfg, tag.Field, v)
	case reflect.Int, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return NewFlagError(tag.Name, err)
		}
		return SetField(cfg, tag.Field, int(v))
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return NewFlagError(tag.Name, err)
		}
		return SetField(cfg, tag.Field, v)
	case reflect.Slice:
		return SetField(cfg, tag.Field, strings.Split(value, ","))
	default:
		// Handle custom types
		switch tag.Type {
		case reflect.TypeOf(Duration{}):
			parsed, err := ParseDuration(value)
			if err != nil {
				return NewFlagError(tag.Name, err)
			}
			return SetField(cfg, tag.Field, parsed)
		case reflect.TypeOf(LogLevel{}):
			parsed, err := ParseLogLevel(value)
			if err != nil {
				return NewFlagError(tag.Name, err)
			}
			return SetField(cfg, tag.Field, parsed)
		case reflect.TypeOf(LogFormat{}):
			parsed, err := ParseLogFormat(value)
			if err != nil {
				return NewFlagError(tag.Name, err)
			}
			return SetField(cfg, tag.Field, parsed)
		case reflect.TypeOf(Enum{}):
			parsed, err := ParseEnum(value, tag.Values)
			if err != nil {
				return NewFlagError(tag.Name, err)
			}
			return SetField(cfg, tag.Field, parsed)
		default:
			return SetField(cfg, tag.Field, value)
		}
	}
}

// ValidateFlags validates flag values against allowed values.
func (r *FlagRegistry) ValidateFlags(cmd *cobra.Command) error {
	for _, tag := range r.tags {
		if len(tag.Values) == 0 {
			continue
		}

		flag := cmd.Flags().Lookup(tag.Name)
		if flag == nil {
			flag = cmd.PersistentFlags().Lookup(tag.Name)
		}
		if flag == nil {
			continue
		}

		if !flag.Changed {
			continue
		}

		value := flag.Value.String()
		if !slices.Contains(tag.Values, value) {
			return NewFlagError(tag.Name, NewEnumError(value, tag.Values))
		}
	}
	return nil
}

// Tags returns all parsed flag tags.
func (r *FlagRegistry) Tags() []FlagTag {
	return r.tags
}

// GenerateHelp generates help text for all flags.
func (r *FlagRegistry) GenerateHelp() string {
	var lines []string
	for _, tag := range r.tags {
		line := fmt.Sprintf("  --%s", tag.Name)
		if tag.Short != "" {
			line += fmt.Sprintf(", -%s", tag.Short)
		}
		line += fmt.Sprintf("\t%s", tag.Help)
		if tag.Default != "" {
			line += fmt.Sprintf(" (default: %s)", tag.Default)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
