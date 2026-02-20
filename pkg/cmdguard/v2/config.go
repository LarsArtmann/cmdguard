package v2

import (
	"fmt"
	"reflect"
	"slices"
)

// Config is the default application configuration.
// Custom configs should embed or mirror this structure.
type Config struct {
	LogLevel   LogLevel  `flag:"log-level" short:"l" default:"info" help:"Log level: debug, info, warn, error"`
	LogFormat  LogFormat `flag:"log-format" default:"text" help:"Log format: text, json"`
	StrictMode bool      `flag:"strict" short:"s" default:"false" help:"Enable strict mode validation"`
	ConfigFile string    `flag:"config" short:"c" default:"" help:"Path to config file"`
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
}

// ValidateConfig validates a config struct.
// Returns all validation errors found.
func ValidateConfig(cfg any) error {
	if cfg == nil {
		return fmt.Errorf("config must not be nil")
	}

	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("config must be a struct, got %T", cfg)
	}

	return validateStruct(v, cfg)
}

// validateStruct validates all fields of a struct.
func validateStruct(v reflect.Value, cfg any) error {
	var errs []error

	tags, err := ParseFlagTags(cfg)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := validateTag(v, tag); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %v", ErrConfigValidation, errs)
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

	return nil
}

// getFieldValue extracts the string value from a field.
func getFieldValue(field reflect.Value) (string, bool) {
	switch field.Kind() {
	case reflect.String:
		return field.String(), true
	default:
		if field.Type() == reflect.TypeOf(Enum{}) ||
			field.Type() == reflect.TypeOf(LogLevel{}) ||
			field.Type() == reflect.TypeOf(LogFormat{}) {
			return field.MethodByName("String").Call(nil)[0].String(), true
		}
		return "", false
	}
}

// MergeConfigs merges multiple config sources.
// Later configs override earlier ones.
func MergeConfigs[T any](configs ...*T) *T {
	if len(configs) == 0 {
		return nil
	}

	result := configs[0]
	if result == nil {
		var zero T
		result = &zero
	}

	for _, cfg := range configs[1:] {
		if cfg == nil {
			continue
		}
		mergeStruct(reflect.ValueOf(result).Elem(), reflect.ValueOf(cfg).Elem())
	}

	return result
}

// mergeStruct merges non-zero fields from src into dst.
func mergeStruct(dst, src reflect.Value) {
	for i := 0; i < dst.NumField(); i++ {
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
