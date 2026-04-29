package v2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

// TypeHandler defines how a flag type is registered, parsed, defaulted, and set.
// Each custom type or primitive kind can have its own handler.
// The registry eliminates the 3-way switch dispatch across flags.go,
// flags_parse.go, config_setfield.go, and config_parsing.go.
type TypeHandler interface {
	// Register adds the flag to the given pflag.FlagSet.
	Register(flags *pflag.FlagSet, tag FlagTag) error
	// Parse converts a string value to the appropriate Go value.
	Parse(value string, tag FlagTag) (any, error)
	// Default returns the default value for this type given the tag.
	Default(tag FlagTag) any
}

// registerStringFlag registers a string flag with optional shorthand.
func registerStringFlag(flags *pflag.FlagSet, name, short, value, usage string) {
	if short != "" {
		_ = flags.StringP(name, short, value, usage)
	} else {
		_ = flags.String(name, value, usage)
	}
}

// TypeHandlerFunc is a functional adapter for TypeHandler where Register is not needed.
type TypeHandlerFunc struct {
	ParseFunc    func(value string, tag FlagTag) (any, error)
	DefaultFunc  func(tag FlagTag) any
	RegisterFunc func(flags *pflag.FlagSet, tag FlagTag) error
}

func (h TypeHandlerFunc) Register(flags *pflag.FlagSet, tag FlagTag) error {
	if h.RegisterFunc != nil {
		return h.RegisterFunc(flags, tag)
	}
	return nil
}

func (h TypeHandlerFunc) Parse(value string, tag FlagTag) (any, error) {
	return h.ParseFunc(value, tag)
}

func (h TypeHandlerFunc) Default(tag FlagTag) any {
	return h.DefaultFunc(tag)
}

// typeRegistry maps reflect.Type or reflect.Kind to TypeHandlers.
type typeRegistry struct {
	byType map[reflect.Type]TypeHandler
	byKind map[reflect.Kind]TypeHandler
}

// globalTypeRegistry is the default registry with all built-in types.
var globalTypeRegistry = newTypeRegistry()

func newTypeRegistry() *typeRegistry {
	r := &typeRegistry{
		byType: make(map[reflect.Type]TypeHandler),
		byKind: make(map[reflect.Kind]TypeHandler),
	}

	r.registerKinds()
	r.registerCustomTypes()

	return r
}

func (r *typeRegistry) registerKinds() {
	r.byKind[reflect.String] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return value, nil
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}

	r.byKind[reflect.Bool] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			def, err := parseBoolDefault(tag.Default)
			if err != nil {
				return fmt.Errorf("invalid bool default for flag %q: %w", tag.Name, err)
			}
			if tag.Short != "" {
				flags.BoolP(tag.Name, tag.Short, def, tag.Help)
			} else {
				flags.Bool(tag.Name, def, tag.Help)
			}
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return strconv.ParseBool(value)
		},
		DefaultFunc: func(tag FlagTag) any {
			v, _ := parseBoolDefault(tag.Default)
			return v
		},
	}

	intHandler := TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			def, err := parseIntDefault(tag.Default)
			if err != nil {
				return fmt.Errorf("invalid int default for flag %q: %w", tag.Name, err)
			}
			if tag.Short != "" {
				flags.IntP(tag.Name, tag.Short, int(def), tag.Help)
			} else {
				flags.Int(tag.Name, int(def), tag.Help)
			}
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			parsed, err := strconv.ParseInt(value, 10, 64)
			return parsed, err
		},
		DefaultFunc: func(tag FlagTag) any {
			v, _ := parseIntDefault(tag.Default)
			return v
		},
	}
	for _, k := range []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64} {
		r.byKind[k] = intHandler
	}

	uintHandler := TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			def, err := parseUintDefault(tag.Default)
			if err != nil {
				return fmt.Errorf("invalid uint default for flag %q: %w", tag.Name, err)
			}
			if tag.Short != "" {
				flags.UintP(tag.Name, tag.Short, uint(def), tag.Help)
			} else {
				flags.Uint(tag.Name, uint(def), tag.Help)
			}
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			parsed, err := strconv.ParseUint(value, 10, 64)
			return parsed, err
		},
		DefaultFunc: func(tag FlagTag) any {
			v, _ := parseUintDefault(tag.Default)
			return v
		},
	}
	for _, k := range []reflect.Kind{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr} {
		r.byKind[k] = uintHandler
	}

	floatHandler := TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			def, err := parseFloat64Default(tag.Default)
			if err != nil {
				return fmt.Errorf("invalid float64 default for flag %q: %w", tag.Name, err)
			}
			if tag.Short != "" {
				flags.Float64P(tag.Name, tag.Short, def, tag.Help)
			} else {
				flags.Float64(tag.Name, def, tag.Help)
			}
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return strconv.ParseFloat(value, 64)
		},
		DefaultFunc: func(tag FlagTag) any {
			v, _ := parseFloat64Default(tag.Default)
			return v
		},
	}
	r.byKind[reflect.Float32] = floatHandler
	r.byKind[reflect.Float64] = floatHandler

	r.byKind[reflect.Slice] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			var def []string
			if tag.Default != "" {
				def = strings.Split(tag.Default, ",")
			}
			if tag.Short != "" {
				flags.StringSliceP(tag.Name, tag.Short, def, tag.Help)
			} else {
				flags.StringSlice(tag.Name, def, tag.Help)
			}
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return strings.Split(value, ","), nil
		},
		DefaultFunc: func(tag FlagTag) any {
			if tag.Default == "" {
				return []string(nil)
			}
			return strings.Split(tag.Default, ",")
		},
	}
}

