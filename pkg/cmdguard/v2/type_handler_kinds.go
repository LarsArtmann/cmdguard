package v2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

func (r *typeRegistry) registerKinds() {
	countHandler := TypeHandlerFunc{
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
	r.countHandler = countHandler

	r.byKind[reflect.String] = TypeHandlerFunc{
		RegisterFunc: registerStringFlagFromTag,
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
