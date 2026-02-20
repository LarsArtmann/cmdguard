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

	return parseStructTags(v.Type())
}

// parseStructTags parses all flag tags from a struct type.
func parseStructTags(t reflect.Type) ([]FlagTag, error) {
	var tags []FlagTag

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := parseFieldTag(field)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}

	return tags, nil
}

// parseFieldTag parses flag tags from a single struct field.
// Returns nil if the field doesn't have a flag tag.
func parseFieldTag(field reflect.StructField) *FlagTag {
	flagTag := field.Tag.Get("flag")

	// Skip fields without flag tag
	if flagTag == "" || flagTag == "-" {
		return nil
	}

	tag := FlagTag{
		Field: field.Name,
		Type:  field.Type,
		Name:  flagTag,
	}

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

	return &tag
}

// DefaultValue returns the default value for a flag based on its type.
func (t FlagTag) DefaultValue() any {
	if t.Default == "" {
		return reflect.Zero(t.Type).Interface()
	}

	return t.parseDefaultValue()
}

// parseDefaultValue parses the default value based on type.
func (t FlagTag) parseDefaultValue() any {
	switch t.Type.Kind() {
	case reflect.String:
		return t.Default
	case reflect.Bool:
		return parseBoolDefault(t.Default)
	case reflect.Int, reflect.Int64:
		return t.parseIntDefault()
	case reflect.Float64:
		return parseFloat64Default(t.Default)
	case reflect.Slice:
		return strings.Split(t.Default, ",")
	default:
		return t.parseCustomDefault()
	}
}

// parseBoolDefault parses a boolean default value.
func parseBoolDefault(s string) bool {
	v, _ := strconv.ParseBool(s)
	return v
}

// parseIntDefault parses an integer default value.
func (t FlagTag) parseIntDefault() any {
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
}

// parseFloat64Default parses a float64 default value.
func parseFloat64Default(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// parseCustomDefault handles custom type defaults.
func (t FlagTag) parseCustomDefault() any {
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

// SetField sets a field value on a config struct using reflection.
func SetField(cfg any, fieldName string, value any) error {
	field, err := getField(cfg, fieldName)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(value)

	// Handle type conversions
	if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))
		return nil
	}

	// Handle string to custom type conversions
	if val.Kind() == reflect.String {
		if err := setStringField(field, val.String()); err != nil {
			return err
		}
		return nil
	}

	// Handle time.Duration to Duration conversion
	if val.Type() == reflect.TypeOf(time.Duration(0)) && field.Type() == reflect.TypeOf(Duration{}) {
		field.Set(reflect.ValueOf(FromDuration(val.Interface().(time.Duration))))
		return nil
	}

	return fmt.Errorf("cannot convert %T to %s", value, field.Type())
}

// getField retrieves a field from config by name.
func getField(cfg any, fieldName string) (reflect.Value, error) {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("config must be a pointer to struct")
	}

	v = v.Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return reflect.Value{}, fmt.Errorf("field %q not found", fieldName)
	}

	if !field.CanSet() {
		return reflect.Value{}, fmt.Errorf("field %q is not settable", fieldName)
	}

	return field, nil
}

// setStringField handles string to custom type conversions.
func setStringField(field reflect.Value, str string) error {
	switch field.Type() {
	case reflect.TypeOf(LogLevel{}):
		return parseAndSetLogLevel(field, str)
	case reflect.TypeOf(LogFormat{}):
		return parseAndSetLogFormat(field, str)
	case reflect.TypeOf(Duration{}):
		return parseAndSetDuration(field, str)
	case reflect.TypeOf(Enum{}):
		field.Set(reflect.ValueOf(Enum{value: str}))
		return nil
	}
	return fmt.Errorf("unsupported string conversion for %s", field.Type())
}

// parseAndSetLogLevel parses and sets a LogLevel field.
func parseAndSetLogLevel(field reflect.Value, str string) error {
	parsed, err := ParseLogLevel(str)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(parsed))
	return nil
}

// parseAndSetLogFormat parses and sets a LogFormat field.
func parseAndSetLogFormat(field reflect.Value, str string) error {
	parsed, err := ParseLogFormat(str)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(parsed))
	return nil
}

// parseAndSetDuration parses and sets a Duration field.
func parseAndSetDuration(field reflect.Value, str string) error {
	parsed, err := ParseDuration(str)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(parsed))
	return nil
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

	// Validate enum values
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
// Returns the value and true if it could be extracted.
func getFieldValue(field reflect.Value) (string, bool) {
	switch field.Kind() {
	case reflect.String:
		return field.String(), true
	default:
		// Handle custom enum types
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
