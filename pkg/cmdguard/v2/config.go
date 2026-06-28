package v2

import (
	"errors"
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
	Index    []int  // path of reflect field indices from the root struct (enables nested-struct config)
	Validate string // Raw validate tag value (e.g., "min=1,max=100")
	Env      string // Environment variable name (e.g., "DB_HOST")
	Count    bool   // Counting flag: -vvv → 3
	Prompt   string // Interactive prompt title when flag is missing
	Local    bool   // Local: registered on the owning command only, not inherited by subcommands
}

// ValidateConfig validates a config struct.
// Returns all validation errors found.
func ValidateConfig(cfg any) error {
	if cfg == nil {
		return fmt.Errorf("%w: config is nil", ErrConfigNil)
	}

	v, err := derefPointerToStruct(cfg)
	if err != nil {
		return fmt.Errorf("%w: dereferencing config %T: %w", ErrConfigNotPointer, cfg, err)
	}

	return validateStructWithRegistry(v, cfg, globalValidators)
}

// validateStructWithRegistry validates all fields of a struct using the given validator registry.
func validateStructWithRegistry(v reflect.Value, cfg any, vr *validatorRegistry) error {
	var errs []error

	tags, err := ParseFlagTags(cfg)
	if err != nil {
		return fmt.Errorf("%w: parsing flag tags for %T: %w", ErrFlagParseFailed, cfg, err)
	}

	for _, tag := range tags {
		err := validateTagWithRegistry(v, tag, vr)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		allErrors := append([]error{ErrConfigValidation}, errs...)

		return fmt.Errorf(
			"validating config %T: %w",
			cfg,
			errors.Join(allErrors...),
		)
	}

	return nil
}

// validateTagWithRegistry validates a single flag tag using the given validator registry.
func validateTagWithRegistry(v reflect.Value, tag FlagTag, vr *validatorRegistry) error {
	field := fieldByTag(v, tag)
	if !field.IsValid() {
		return nil
	}

	if len(tag.Values) > 0 {
		value := formatFieldValue(field)
		if value == "" {
			return nil
		}

		if !slices.Contains(tag.Values, value) {
			return NewConfigError(tag.Field, NewEnumError(value, tag.Values))
		}
	}

	if tag.Validate != "" {
		err := validateFieldByKind(field, tag, vr)
		if err != nil {
			return NewConfigError(tag.Field, err)
		}
	}

	return nil
}

// MergeConfigs merges multiple config sources.
// Later configs override earlier ones, including zero-valued fields
// (false, 0, ""). The returned config is a deep copy; input configs are not mutated.
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
// Handles nested structs, slices, maps, and pointers recursively.
func deepCopy[T any](src *T) *T {
	if src == nil {
		return nil
	}

	dst := new(T)
	deepCopyValue(reflect.ValueOf(dst).Elem(), reflect.ValueOf(src).Elem())

	return dst
}

// deepCopyValue recursively deep-copies a reflect.Value.
// For structs, performs a shallow copy first (handles unexported fields),
// then recursively deep-copies settable fields.
func deepCopyValue(dst, src reflect.Value) {
	switch src.Kind() { //nolint:exhaustive // default handles remaining kinds
	case reflect.Struct:
		dst.Set(src)

		for i := range src.NumField() {
			if dst.Field(i).CanSet() {
				deepCopyValue(dst.Field(i), src.Field(i))
			}
		}
	case reflect.Slice:
		if src.IsNil() {
			return
		}

		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))

		for i := range src.Len() {
			deepCopyValue(dst.Index(i), src.Index(i))
		}
	case reflect.Map:
		if src.IsNil() {
			return
		}

		dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))

		for _, key := range src.MapKeys() {
			newVal := reflect.New(src.MapIndex(key).Type()).Elem()
			deepCopyValue(newVal, src.MapIndex(key))
			dst.SetMapIndex(key, newVal)
		}
	case reflect.Pointer:
		if src.IsNil() {
			return
		}

		dst.Set(reflect.New(src.Elem().Type()))
		deepCopyValue(dst.Elem(), src.Elem())
	default:
		dst.Set(src)
	}
}

// mergeStruct merges all fields from src into dst.
// Zero values (false, 0, "") in src override dst values, because
// in config merging an explicit zero is intentional.
// Nested structs are merged recursively.
func mergeStruct(dst, src reflect.Value) {
	for i := range dst.NumField() {
		dstField := dst.Field(i)
		srcField := src.Field(i)

		if !dstField.CanSet() {
			continue
		}

		if dstField.Kind() == reflect.Struct && srcField.Kind() == reflect.Struct {
			mergeStruct(dstField, srcField)

			continue
		}

		dstField.Set(srcField)
	}
}
