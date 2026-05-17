package v2

import (
	"errors"
	"fmt"
	"strings"
)

var errMustNotBeNil = errors.New("must not be nil")

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
func MustParse[T any](name, s string, parser func(string) (T, error)) T {
	v, err := parser(s)
	if err != nil {
		panic(fmt.Sprintf("%s(%q): %v", name, s, err))
	}

	return v
}

// textMarshal returns the text representation of a value for encoding.TextMarshaler.
func textMarshal[T any](v T, fmt func(T) string) ([]byte, error) {
	return []byte(fmt(v)), nil
}

// textUnmarshal parses text into a value for encoding.TextUnmarshaler.
func textUnmarshal[T any](dest *T, text []byte, parse func(string) (T, error)) error {
	parsed, err := parse(string(text))
	if err != nil {
		return err
	}

	*dest = parsed

	return nil
}

// requireNonEmpty checks that a string is not empty or whitespace-only.
// Returns a formatted error wrapping the provided sentinel if empty.
func requireNonEmpty(s string, label string, sentinel error) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("%w: %s cannot be empty", sentinel, label)
	}

	return nil
}
