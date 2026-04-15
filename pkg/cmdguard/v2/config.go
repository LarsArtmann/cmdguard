package v2

import (
	"fmt"
	"reflect"
	"slices"
)

// Config is the default application configuration.
// Custom configs should embed or mirror this structure.
type Config struct {
	LogLevel   LogLevel  `default:"info"  flag:"log-level"  help:"Log level: debug, info, warn, error" short:"l"`
	LogFormat  LogFormat `default:"text"  flag:"log-format" help:"Log format: text, json"`
	StrictMode bool      `default:"false" flag:"strict"     help:"Enable strict mode validation"       short:"s"`
	ConfigFile string    `default:""      flag:"config"     help:"Path to config file"                 short:"c"`
}

// FlagTag represents parsed struct tag information for a flag.
type FlagTag struct {
	Name     string
	Short    string
	Default  string
	Help     string
	Values   []string // For enums
	Required bool
	Field    string
	Type     reflect.Type
	Validate string // Raw validate tag value (e.g., "min=1,max=100")
}

// ValidateConfig validates a config struct.
// Returns all validation errors found.
func ValidateConfig(cfg any) error {
	if cfg == nil {
		return ErrConfigNil
	}

	v, err := derefPointerToStruct(cfg)
	if err != nil {
		return err
	}

	return validateStruct(v, cfg)
}

// validateStruct validates all fields of a struct.
func validateStruct(v reflect.Value, cfg any) error {
	var errs []error

	tags, err := ParseFlagTags(cfg)
	if err != nil {
		return fmt.Errorf("parsing flag tags for %T: %w", cfg, err)
	}

	for _, tag := range tags {
		err := validateTag(v, tag)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("validating config %T: %w: %v", cfg, ErrConfigValidation, errs)
	}

	return nil
}

// validateTag validates a single flag tag against its field.
func validateTag(v reflect.Value, tag FlagTag) error {
	field := v.FieldByName(tag.Field)
	if !field.IsValid() {
		return nil
	}

	if len(tag.Values) > 0 {
		value, ok := getFieldValue(field)
		if !ok {
			return nil
		}

		if !slices.Contains(tag.Values, value) {
			return NewConfigError(tag.Field, NewEnumError(value, tag.Values))
		}
	}

	if tag.Validate != "" {
		err := validateFieldByKind(field, tag)
		if err != nil {
			return NewConfigError(tag.Field, err)
		}
	}

	return nil
}

// getFieldValue extracts the string value from a field.
func getFieldValue(field reflect.Value) (string, bool) {
	switch field.Kind() {
	case reflect.String:
		return field.String(), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Bool, reflect.Complex64, reflect.Complex128:
		if field.Type() == reflect.TypeFor[Enum]() ||
			field.Type() == reflect.TypeFor[LogLevel]() ||
			field.Type() == reflect.TypeFor[LogFormat]() {
			return field.MethodByName("String").Call(nil)[0].String(), true
		}

		fallthrough
	case reflect.Invalid, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
		return "", false
	default:
		return "", false
	}
}

// MergeConfigs merges multiple config sources.
// Later configs override earlier ones.
// The returned config is a deep copy; input configs are not mutated.
func MergeConfigs[T any](configs ...*T) *T {
	if len(configs) == 0 {
		return nil
	}

	var result *T

	if configs[0] == nil {
		var zero T

		result = &zero
	} else {
		result = deepCopy(configs[0])
	}

	for _, cfg := range configs[1:] {
		if cfg == nil {
			continue
		}

		mergeStruct(reflect.ValueOf(result).Elem(), reflect.ValueOf(cfg).Elem())
	}

	return result
}

// deepCopy creates a deep copy of a struct pointer via reflection.
func deepCopy[T any](src *T) *T {
	if src == nil {
		return nil
	}

	dst := *src

	return &dst
}

// mergeStruct merges non-zero fields from src into dst.
func mergeStruct(dst, src reflect.Value) {
	for i := range dst.NumField() {
		dstField := dst.Field(i)
		srcField := src.Field(i)

		if !dstField.CanSet() {
			continue
		}

		if srcField.IsZero() {
			continue
		}

		if dstField.Kind() == reflect.Struct && srcField.Kind() == reflect.Struct {
			mergeStruct(dstField, srcField)

			continue
		}

		dstField.Set(srcField)
	}
}
