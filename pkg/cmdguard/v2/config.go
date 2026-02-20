package v2

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
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

// ParseFlagTags extracts flag information from a config struct.
// The struct must have `flag` tags on its fields.
func ParseFlagTags(cfg any) ([]FlagTag, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config must not be nil")
	}

	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("config must be a struct, got %T", cfg)
	}

	t := v.Type()
	var tags []FlagTag

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		flagTag := field.Tag.Get("flag")

		// Skip fields without flag tag
		if flagTag == "" || flagTag == "-" {
			continue
		}

		tag := FlagTag{
			Field: field.Name,
			Type:  field.Type,
		}

		// Parse flag name
		tag.Name = flagTag

		// Parse short form
		if short := field.Tag.Get("short"); short != "" {
			tag.Short = short
		}

		// Parse default value
		if def := field.Tag.Get("default"); def != "" {
			tag.Default = def
		}

		// Parse help text
		if help := field.Tag.Get("help"); help != "" {
			tag.Help = help
		}

		// Parse allowed values (for enums)
		if values := field.Tag.Get("values"); values != "" {
			tag.Values = strings.Split(values, ",")
		}

		// Parse required tag
		if req := field.Tag.Get("required"); req != "" {
			tag.Required, _ = strconv.ParseBool(req)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

// DefaultValue returns the default value for a flag based on its type.
func (t FlagTag) DefaultValue() any {
	if t.Default == "" {
		return reflect.Zero(t.Type).Interface()
	}

	switch t.Type.Kind() {
	case reflect.String:
		return t.Default
	case reflect.Bool:
		v, _ := strconv.ParseBool(t.Default)
		return v
	case reflect.Int, reflect.Int64:
		// Check if it's a Duration type
		if t.Type == reflect.TypeOf(Duration{}) {
			d, err := ParseDuration(t.Default)
			if err != nil {
				return Duration{}
			}
			return d
		}
		v, _ := strconv.ParseInt(t.Default, 10, 64)
		return int(v)
	case reflect.Float64:
		v, _ := strconv.ParseFloat(t.Default, 64)
		return v
	case reflect.Slice:
		// For slices, parse comma-separated values
		return strings.Split(t.Default, ",")
	default:
		// Handle custom types (Enum, Duration, LogLevel, LogFormat)
		switch t.Type {
		case reflect.TypeOf(Duration{}):
			d, err := ParseDuration(t.Default)
			if err != nil {
				return Duration{}
			}
			return d
		case reflect.TypeOf(Enum{}), reflect.TypeOf(LogLevel{}), reflect.TypeOf(LogFormat{}):
			return t.Default
		default:
			return t.Default
		}
	}
}

// SetField sets a field value on a config struct using reflection.
func SetField(cfg any, fieldName string, value any) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to struct")
	}

	v = v.Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %q not found", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("field %q is not settable", fieldName)
	}

	val := reflect.ValueOf(value)

	// Handle type conversions
	if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))
		return nil
	}

	// Handle string to custom type conversions
	if val.Kind() == reflect.String {
		str := val.String()
		switch field.Type() {
		case reflect.TypeOf(LogLevel{}):
			parsed, err := ParseLogLevel(str)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(parsed))
			return nil
		case reflect.TypeOf(LogFormat{}):
			parsed, err := ParseLogFormat(str)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(parsed))
			return nil
		case reflect.TypeOf(Duration{}):
			parsed, err := ParseDuration(str)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(parsed))
			return nil
		case reflect.TypeOf(Enum{}):
			// For generic enums, just set the string value
			field.Set(reflect.ValueOf(Enum{value: str}))
			return nil
		}
	}

	// Handle time.Duration to Duration conversion
	if val.Type() == reflect.TypeOf(time.Duration(0)) && field.Type() == reflect.TypeOf(Duration{}) {
		field.Set(reflect.ValueOf(FromDuration(val.Interface().(time.Duration))))
		return nil
	}

	return fmt.Errorf("cannot convert %T to %s", value, field.Type())
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

	var errs []error

	tags, err := ParseFlagTags(cfg)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		field := v.FieldByName(tag.Field)
		if !field.IsValid() {
			continue
		}

		// Validate enum values
		if len(tag.Values) > 0 {
			var value string
			switch field.Kind() {
			case reflect.String:
				value = field.String()
			default:
				// Handle custom enum types
				if field.Type() == reflect.TypeOf(Enum{}) ||
					field.Type() == reflect.TypeOf(LogLevel{}) ||
					field.Type() == reflect.TypeOf(LogFormat{}) {
					value = field.MethodByName("String").Call(nil)[0].String()
				} else {
					continue
				}
			}

			if !slices.Contains(tag.Values, value) {
				errs = append(errs, NewConfigError(tag.Field, NewEnumError(value, tag.Values)))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %v", ErrConfigValidation, errs)
	}

	return nil
}

// MergeConfigs merges multiple config sources.
// Later configs override earlier ones.
func MergeConfigs[T any](configs ...*T) *T {
	if len(configs) == 0 {
		return nil
	}

	// Start with the first config
	result := configs[0]
	if result == nil {
		var zero T
		result = &zero
	}

	// Merge subsequent configs
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

		// Skip zero values in source
		if srcField.IsZero() {
			continue
		}

		// Recursively merge nested structs
		if dstField.Kind() == reflect.Struct && srcField.Kind() == reflect.Struct {
			mergeStruct(dstField, srcField)
			continue
		}

		// Copy non-zero source to destination
		dstField.Set(srcField)
	}
}
