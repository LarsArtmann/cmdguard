package v2

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
)

// FlagTypeConstraint validates that F is a valid flag type at initialization time.
// Valid types are: struct{} (NoFlags), any struct, or pointer to struct.
// This enforces type safety for the F type parameter in GuardedCommand and Command.
// Returns an error if F is an invalid type (e.g., int, string, slice, map).
func FlagTypeConstraint[F any]() error {
	var zero F

	t := reflect.TypeOf(zero)

	// Nil type means F is an untyped nil interface - not valid
	if t == nil {
		return fmt.Errorf(
			"%w: flag type F must be a struct or pointer to struct, got untyped nil",
			ErrInvalidFlagType,
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

		return fmt.Errorf(
			"%w: flag type F must be struct or *struct, got %s",
			ErrInvalidFlagType,
			t,
		)
	default:
		return fmt.Errorf(
			"%w: flag type F must be struct or *struct, got %s",
			ErrInvalidFlagType,
			t,
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
		return reflect.New(t.Elem()).Interface().(F)
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
	}

	return false
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

	// Handle pointer to struct
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			var zero F

			return zero
		}
		// Create new pointer to same type
		newPtr := reflect.New(v.Elem().Type())
		// Copy the value
		newPtr.Elem().Set(v.Elem())

		return newPtr.Interface().(F)
	}

	// Handle struct directly
	if v.Kind() == reflect.Struct {
		newStruct := reflect.New(v.Type()).Elem()
		newStruct.Set(v)

		return newStruct.Interface().(F)
	}

	// For other types, return as-is (can't clone safely)
	return flags
}

// cloneAndParseFlags clones flags once and parses them.
// This is the optimized single-entry point for flag handling during execution.
// If flags is nil, creates a new instance of F to parse into.
func cloneAndParseFlags[F any](c *cobra.Command, flags F, registry *FlagRegistry) (F, error) {
	var flagsCopy F

	var flagsPtr any // Pointer to flags for parsing (SetField requires pointer)

	// If flags is nil, create a new instance of the flag type
	if isNilPointer(flags) {
		// Create new instance using reflection
		var zero F

		t := reflect.TypeOf(zero)
		if t == nil {
			// F is an interface type with nil value - can't create
			return zero, nil
		}

		if t.Kind() == reflect.Pointer {
			// Create new instance of the underlying type
			newVal := reflect.New(t.Elem())
			flagsCopy = newVal.Interface().(F)
			flagsPtr = flagsCopy
		} else {
			// F is a struct type (like NoFlags) - create pointer for parsing
			newPtr := reflect.New(t)
			flagsPtr = newPtr.Interface()
			flagsCopy = newPtr.Elem().Interface().(F)
		}
	} else {
		// Clone the flags struct
		flagsCopy = cloneFlags(flags)
		if any(flagsCopy) == nil {
			flagsCopy = flags
		}

		// Create pointer for parsing
		t := reflect.TypeOf(flagsCopy)
		if t.Kind() == reflect.Pointer {
			flagsPtr = flagsCopy
		} else {
			// F is a struct - create pointer for parsing
			newPtr := reflect.New(t)
			newPtr.Elem().Set(reflect.ValueOf(flagsCopy))
			flagsPtr = newPtr.Interface()
		}
	}

	// Parse command-line values into the flags
	if registry != nil {
		err := registry.ParseFlags(c, flagsPtr)
		if err != nil {
			return flagsCopy, fmt.Errorf("parse flags: %w", err)
		}
		// Copy parsed values back to flagsCopy if it was a struct
		t := reflect.TypeOf(flagsCopy)
		if t != nil && t.Kind() != reflect.Pointer {
			// flagsPtr is *F, dereference to get the parsed values
			flagsCopy = reflect.ValueOf(flagsPtr).Elem().Interface().(F)
		}
	}

	return flagsCopy, nil
}
