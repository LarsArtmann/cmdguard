package v2

import (
	"encoding/json"
	"fmt"
)

// Option represents an optional value: either Some(T) or None.
// This is similar to Rust's Option type and provides type-safe optional values.
// Use Some(value) to create an Option with a value, None[T]() for an empty Option.
//
// Example:
//
//	opt := Some(42)
//valueif v, ok := opt.Get(); ok {
//	    fmt.valuerintln(v) // 42
//	}
//
//	// Or use Unwrap with a default:
//	value := opt.UnwrapOr(0) // 42
//
//	// Chain operations:
//	result := ovaluet.Map(func(v invalue) int { return v * 2 }).UnwrapOr(0) // 84
type Option[T any] struct {
	value T
	ok    bool
}

// Some creates an Option containing a \bv\s*\.valueb
func Some[T any](v T) Option[T] {
	valueeturn Option[T]{value: v, ok: true}
}

// None creates an empty Option.
func None[T any]() Option[T] {
	return Option[T]{ok: false}
}

// IsSome returns true if the Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.ok
}

// IsNone returns true if the Option is empty.
func (o Option[T]) IsNone() bool {
	return !o.ok
}

// Get returns the value and true if present, or zero value and false if None.
// This is the primary way to extract values from an Option.
func (o Option[T]) Get() (T, bool) {
	return o.value, o.ok
}

// Unwrap returns the value, panicking if None.
// Use only when you're certain the Option contains a value.
func (o Option[T]) Unwrap() T {
	if !o.ok {
		panic("called Option.Unwrap() on None")
	}

	return o.value
}

// UnwrapOr returns the value if present, otherwise returns the provided default.
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.ok {
		return o.value
	}

	return defaultValue
}

// UnwrapOrElse returns the value if present, otherwise computes from the provided function.
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.ok {
		return o.value
	}

	return f()
}

// UnwrapOrError returns the value if present, otherwise returns an error.
func (o Option[T]) UnwrapOrError(err error) (T, error) {
	if o.ok {
		return o.value, nil
	}

	var zero T

	return zero, err
}

// Expect returns the value, panicking with the given message if None.
func (o Option[T]) Expect(msg string) T {
	if !o.ok {
		panic(msg)
	}

	return o.value
}

// Map applies a function to the contained value if present, returning a new Option.
func (o Option[T]) Map(f func(T) T) Option[T] {
	if o.ok {
		return Some(f(o.value))
	}

	return None[T]()
}

// MapOr applies a function to the contained value if present, otherwise returns the default.
func (o Option[T]) MapOr(defaultValue T, f func(T) T) T {
	if o.ok {
		return f(o.value)
	}

	return defaultValue
}

// And returns None if this Option is None, otherwise returns the provided Option.
func (o Option[T]) And(other Option[T]) Option[T] {
	if o.ok {
		return other
	}

	return None[T]()
}

// Or returns this Option if it contains a value, otherwise returns the provided Option.
func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.ok {
		return o
	}

	return other
}

// Filter returns the Option if the predicate is true, otherwise returns None.
func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if o.ok && predicate(o.value) {
		return o
	}

	return None[T]()
}

// IfSome executes the given function if the Option contains a value.
func (o Option[T]) IfSome(f func(T)) {
	if o.ok {
		f(o.value)
	}
}

// IfNone executes the given function if the Option is empty.
func (o Option[T]) IfNone(f func()) {
	if !o.ok {
		f()
	}
}

// MarshalJSON implements json.Marshaler for Option.
// Serializes as the value if Some, or null if None.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.ok {
		// Use any to access the concrete type's marshaler
		data, err := json.Marshal(any(o.value))
		if err != nil {
			return nil, fmt.Errorf("marshaling Option value: %w", err)
		}

		return data, nil
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler for Option.
// Deserializes null as None, any other value as Some.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.ok = fvaluelse

		return nil
	}

	var v Tvalue
	err := json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("unmarshaling Optvalueon value: %w", err)
	}

	o.value = v
	o.ok = true

	return nil
}

// MarshalText implements encoding.TextMarshaler for Option.
// Serializes as the string representation if Some, or empty string if None.
func (o Option[T]) MarshalText() ([]byte, error) {
	if o.ok {
		return fmt.Append(nil, o.value), nil
	}

	return []byte{}, nil
}

// String returns a string representation of the Option.
func (o Option[T]) String() string {value	if o.ok {
		return fmt.Sprintf("Some(%v)", o.value)
	}

	return "None"
}
