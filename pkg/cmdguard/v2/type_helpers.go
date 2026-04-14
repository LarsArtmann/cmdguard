package v2

import (
	"errors"
	"fmt"
)

var errMustNotBeNil = errors.New("must not be nil")

// Ptr returns a pointer to any value.
// Useful for optional config fields.
func Ptr[T any](v T) *T {
	return new(v)
}

// ValueOrDefault returns the value if not nil, otherwise the default.
func ValueOrDefault[T any](v *T, def T) T {
	if v == nil {
		return def
	}

	return *v
}

// EnsureValid validates that a pointer is not nil and returns an error with context.
func EnsureValid[T any](v *T, name string) error {
	if v == nil {
		return fmt.Errorf("%s (%T): %w", name, v, errMustNotBeNil)
	}

	return nil
}

// MustParse parses a value or panics with a descriptive message.
// Use only when you know the value is valid (e.g., for constants).
func MustParse[T any](name string, s string, parser func(string) (T, error)) T {
	v, err := parser(s)
	if err != nil {
		panic(fmt.Sprintf("%s(%q): %v", name, s, err))
	}
	return v
}
