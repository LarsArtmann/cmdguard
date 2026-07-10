package v3

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

func (r *typeRegistry) registerKinds() {
	r.registerCountHandler()
	r.registerStringKind()
	r.registerBoolKind()
	r.registerIntKinds()
	r.registerUintKinds()
	r.registerFloatKinds()
	r.registerSliceKind()
}

func (r *typeRegistry) registerCountHandler() {
	r.countHandler = TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			if tag.Short != "" {
				flags.CountP(tag.Name, tag.Short, tag.Help)
			} else {
				flags.Count(tag.Name, tag.Help)
			}

			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return strconv.Atoi(value)
		},
		DefaultFunc: func(_ FlagTag) any {
			return 0
		},
	}
}

func (r *typeRegistry) registerStringKind() {
	r.byKind[reflect.String] = TypeHandlerFunc{
		RegisterFunc: registerStringFlagFromTag,
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return value, nil
		},
		DefaultFunc: func(tag FlagTag) any {
			return tag.Default
		},
	}
}

func (r *typeRegistry) registerBoolKind() {
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
}

func (r *typeRegistry) registerIntKinds() {
	// Integer handlers are built per-kind so that flag values are validated
	// against the field's actual bit-width.
	for _, k := range []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64} {
		r.byKind[k] = makeIntKindHandler(intBitSize(k))
	}
}

func (r *typeRegistry) registerUintKinds() {
	for _, k := range []reflect.Kind{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr} {
		r.byKind[k] = makeUintKindHandler(uintBitSize(k))
	}
}

func (r *typeRegistry) registerFloatKinds() {
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
}

func (r *typeRegistry) registerSliceKind() {
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
