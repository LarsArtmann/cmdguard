package v2

import (
	"errors"
	"fmt"
)

// Result is a value-or-error sum type (like Rust's Result<T, E>).
// It makes error handling explicit at the type level: a Result is either Ok
// (carrying a value) or Err (carrying an error). All accessors are panic-free.
type Result[T any] struct {
	value T
	err   error
}

// Ok constructs a successful Result carrying value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err constructs a failed Result carrying err.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// IsOk reports whether the Result carries a value (no error).
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr reports whether the Result carries an error.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Value returns the carried value and a bool indicating success.
// When false, the value is the zero value of T.
func (r Result[T]) Value() (T, bool) {
	return r.value, r.err == nil
}

// MustValue returns the carried value. Panics if the Result is an error.
// Use only when the caller has already checked IsOk or is certain of success.
func (r Result[T]) MustValue() T {
	if r.err != nil {
		panic(fmt.Sprintf("Result.MustValue on Err: %v", r.err))
	}

	return r.value
}

// ErrValue returns the carried error, or nil if Ok.
func (r Result[T]) ErrValue() error {
	return r.err
}

// UnwrapOr returns the carried value, or fallback if Err.
func (r Result[T]) UnwrapOr(fallback T) T {
	if r.err == nil {
		return r.value
	}

	return fallback
}

// UnwrapOrElse returns the carried value, or the result of fn if Err.
func (r Result[T]) UnwrapOrElse(fn func(error) T) T {
	if r.err == nil {
		return r.value
	}

	return fn(r.err)
}

// Map applies fn to the carried value if Ok, returning a new Result.
// If Err, the error is propagated unchanged.
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.err != nil {
		return r
	}

	return Ok(fn(r.value))
}

// AndThen (flat-map) applies fn to the carried value if Ok.
// If Err, the error is propagated unchanged. Use for chaining fallible operations.
func (r Result[T]) AndThen(fn func(T) Result[T]) Result[T] {
	if r.err != nil {
		return r
	}

	return fn(r.value)
}

// Validated accumulates zero or more validation errors alongside a value.
// Unlike Result (which is Ok-or-Err), Validated collects ALL errors, making it
// ideal for validating multiple independent fields before reporting.
type Validated[T any] struct {
	value T
	errs  []error
}

// Valid constructs a Validated with no errors.
func Valid[T any](value T) Validated[T] {
	return Validated[T]{value: value}
}

// Invalid constructs a Validated with the given errors.
func Invalid[T any](value T, errs ...error) Validated[T] {
	return Validated[T]{value: value, errs: errs}
}

// IsValid reports whether there are no accumulated errors.
func (v Validated[T]) IsValid() bool {
	return len(v.errs) == 0
}

// Errors returns the accumulated errors (nil if valid).
func (v Validated[T]) Errors() []error {
	return v.errs
}

// Value returns the carried value and a bool indicating validity.
func (v Validated[T]) Value() (T, bool) {
	return v.value, len(v.errs) == 0
}

// AddErr appends an error to the accumulator (mutates receiver).
func (v *Validated[T]) AddErr(err error) {
	if err != nil {
		v.errs = append(v.errs, err)
	}
}

// Combine merges another Validated's errors into this one.
func (v *Validated[T]) Combine(other Validated[T]) {
	v.errs = append(v.errs, other.errs...)
}

// ToResult converts to a Result: Ok if valid, Err if any errors (joined).
func (v Validated[T]) ToResult() Result[T] {
	if len(v.errs) == 0 {
		return Ok(v.value)
	}

	return Err[T](errors.Join(v.errs...))
}
