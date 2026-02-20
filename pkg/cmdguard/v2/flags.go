package v2

import (
	"fmt"
	"reflect"
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

// ValidateFlags validates flag values against allowed values and checks required flags.
func (r *FlagRegistry) ValidateFlags(cmd *cobra.Command) error {
	for _, tag := range r.tags {
		// Check required flags
		if tag.Required {
			flag := cmd.Flags().Lookup(tag.Name)
			if flag == nil {
				flag = cmd.PersistentFlags().Lookup(tag.Name)
			}
			if flag == nil || !flag.Changed {
				return NewFlagError(tag.Name, fmt.Errorf("required flag not set"))
			}
		}

		// Validate enum values
		if len(tag.Values) > 0 {
			flag := cmd.Flags().Lookup(tag.Name)
			if flag == nil {
				flag = cmd.PersistentFlags().Lookup(tag.Name)
			}
			if flag != nil && flag.Changed {
				value := flag.Value.String()
				found := false
				for _, allowed := range tag.Values {
					if value == allowed {
						found = true
						break
					}
				}
				if !found {
					return NewFlagError(tag.Name, fmt.Errorf("invalid value %q, must be one of: %v", value, tag.Values))
				}
			}
		}
	}
	return nil
}