func (r *typeRegistry) registerCustomTypes() {
	enumHelp := func(tag FlagTag) string {
		help := tag.Help
		if len(tag.Values) > 0 {
			help = fmt.Sprintf("%s (one of: %s)", tag.Help, strings.Join(tag.Values, ", "))
		}
		return help
	}

	r.byType[reflect.TypeFor[Duration]()] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return ParseDuration(value)
		},
		DefaultFunc: func(tag FlagTag) any {
			d, err := ParseDuration(tag.Default)
			if err != nil {
				return Duration{}
			}
			return d
		},
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
	r.byType[reflect.TypeFor[LogLevel]()] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			help := tag.Help
			if len(tag.Values) > 0 {
				help = fmt.Sprintf("%s (one of: %s)", tag.Help, strings.Join(tag.Values, ", "))
			} else {
				help = fmt.Sprintf("%s (one of: %s)", tag.Help, strings.Join(logLevelAllowed, ", "))
			}
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, help)
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return ParseLogLevel(value)
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}
	r.byType[reflect.TypeFor[LogFormat]()] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			help := tag.Help
			if len(tag.Values) > 0 {
				help = fmt.Sprintf("%s (one of: %s)", tag.Help, strings.Join(tag.Values, ", "))
			} else {
				help = fmt.Sprintf("%s (one of: %s)", tag.Help, strings.Join(logFormatAllowed, ", "))
			}
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, help)
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return ParseLogFormat(value)
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}

	r.byType[reflect.TypeFor[URL]()] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return ParseURL(value)
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}

	r.byType[reflect.TypeFor[Email]()] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return ParseEmail(value)
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}

	r.byType[reflect.TypeFor[Port]()] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return ParsePort(value)
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}

	r.byType[reflect.TypeFor[FilePath]()] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return ParseFilePath(value, false)
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}

	r.byType[reflect.TypeFor[HostPort]()] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)
			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return ParseHostPort(value)
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}
}

// lookupHandler finds the TypeHandler for a given reflect.Type.
// Checks exact type match first, then falls back to kind-based lookup.
func (r *typeRegistry) lookupHandler(typ reflect.Type) (TypeHandler, bool) {
	if h, ok := r.byType[typ]; ok {
		return h, true
	}
	if h, ok := r.byKind[typ.Kind()]; ok {
		return h, true
	}
	return nil, false
}

// RegisterTypeHandler registers a custom TypeHandler for a specific reflect.Type.
// This allows users to extend the flag system with their own custom types.
func RegisterTypeHandler(typ reflect.Type, handler TypeHandler) {
	globalTypeRegistry.byType[typ] = handler
}

// handledByTypeRegistry checks whether the given type has a handler in the registry.
func handledByTypeRegistry(typ reflect.Type) bool {
	_, ok := globalTypeRegistry.lookupHandler(typ)
	return ok
}

// dispatchRegister dispatches flag registration to the TypeHandler registry.
func dispatchRegister(flags *pflag.FlagSet, tag FlagTag) error {
	h, ok := globalTypeRegistry.lookupHandler(tag.Type)
	if !ok {
		registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)
		return nil
	}
	return h.Register(flags, tag)
}

// dispatchParse dispatches value parsing to the TypeHandler registry.
func dispatchParse(value string, tag FlagTag) (any, error) {
	h, ok := globalTypeRegistry.lookupHandler(tag.Type)
	if !ok {
		return value, nil
	}
	return h.Parse(value, tag)
}

// dispatchDefault dispatches default value computation to the TypeHandler registry.
func dispatchDefault(tag FlagTag) any {
	if tag.Default == "" {
		return reflect.Zero(tag.Type).Interface()
	}
	h, ok := globalTypeRegistry.lookupHandler(tag.Type)
	if !ok {
		return tag.Default
	}
	def := h.Default(tag)
	// Convert to the exact field type for numeric kinds
	defVal := reflect.ValueOf(def)
	if defVal.IsValid() && defVal.Type().ConvertibleTo(tag.Type) {
		return defVal.Convert(tag.Type).Interface()
	}
	return def
}

// RegisterGoDurationHandler registers a TypeHandler for time.Duration fields.
// This allows using Go's native time.Duration type directly in flag structs.
func RegisterGoDurationHandler() {
	globalTypeRegistry.byType[reflect.TypeFor[time.Duration]()] = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			def, _ := time.ParseDuration(tag.Default)
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
	}
}
