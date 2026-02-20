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

// maxEditDistance is the threshold for flag name suggestions.
const maxEditDistance = 3

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
	flag, err := r.lookupFlag(cmd, tag)
	if err != nil {
		return err
	}

	// Skip if flag wasn't changed and we're not using defaults
	if !flag.Changed && tag.Default == "" {
		return nil
	}

	value := flag.Value.String()
	return r.parseAndSetValue(cfg, tag, value)
}

// lookupFlag finds a flag in the command.
func (r *FlagRegistry) lookupFlag(cmd *cobra.Command, tag FlagTag) (*pflag.Flag, error) {
	flag := cmd.Flags().Lookup(tag.Name)
	if flag == nil {
		// Try persistent flags
		flag = cmd.PersistentFlags().Lookup(tag.Name)
	}
	if flag == nil {
		return nil, NewFlagError(tag.Name, fmt.Errorf("flag not found"))
	}
	return flag, nil
}

// parseAndSetValue parses the flag value based on type and sets it on config.
func (r *FlagRegistry) parseAndSetValue(cfg any, tag FlagTag, value string) error {
	// Parse and set the value based on type
	switch tag.Type.Kind() {
	case reflect.String:
		return SetField(cfg, tag.Field, value)
	case reflect.Bool:
		return r.parseAndSetBool(cfg, tag, value)
	case reflect.Int, reflect.Int64:
		return r.parseAndSetInt(cfg, tag, value)
	case reflect.Float64:
		return r.parseAndSetFloat64(cfg, tag, value)
	case reflect.Slice:
		return SetField(cfg, tag.Field, strings.Split(value, ","))
	default:
		return r.parseAndSetCustom(cfg, tag, value)
	}
}

// parseAndSetBool parses and sets a boolean value.
func (r *FlagRegistry) parseAndSetBool(cfg any, tag FlagTag, value string) error {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return NewFlagError(tag.Name, err)
	}
	return SetField(cfg, tag.Field, v)
}

// parseAndSetInt parses and sets an integer value.
func (r *FlagRegistry) parseAndSetInt(cfg any, tag FlagTag, value string) error {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return NewFlagError(tag.Name, err)
	}
	return SetField(cfg, tag.Field, int(v))
}

// parseAndSetFloat64 parses and sets a float64 value.
func (r *FlagRegistry) parseAndSetFloat64(cfg any, tag FlagTag, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return NewFlagError(tag.Name, err)
	}
	return SetField(cfg, tag.Field, v)
}

// parseAndSetCustom handles custom type parsing.
func (r *FlagRegistry) parseAndSetCustom(cfg any, tag FlagTag, value string) error {
	switch tag.Type {
	case reflect.TypeOf(Duration{}):
		return r.parseAndSetDuration(cfg, tag, value)
	case reflect.TypeOf(LogLevel{}):
		return r.parseAndSetLogLevel(cfg, tag, value)
	case reflect.TypeOf(LogFormat{}):
		return r.parseAndSetLogFormat(cfg, tag, value)
	case reflect.TypeOf(Enum{}):
		return r.parseAndSetEnum(cfg, tag, value)
	default:
		return SetField(cfg, tag.Field, value)
	}
}

// parseAndSetDuration parses and sets a Duration value.
func (r *FlagRegistry) parseAndSetDuration(cfg any, tag FlagTag, value string) error {
	parsed, err := ParseDuration(value)
	if err != nil {
		return NewFlagError(tag.Name, err)
	}
	return SetField(cfg, tag.Field, parsed)
}

// parseAndSetLogLevel parses and sets a LogLevel value.
func (r *FlagRegistry) parseAndSetLogLevel(cfg any, tag FlagTag, value string) error {
	parsed, err := ParseLogLevel(value)
	if err != nil {
		return NewFlagError(tag.Name, err)
	}
	return SetField(cfg, tag.Field, parsed)
}

// parseAndSetLogFormat parses and sets a LogFormat value.
func (r *FlagRegistry) parseAndSetLogFormat(cfg any, tag FlagTag, value string) error {
	parsed, err := ParseLogFormat(value)
	if err != nil {
		return NewFlagError(tag.Name, err)
	}
	return SetField(cfg, tag.Field, parsed)
}

// parseAndSetEnum parses and sets an Enum value.
func (r *FlagRegistry) parseAndSetEnum(cfg any, tag FlagTag, value string) error {
	parsed, err := ParseEnum(value, tag.Values)
	if err != nil {
		return NewFlagError(tag.Name, err)
	}
	return SetField(cfg, tag.Field, parsed)
}

// ValidateFlags validates flag values against allowed values and checks required flags.
func (r *FlagRegistry) ValidateFlags(cmd *cobra.Command) error {
	for _, tag := range r.tags {
		flag := cmd.Flags().Lookup(tag.Name)
		if flag == nil {
			flag = cmd.PersistentFlags().Lookup(tag.Name)
		}
		if flag == nil {
			continue
		}

		// Check required flags
		if tag.Required && !flag.Changed {
			return NewFlagError(tag.Name, fmt.Errorf("required flag not set"))
		}

		// Validate enum values
		if len(tag.Values) > 0 && flag.Changed {
			value := flag.Value.String()
			if !slices.Contains(tag.Values, value) {
				return NewFlagError(tag.Name, NewEnumError(value, tag.Values))
			}
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

// FlagNames returns all registered flag names for suggestion purposes.
func (r *FlagRegistry) FlagNames() []string {
	names := make([]string, len(r.tags))
	for i, tag := range r.tags {
		names[i] = tag.Name
	}
	return names
}

// SuggestFlag returns the best matching flag name for a potentially misspelled input.
// Returns empty string if no good match is found.
func SuggestFlag(validNames []string, input string) string {
	if len(validNames) == 0 {
		return ""
	}

	bestMatch := ""
	bestDist := maxEditDistance + 1

	for _, name := range validNames {
		dist := editDistance(input, name)
		if dist < bestDist {
			bestDist = dist
			bestMatch = name
		}
	}

	// Only return a match if it's close enough
	if bestDist <= maxEditDistance {
		return bestMatch
	}
	return ""
}

// editDistance computes the Levenshtein distance between two strings.
func editDistance(a, b string) int {
	aLen, bLen := len(a), len(b)
	if aLen == 0 {
		return bLen
	}
	if bLen == 0 {
		return aLen
	}

	// Use a single row for space optimization
	prev := make([]int, bLen+1)
	curr := make([]int, bLen+1)

	for j := 0; j <= bLen; j++ {
		prev[j] = j
	}

	for i := 1; i <= aLen; i++ {
		curr[0] = i
		for j := 1; j <= bLen; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(
				prev[j]+1,      // deletion
				curr[j-1]+1,    // insertion
				prev[j-1]+cost, // substitution
			)
		}
		prev, curr = curr, prev
	}

	return prev[bLen]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
