package v3

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
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
		return nil, fmt.Errorf("%w: config is nil", ErrConfigNil)
	}

	v, err := derefPointerToStruct(cfg)
	if err != nil {
		return nil, err
	}

	return parseStructTags(v.Type())
}

// parseStructTags parses all flag tags from a struct type, recursing into
// nested struct fields that do not carry their own flag tag. This lets users
// organize config into nested groups while keeping flat CLI flag names.
func parseStructTags(t reflect.Type) ([]FlagTag, error) {
	return parseStructTagsAtPath(t, nil)
}

// parseStructTagsAtPath walks a struct type at the given reflect index path,
// collecting flag tags from both top-level and nested struct fields.
func parseStructTagsAtPath(t reflect.Type, parentIndex []int) ([]FlagTag, error) {
	var tags []FlagTag

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		tag, ok, err := parseFieldFlag(field)
		if err != nil {
			return nil, err
		}

		if ok {
			tag.Index = appendIndex(parentIndex, i)
			tags = append(tags, tag)

			continue
		}

		// Recurse into nested struct fields without their own flag tag so that
		// flag-bearing fields inside nested config groups are discovered.
		ft := field.Type
		for ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct && shouldRecurseInto(ft) {
			nested, err := parseStructTagsAtPath(ft, appendIndex(parentIndex, i))
			if err != nil {
				return nil, err
			}

			tags = append(tags, nested...)
		}
	}

	return tags, nil
}

// appendIndex returns a new slice with i appended to parentIndex (non-destructive).
func appendIndex(parentIndex []int, i int) []int {
	result := make([]int, 0, len(parentIndex)+1)
	result = append(result, parentIndex...)
	result = append(result, i)

	return result
}

// shouldRecurseInto reports whether a struct type is a nested config group (to
// recurse into) rather than an opaque value type handled by the type registry.
func shouldRecurseInto(t reflect.Type) bool {
	switch t {
	case reflect.TypeFor[Duration](),
		reflect.TypeFor[Port](),
		reflect.TypeFor[Email](),
		reflect.TypeFor[URL](),
		reflect.TypeFor[HostPort](),
		reflect.TypeFor[Enum](),
		reflect.TypeFor[FilePath](),
		reflect.TypeFor[LogLevel](),
		reflect.TypeFor[time.Duration](),
		reflect.TypeFor[time.Time]():
		return false
	}

	return true
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
		if len(short) != 1 {
			return FlagTag{}, false, fmt.Errorf(
				"field %q: short tag must be a single character, got %q: %w",
				field.Name, short, ErrInvalidFlagType,
			)
		}

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

	// Parse required, count, and local boolean tags in one pass.
	bools, err := parseBoolTags(field)
	if err != nil {
		return FlagTag{}, false, err
	}

	tag.Required = bools.Required
	tag.Count = bools.Count
	tag.Local = bools.Local
	tag.Hidden = bools.Hidden

	// Parse validate tag
	if validate := field.Tag.Get("validate"); validate != "" {
		tag.Validate = validate
	}

	// Parse env tag
	if env := field.Tag.Get("env"); env != "" {
		tag.Env = env
	}

	// Parse prompt tag (interactive prompt when flag is missing)
	if prompt := field.Tag.Get("prompt"); prompt != "" {
		tag.Prompt = prompt
	}

	return tag, true, nil
}

// boolTagSet holds the parsed boolean modifier struct tags.
type boolTagSet struct {
	Required bool
	Count    bool
	Local    bool
	Hidden   bool
}

// parseBoolTags parses the required, count, and local boolean struct tags.
// Absent tags keep their zero value (false). An invalid value returns an error
// naming the offending tag key, matching the historical per-tag error messages.
func parseBoolTags(field reflect.StructField) (boolTagSet, error) {
	var out boolTagSet

	for _, t := range []struct {
		key  string
		dest *bool
	}{
		{"required", &out.Required},
		{"count", &out.Count},
		{"local", &out.Local},
		{"hidden", &out.Hidden},
	} {
		raw := field.Tag.Get(t.key)
		if raw == "" {
			continue
		}

		value, err := strconv.ParseBool(raw)
		if err != nil {
			return boolTagSet{}, fmt.Errorf(
				"field %q: invalid %s tag %q: %w",
				field.Name, t.key, raw, err,
			)
		}

		*t.dest = value
	}

	return out, nil
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

// DefaultValue returns the default value for a flag based on its type.
// Uses the global type registry for backward compatibility.
func (t FlagTag) DefaultValue() any {
	return dispatchDefault(globalTypeRegistry, t)
}
