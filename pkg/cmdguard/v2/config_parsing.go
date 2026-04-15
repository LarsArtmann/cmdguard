package v2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// derefPointerToStruct dereferences a pointer to a struct and returns its value.
// Returns an error if the input is not a pointer to a struct.
func derefPointerToStruct(cfg any) (reflect.Value, error) {
	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("%w: expected struct, got %T", ErrInvalidFlagType, cfg)
	}

	return v, nil
}

// ParseFlagTags extracts flag information from a config struct.
// The struct must have `flag` tags on its fields.
func ParseFlagTags(cfg any) ([]FlagTag, error) {
	if cfg == nil {
		return nil, ErrConfigNil
	}

	v, err := derefPointerToStruct(cfg)
	if err != nil {
		return nil, err
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

// parseUintDefault parses an unsigned integer default value.
func (t FlagTag) parseUintDefault() any {
	v, _ := strconv.ParseUint(t.Default, 10, 64)

	return uint(v)
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return t.parseIntDefault()
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return t.parseUintDefault()
	case reflect.Float32, reflect.Float64:
		return parseFloat64Default(t.Default)
	case reflect.Slice:
		return strings.Split(t.Default, ",")
	case reflect.Invalid, reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Struct, reflect.UnsafePointer:
		return t.parseCustomDefault()
	default:
		return t.parseCustomDefault()
	}
}
