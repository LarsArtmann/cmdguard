package v2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ParseFlags populates a config struct from parsed flags.
func (r *FlagRegistry) ParseFlags(cmd *cobra.Command, cfg any) error {
	for _, tag := range r.tags {
		err := r.parseFlag(cmd, cfg, tag)
		if err != nil {
			return fmt.Errorf("parsing flags for command %q: %w", cmd.Use, err)
		}
	}

	return nil
}

// parseFlag reads a flag value and sets it on the config struct.
func (r *FlagRegistry) parseFlag(cmd *cobra.Command, cfg any, tag FlagTag) error {
	flag, err := r.lookupFlag(cmd, tag)
	if err != nil {
		return fmt.Errorf("looking up flag %q on command %q: %w", tag.Name, cmd.Use, err)
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
		return nil, fmt.Errorf(
			"flag %q not found in command %q: %w",
			tag.Name,
			cmd.Use,
			ErrFlagNotFound,
		)
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return r.parseAndSetInt(cfg, tag, value)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return r.parseAndSetUint(cfg, tag, value)
	case reflect.Float32, reflect.Float64:
		return r.parseAndSetFloat64(cfg, tag, value)
	case reflect.Slice:
		return SetField(cfg, tag.Field, strings.Split(value, ","))
	case reflect.Invalid, reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Struct, reflect.UnsafePointer:
		return r.parseAndSetCustom(cfg, tag, value)
	default:
		return r.parseAndSetCustom(cfg, tag, value)
	}
}

// parseFlagValue is a helper for parsing values with error handling.
func parseFlagValue(
	cfg any,
	tag FlagTag,
	value string,
	typeName string,
	parser func(string) (any, error),
) error {
	v, err := parser(value)
	if err != nil {
		return fmt.Errorf("parsing %s flag %q with value %q: %w", typeName, tag.Name, value, err)
	}

	return SetField(cfg, tag.Field, v)
}

// parseAndSetBool parses and sets a boolean value.
func (r *FlagRegistry) parseAndSetBool(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "bool", func(v string) (any, error) {
		return strconv.ParseBool(v)
	})
}

// parseAndSetInt parses and sets an integer value.
func (r *FlagRegistry) parseAndSetInt(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "int", func(v string) (any, error) {
		parsed, err := strconv.ParseInt(v, 10, 64)
		return int(parsed), err
	})
}

// parseAndSetUint parses and sets an unsigned integer value.
func (r *FlagRegistry) parseAndSetUint(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "uint", func(v string) (any, error) {
		parsed, err := strconv.ParseUint(v, 10, 64)
		return uint(parsed), err
	})
}

func (r *FlagRegistry) parseAndSetFloat64(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "float64", func(v string) (any, error) {
		return strconv.ParseFloat(v, 64)
	})
}

// parseAndSetCustom handles custom type parsing.
func (r *FlagRegistry) parseAndSetCustom(cfg any, tag FlagTag, value string) error {
	switch tag.Type {
	case reflect.TypeFor[Duration]():
		return r.parseAndSetDuration(cfg, tag, value)
	case reflect.TypeFor[LogLevel]():
		return r.parseAndSetLogLevel(cfg, tag, value)
	case reflect.TypeFor[LogFormat]():
		return r.parseAndSetLogFormat(cfg, tag, value)
	case reflect.TypeFor[Enum]():
		return r.parseAndSetEnum(cfg, tag, value)
	case reflect.TypeFor[URL]():
		return r.parseAndSetURL(cfg, tag, value)
	case reflect.TypeFor[Email]():
		return r.parseAndSetEmail(cfg, tag, value)
	case reflect.TypeFor[Port]():
		return r.parseAndSetPort(cfg, tag, value)
	case reflect.TypeFor[FilePath]():
		return r.parseAndSetFilePath(cfg, tag, value)
	case reflect.TypeFor[HostPort]():
		return r.parseAndSetHostPort(cfg, tag, value)
	default:
		return SetField(cfg, tag.Field, value)
	}
}

// parseAndSetDuration parses and sets a Duration value.
func (r *FlagRegistry) parseAndSetDuration(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "duration", func(v string) (any, error) {
		return ParseDuration(v)
	})
}

// parseAndSetLogLevel parses and sets a LogLevel value.
func (r *FlagRegistry) parseAndSetLogLevel(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "log level", func(v string) (any, error) {
		return ParseLogLevel(v)
	})
}

// parseAndSetLogFormat parses and sets a LogFormat value.
func (r *FlagRegistry) parseAndSetLogFormat(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "log format", func(v string) (any, error) {
		return ParseLogFormat(v)
	})
}

// parseAndSetEnum parses and sets an Enum value.
func (r *FlagRegistry) parseAndSetEnum(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "enum", func(v string) (any, error) {
		return ParseEnum(v, tag.Values)
	})
}

// parseAndSetURL parses and sets a URL value.
func (r *FlagRegistry) parseAndSetURL(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "URL", func(v string) (any, error) {
		return ParseURL(v)
	})
}

// parseAndSetEmail parses and sets an Email value.
func (r *FlagRegistry) parseAndSetEmail(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "email", func(v string) (any, error) {
		return ParseEmail(v)
	})
}

// parseAndSetPort parses and sets a Port value.
func (r *FlagRegistry) parseAndSetPort(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "port", func(v string) (any, error) {
		return ParsePort(v)
	})
}

// parseAndSetFilePath parses and sets a FilePath value.
func (r *FlagRegistry) parseAndSetFilePath(cfg any, tag FlagTag, value string) error {
	// Note: FilePath parsing does NOT check if the path exists
	return parseFlagValue(cfg, tag, value, "file path", func(v string) (any, error) {
		return ParseFilePath(v, false)
	})
}

// parseAndSetHostPort parses and sets a HostPort value.
func (r *FlagRegistry) parseAndSetHostPort(cfg any, tag FlagTag, value string) error {
	return parseFlagValue(cfg, tag, value, "host:port", func(v string) (any, error) {
		return ParseHostPort(v)
	})
}
