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
		tag, ok, err := parseFieldFlag(field)
		if err != nil {
			return nil, err
		}

		if ok {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

// parseFieldFlag parses flag tags from a single struct field.
// Returns (FlagTag, true, nil) if a flag tag was found.
// Returns (FlagTag{}, false, nil) if the field should be skipped.
func parseFieldFlag(field reflect.StructField) (FlagTag, bool, error) {
	flagTag := field.Tag.Get("flag")

	// Skip fields without flag tag
	if flagTag == "" || flagTag == "-" {
		return FlagTag{}, false, nil
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
		required, err := strconv.ParseBool(req)
		if err != nil {
			return FlagTag{}, false, fmt.Errorf("field %q: invalid required tag %q: %w", field.Name, req, err)
		}

		tag.Required = required
	}

	// Parse validate tag
	if validate := field.Tag.Get("validate"); validate != "" {
		tag.Validate = validate
	}

	return tag, true, nil
}

// parseBoolDefault parses a boolean default value.
// Returns an error if the string is not empty and not a valid boolean.
func parseBoolDefault(s string) (bool, error) {
	if s == "" {
		return false, nil
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("failed to parse bool %q: %w", s, err)
	}

	return v, nil
}

// parseIntDefault parses an integer default value.
func parseIntDefault(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int %q: %w", s, err)
	}

	return v, nil
}

// parseUintDefault parses an unsigned integer default value.
func parseUintDefault(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}

	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uint %q: %w", s, err)
	}

	return v, nil
}

// parseFloat64Default parses a float64 default value.
func parseFloat64Default(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float %q: %w", s, err)
	}

	return v, nil
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
		v, _ := parseBoolDefault(t.Default)

		return v
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t.Type == reflect.TypeFor[Duration]() {
			d, err := ParseDuration(t.Default)
			if err != nil {
				return Duration{}
			}

			return d
		}

		v, _ := parseIntDefault(t.Default)

		return int(v)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		v, _ := parseUintDefault(t.Default)

		return uint(v)
	case reflect.Float32, reflect.Float64:
		v, _ := parseFloat64Default(t.Default)

		return v
	case reflect.Slice:
		return strings.Split(t.Default, ",")
	case reflect.Invalid, reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Struct, reflect.UnsafePointer:
		return t.parseCustomDefault()
	default:
		return t.parseCustomDefault()
	}
}
