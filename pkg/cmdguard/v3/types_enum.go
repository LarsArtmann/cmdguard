package v3

import (
	"slices"
)

// Enum provides type-safe enum values with validation.
// Use this for config fields that must be one of a set of allowed values.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type Enum struct {
	value   string
	allowed []string
}

// ParseEnum creates a new Enum from a string value.
// Returns an error if the value is not in the allowed list.
func ParseEnum(value string, allowed []string) (Enum, error) {
	if slices.Contains(allowed, value) {
		return Enum{value: value, allowed: allowed}, nil
	}

	return Enum{}, NewEnumError(value, allowed)
}

// String returns the enum value as a string.
func (e Enum) String() string {
	return e.value
}

// Value returns the enum value.
func (e Enum) Value() string {
	return e.value
}

// Allowed returns a copy of the list of allowed values.
func (e Enum) Allowed() []string {
	return slices.Clone(e.allowed)
}

// IsEmpty returns true if the enum has no value.
func (e Enum) IsEmpty() bool {
	return e.value == ""
}

// MarshalText implements encoding.TextMarshaler for Enum.
// Hand-written instead of using textMarshal because Enum.String() is simple enough.
func (e Enum) MarshalText() ([]byte, error) {
	return []byte(e.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for Enum.
// Hand-written because textUnmarshal can't handle the bootstrap behavior
// (accepting any value when allowed list is empty).
func (e *Enum) UnmarshalText(text []byte) error {
	value := string(text)
	if slices.Contains(e.allowed, value) {
		e.value = value

		return nil
	}

	if len(e.allowed) == 0 {
		// If no allowed values set yet, accept any value and initialize allowed list
		e.value = value
		e.allowed = []string{value}

		return nil
	}

	return NewEnumError(value, e.allowed)
}
