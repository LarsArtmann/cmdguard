package v2

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"
)

// FlagTypeConstraint validates that F is a valid flag type at initialization time.
// Valid types are: struct{} (NoFlags), any struct, or pointer to struct.
// This enforces type safety for the F type parameter in Command.
// Returns an error if F is an invalid type (e.g., int, string, slice, map).
func FlagTypeConstraint[F any]() error {
	var zero F

	t := reflect.TypeOf(zero)

	// Nil type means F is an untyped nil interface - not valid
	if t == nil {
		return fmt.Errorf(
			"%w: flag type F must be a struct or pointer to struct, got untyped nil, type=%T",
			ErrInvalidFlagType,
			zero,
		)
	}

	switch t.Kind() {
	case reflect.Struct:
		// struct{} (NoFlags) or any struct is valid
		return nil
	case reflect.Pointer:
		// Must be pointer to struct
		if t.Elem().Kind() == reflect.Struct {
			return nil
		}

		fallthrough
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice,
		reflect.String, reflect.UnsafePointer:
		return fmt.Errorf(
			"%w: flag type F must be struct or *struct, got %s, type=%T",
			ErrInvalidFlagType,
			t,
			zero,
		)
	default:
		return fmt.Errorf(
			"%w: flag type F must be struct or *struct, got %s, type=%T",
			ErrInvalidFlagType,
			t,
			zero,
		)
	}
}

// createFlagPrototype creates a flag prototype from the flags value.
func createFlagPrototype[F any](flags F) F {
	if !isNilPointer(flags) {
		return flags
	}

	var zero F

	t := reflect.TypeOf(zero)
	if t != nil && t.Kind() == reflect.Pointer {
		if proto, ok := reflect.New(t.Elem()).Interface().(F); ok {
			return proto
		}
	}

	return zero
}

// isNilPointer checks if a value is a nil pointer or nil interface.
// This is needed because `any(nil) != nil` is true for typed nil pointers.
func isNilPointer(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Array, reflect.String, reflect.Struct, reflect.UnsafePointer:
		return false
	default:
		return false
	}
}

// cloneFlags creates a copy of a flags struct using reflection.
// This ensures each command execution gets its own flag instance.
// Returns the zero value of F if cloning fails or input is nil.
func cloneFlags[F any](flags F) F {
	if isNilPointer(flags) {
		var zero F

		return zero
	}

	// Use reflection to create a new instance
	v := reflect.ValueOf(flags)

	switch v.Kind() { //nolint:exhaustive // default handles remaining kinds
	case reflect.Pointer:
		if v.IsNil() {
			var zero F

			return zero
		}

		newPtr := reflect.New(v.Elem().Type())
		deepCopyValue(newPtr.Elem(), v.Elem())

		if cloned, ok := newPtr.Interface().(F); ok {
			return cloned
		}

		var zero F

		return zero
	case reflect.Struct:
		newStruct := reflect.New(v.Type()).Elem()
		deepCopyValue(newStruct, v)

		if cloned, ok := newStruct.Interface().(F); ok {
			return cloned
		}

		var zero F

		return zero
	default:
	}

	// For other types, return as-is (can't clone safely)
	return flags
}

// createNilFlags creates a new flag instance when flags is a nil pointer.
// Returns (flagsCopy, flagsPtr, error).
func createNilFlags[F any]() (F, any, error) {
	var zero F

	t := reflect.TypeOf(zero)
	if t == nil {
		return zero, nil, nil
	}

	if t.Kind() == reflect.Pointer {
		newVal := reflect.New(t.Elem())

		fc, ok := newVal.Interface().(F)
		if !ok {
			return zero, nil, fmt.Errorf(
				"cloneAndParseFlags: failed to create flag instance for type %T: %w",
				zero,
				ErrFlagInstance,
			)
		}

		return fc, fc, nil
	}

	newPtr := reflect.New(t)

	fc, ok := newPtr.Elem().Interface().(F)
	if !ok {
		return zero, nil, fmt.Errorf(
			"cloneAndParseFlags: failed to create flag instance for type %T: %w",
			zero,
			ErrFlagInstance,
		)
	}

	return fc, newPtr.Interface(), nil
}

// flagsToPtr converts a flags value to a pointer for parsing.
// If flags is already a pointer, returns it directly.
// If flags is a struct, creates a pointer to a copy.
func flagsToPtr[F any](flags F) (F, any) {
	t := reflect.TypeOf(flags)
	if t.Kind() == reflect.Pointer {
		return flags, flags
	}

	newPtr := reflect.New(t)
	newPtr.Elem().Set(reflect.ValueOf(flags))

	return flags, newPtr.Interface()
}

// parseAndSyncFlags parses flags via the registry and syncs parsed values back.
func parseAndSyncFlags[F any](
	c *cobra.Command, flags F, flagsPtr any, registry *FlagRegistry,
) (F, error) {
	if registry == nil {
		return flags, nil
	}

	err := registry.ParseFlags(c, flagsPtr)
	if err != nil {
		return flags, fmt.Errorf(
			"parse flags: command=%q, registry=%T, flags=%T: %w",
			c.Name(),
			registry,
			flags,
			err,
		)
	}

	t := reflect.TypeOf(flags)
	if t != nil && t.Kind() != reflect.Pointer {
		if fc, ok := reflect.ValueOf(flagsPtr).Elem().Interface().(F); ok {
			return fc, nil
		}
	}

	return flags, nil
}

// cloneAndParseFlags clones flags once and parses them.
// This is the optimized single-entry point for flag handling during execution.
// If flags is nil, creates a new instance of F to parse into.
func cloneAndParseFlags[F any](c *cobra.Command, flags F, registry *FlagRegistry) (F, error) {
	var flagsCopy F

	var flagsPtr any

	if isNilPointer(flags) {
		var err error

		flagsCopy, flagsPtr, err = createNilFlags[F]()
		if err != nil {
			return flagsCopy, err
		}
	} else {
		flagsCopy = cloneFlags(flags)
		if any(flagsCopy) == nil {
			flagsCopy = flags
		}

		flagsCopy, flagsPtr = flagsToPtr(flagsCopy)
	}

	return parseAndSyncFlags(c, flagsCopy, flagsPtr, registry)
}

// formatFieldValue converts a reflect.Value to its string representation.
// Handles primitives, pointers, interfaces, and types implementing fmt.Stringer.
func formatFieldValue(field reflect.Value) string {
	if !field.IsValid() {
		return ""
	}

	if field.Kind() == reflect.Pointer || field.Kind() == reflect.Interface {
		if field.Elem().IsValid() {
			return formatFieldValue(field.Elem())
		}

		return ""
	}

	if s, ok := field.Interface().(fmt.Stringer); ok {
		return s.String()
	}

	switch field.Kind() { //nolint:exhaustive // default handles remaining kinds
	case reflect.String:
		return field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(field.Bool())
	default:
		return fmt.Sprintf("%v", field.Interface())
	}
}
