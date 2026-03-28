package v2

import "fmt"

// Ptr returns a pointer to any value.
// Useful for optional config fields.
//
//go:fix inline
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
		return fmt.Errorf("%s (%T): must not be nil", name, v)
	}

	return nil
}
