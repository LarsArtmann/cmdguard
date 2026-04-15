package v2

import (
	"encoding/json"
	"fmt"
)

// Result represents either a successful value (Ok) or an error (Err).
// This is similar to Rust's Result type and provides type-safe error handling
// without requiring separate return values.
//
// Use Ok(v) to create a successful Result, Err[T](e) for an error Result.
//
// Example:
//
//	res := Ok(42)
//	if v, err := res.Get(); err == nil {
//	    fmt.Println(v) // 42
//	}
//
//	// Chain operations:
//	doubled := Ok(21).Map(func(v int) int { return v * 2 })
//	fmt.Println(doubled.UnwrapOr(0)) // 42
//
//	// Error handling:
//	res := Err[int](fmt.Errorf("something failed"))
//	fmt.Println(res.UnwrapOr(-1)) // -1
type Result[T any] struct {
	value T
	err   error
	ok    bool
}

// Ok creates a successful Result containing a value.
func Ok[T any](v T) Result[T] {
	return Result[T]{value: v, ok: true}
}

// Err creates an error Result.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, ok: false}
}

// IsOk returns true if the Result contains a value.
func (r Result[T]) IsOk() bool {
	return r.ok
}

// IsErr returns true if the Result contains an error.
func (r Result[T]) IsErr() bool {
	return !r.ok
}

// Get returns the value and error.
// If Ok, returns (value, nil). If Err, returns (zero, error).
func (r Result[T]) Get() (T, error) {
	if r.ok {
		return r.value, nil
	}

	return r.value, r.err
}

// Unwrap returns the value, panicking if Err.
// Use only when you're certain the Result is Ok.
func (r Result[T]) Unwrap() T {
	if !r.ok {
		panic(fmt.Sprintf("called Result.Unwrap() on Err: %v", r.err))
	}

	return r.value
}

// UnwrapOr returns the value if Ok, otherwise returns the provided default.
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.ok {
		return r.value
	}

	return defaultValue
}

// UnwrapOrElse returns the value if Ok, otherwise computes from the provided function.
func (r Result[T]) UnwrapOrElse(f func(error) T) T {
	if r.ok {
		return r.value
	}

	return f(r.err)
}

// UnwrapErr returns the error, panicking if Ok.
func (r Result[T]) UnwrapErr() error {
	if r.ok {
		panic(fmt.Sprintf("called Result.UnwrapErr() on Ok: %v", r.value))
	}

	return r.err
}

// Expect returns the value, panicking with the given message if Err.
func (r Result[T]) Expect(msg string) T {
	if !r.ok {
		panic(fmt.Sprintf("%s: %v", msg, r.err))
	}

	return r.value
}

// ExpectErr returns the error, panicking with the given message if Ok.
func (r Result[T]) ExpectErr(msg string) error {
	if r.ok {
		panic(fmt.Sprintf("%s: %v", msg, r.value))
	}

	return r.err
}

// Map applies a function to the contained value if Ok, returning a new Result.
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.ok {
		return Ok(f(r.value))
	}

	return Err[T](r.err)
}

// MapErr applies a function to the contained error if Err, returning a new Result.
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.ok {
		return Ok(r.value)
	}

	return Err[T](f(r.err))
}

// MapOr applies a function to the contained value if Ok, otherwise returns the default.
func (r Result[T]) MapOr(defaultValue T, f func(T) T) T {
	if r.ok {
		return f(r.value)
	}

	return defaultValue
}

// And returns Err if this Result is Err, otherwise returns the provided Result.
func (r Result[T]) And(other Result[T]) Result[T] {
	if r.ok {
		return other
	}

	return Err[T](r.err)
}

// Or returns this Result if Ok, otherwise returns the provided Result.
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.ok {
		return r
	}

	return other
}

// IfOk executes the given function if the Result is Ok.
func (r Result[T]) IfOk(f func(T)) {
	if r.ok {
		f(r.value)
	}
}

// IfErr executes the given function if the Result is Err.
func (r Result[T]) IfErr(f func(error)) {
	if !r.ok {
		f(r.err)
	}
}

// ToOption converts a Result to an Option.
// Ok(v) becomes Some(v), Err becomes None.
func (r Result[T]) ToOption() Option[T] {
	if r.ok {
		return Some(r.value)
	}

	return None[T]()
}

// String returns a string representation of the Result.
func (r Result[T]) String() string {
	if r.ok {
		return fmt.Sprintf("Ok(%v)", r.value)
	}

	return fmt.Sprintf("Err(%v)", r.err)
}

// MarshalJSON implements json.Marshaler for Result.
// Serializes as the value if Ok, or as {"error": "msg"} if Err.
func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.ok {
		data, err := json.Marshal(any(r.value))
		if err != nil {
			return nil, fmt.Errorf("marshaling Result value: %w", err)
		}

		return data, nil
	}

	errData, err := json.Marshal(map[string]string{"error": r.err.Error()})
	if err != nil {
		return nil, fmt.Errorf("marshaling Result error: %w", err)
	}

	return errData, nil
}

// ResultFrom converts a (value, error) pair into a Result.
// This is the bridge between Go's conventional error handling and Result.
func ResultFrom[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}

	return Ok(value)
}

// ToPair converts a Result back to a (value, error) pair.
// This is the bridge between Result and Go's conventional error handling.
func (r Result[T]) ToPair() (T, error) {
	if r.ok {
		return r.value, nil
	}

	var zero T

	return zero, r.err
}
