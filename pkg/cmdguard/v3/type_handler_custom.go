package v3

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

func (r *typeRegistry) registerCustomTypes() {
	r.registerEnumTypes()
	r.registerValueTypes()
}

func (r *typeRegistry) registerEnumTypes() {
	enumHelp := func(tag FlagTag) string {
		if len(tag.Values) > 0 {
			return fmt.Sprintf("%s (one of: %s)", tag.Help, strings.Join(tag.Values, ", "))
		}

		return tag.Help
	}

	enumHandler := TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, enumHelp(tag))

			return nil
		},
		ParseFunc: func(value string, tag FlagTag) (any, error) {
			return ParseEnum(value, tag.Values)
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}
	r.byType[reflect.TypeFor[Enum]()] = enumHandler

	makeEnumLikeHandler := func(
		parseFunc func(string) (any, error),
		defaultAllowed []string,
	) TypeHandlerFunc {
		return TypeHandlerFunc{
			RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
				allowed := tag.Values
				if len(allowed) == 0 {
					allowed = defaultAllowed
				}

				help := fmt.Sprintf("%s (one of: %s)", tag.Help, strings.Join(allowed, ", "))
				registerStringFlag(flags, tag.Name, tag.Short, tag.Default, help)

				return nil
			},
			ParseFunc: func(value string, _ FlagTag) (any, error) {
				return parseFunc(value)
			},
			DefaultFunc: func(tag FlagTag) any {
				return tag.Default
			},
		}
	}

	r.byType[reflect.TypeFor[LogLevel]()] = makeEnumLikeHandler(
		func(v string) (any, error) { return ParseLogLevel(v) },
		logLevelAllowed,
	)
	r.byType[reflect.TypeFor[LogFormat]()] = makeEnumLikeHandler(
		func(v string) (any, error) { return ParseLogFormat(v) },
		logFormatAllowed,
	)
}

func (r *typeRegistry) registerValueTypes() {
	for _, entry := range []struct {
		typ        reflect.Type
		parse      func(string) (any, error)
		defaultVal func(FlagTag) any
	}{
		{reflect.TypeFor[Duration](), func(v string) (any, error) { return ParseDuration(v) }, func(tag FlagTag) any {
			d, err := ParseDuration(tag.Default)
			if err != nil {
				return Duration{}
			}

			return d
		}},
		{reflect.TypeFor[URL](), func(v string) (any, error) { return ParseURL(v) }, stringDefault},
		{reflect.TypeFor[Email](), func(v string) (any, error) { return ParseEmail(v) }, stringDefault},
		{reflect.TypeFor[Port](), func(v string) (any, error) { return ParsePort(v) }, stringDefault},
		{reflect.TypeFor[FilePath](), func(v string) (any, error) { return ParseFilePath(v, false) }, stringDefault},
		{reflect.TypeFor[HostPort](), func(v string) (any, error) { return ParseHostPort(v) }, stringDefault},
	} {
		typ, parse, def := entry.typ, entry.parse, entry.defaultVal
		r.byType[typ] = TypeHandlerFunc{
			RegisterFunc: registerStringFlagFromTag,
			ParseFunc: func(value string, _ FlagTag) (any, error) {
				return parse(value)
			},
			DefaultFunc: def,
		}
	}
}

func stringDefault(tag FlagTag) any {
	return tag.Default
}

// RegisterGoDurationHandler registers a TypeHandler for time.Duration fields
// in the global defaults template. New FlagRegistries will include this handler.
func RegisterGoDurationHandler() {
	globalTypeRegistry.registerGoDurationHandler()
}

// RegisterGoDurationHandler registers a TypeHandler for time.Duration fields
// on this registry instance.
func (r *FlagRegistry) RegisterGoDurationHandler() {
	r.types.registerGoDurationHandler()
}

func (r *typeRegistry) registerGoDurationHandler() {
	r.register(reflect.TypeFor[time.Duration](), TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			def, err := time.ParseDuration(tag.Default)
			if err != nil && tag.Default != "" {
				return fmt.Errorf("invalid duration default for flag %q: %w", tag.Name, err)
			}

			if tag.Short != "" {
				flags.DurationP(tag.Name, tag.Short, def, tag.Help)
			} else {
				flags.Duration(tag.Name, def, tag.Help)
			}

			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return time.ParseDuration(value)
		},
		DefaultFunc: func(tag FlagTag) any {
			d, err := time.ParseDuration(tag.Default)
			if err != nil {
				return time.Duration(0)
			}

			return d
		},
	})
}
