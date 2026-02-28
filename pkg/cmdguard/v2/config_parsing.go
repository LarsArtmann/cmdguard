package v2

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ParseFlagTags extracts flag information from a config struct.
// The struct must have `flag` tags on its fields.
func ParseFlagTags(cfg any) ([]FlagTag, error) {
	if cfg == nil {
		return nil, errors.New("config must not be nil")
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

	for field := range t.Fields() {
		tag := parseFieldFlag(field)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}

	return tags, nil
}

// parseFieldFlag parses flag tags from a single struct field.
// Returns nil if the field doesn't have a flag tag.
func parseFieldFlag(field reflect.StructField) *FlagTag {
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

// parseBoolDefault parses a boolean default value.
func parseBoolDefault(s string) bool {
	v, _ := strconv.ParseBool(s)

	return v
}

// parseIntDefault parses an integer default value.
func (t FlagTag) parseIntDefault() any {
	// Check if it's a Duration type
	if t.Type == reflect.TypeFor[Duration]() {
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
	case reflect.TypeFor[Duration]():
		d, err := ParseDuration(t.Default)
		if err != nil {
			return Duration{}
		}

		return d
	case reflect.TypeFor[Enum](), reflect.TypeFor[LogLevel](), reflect.TypeFor[LogFormat]():
		return t.Default
	default:
		return t.Default
	}
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
