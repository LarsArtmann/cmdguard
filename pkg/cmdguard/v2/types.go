package v2

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"time"
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

// Allowed returns the list of allowed values.
func (e Enum) Allowed() []string {
	return e.allowed
}

// IsEmpty returns true if the enum has no value.
func (e Enum) IsEmpty() bool {
	return e.value == ""
}

// Duration wraps time.Duration with parsing validation.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type Duration struct {
	duration time.Duration
}

// ParseDuration creates a new Duration from a string.
// Accepts formats like "30s", "5m", "1h30m", etc.
func ParseDuration(s string) (Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return Duration{}, NewDurationError(s, err)
	}

	return Duration{duration: d}, nil
}

// FromDuration creates a Duration from a time.Duration.
func FromDuration(d time.Duration) Duration {
	return Duration{duration: d}
}

// Duration returns the underlying time.Duration.
func (d Duration) Duration() time.Duration {
	return d.duration
}

// String returns the duration as a string.
func (d Duration) String() string {
	return d.duration.String()
}

// IsZero returns true if the duration is zero.
func (d Duration) IsZero() bool {
	return d.duration == 0
}

// Milliseconds returns the duration in milliseconds.
func (d Duration) Milliseconds() int64 {
	return d.duration.Milliseconds()
}

// Seconds returns the duration in seconds as a float64.
func (d Duration) Seconds() float64 {
	return d.duration.Seconds()
}

// LogLevel is a type-safe log level enum.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type LogLevel Enum

// Log level constants.
//
//nolint:gochecknoglobals // Predefined constants for type-safe defaults
var (
	LogLevelDebug = LogLevel{value: "debug", allowed: []string{"debug", "info", "warn", "error"}}
	LogLevelInfo  = LogLevel{value: "info", allowed: []string{"debug", "info", "warn", "error"}}
	LogLevelWarn  = LogLevel{value: "warn", allowed: []string{"debug", "info", "warn", "error"}}
	LogLevelError = LogLevel{value: "error", allowed: []string{"debug", "info", "warn", "error"}}
)

// ParseLogLevel creates a LogLevel from a string.
func ParseLogLevel(s string) (LogLevel, error) {
	e, err := ParseEnum(s, []string{"debug", "info", "warn", "error"})
	if err != nil {
		return LogLevel{}, err
	}

	return LogLevel(e), nil
}

// String returns the log level as a string.
func (l LogLevel) String() string {
	return l.value
}

// SlogLevel converts LogLevel to slog.Level for use with log/slog.
func (l LogLevel) SlogLevel() slog.Level {
	switch l.value {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// LogFormat is a type-safe log format enum.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type LogFormat Enum

// Log format constants.
//
//nolint:gochecknoglobals // Predefined constants for type-safe defaults
var (
	LogFormatText = LogFormat{value: "text", allowed: []string{"text", "json"}}
	LogFormatJSON = LogFormat{value: "json", allowed: []string{"text", "json"}}
)

// ParseLogFormat creates a LogFormat from a string.
func ParseLogFormat(s string) (LogFormat, error) {
	e, err := ParseEnum(s, []string{"text", "json"})
	if err != nil {
		return LogFormat{}, err
	}

	return LogFormat(e), nil
}

// String returns the log format as a string.
func (f LogFormat) String() string {
	return f.value
}

// TextMarshaler implementations for JSON/YAML encoding

// MarshalText implements encoding.TextMarshaler for Enum.
func (e Enum) MarshalText() ([]byte, error) {
	return []byte(e.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for Enum.
// If no allowed values are defined, accepts any value and sets allowed to contain it.
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

// MarshalText implements encoding.TextMarshaler for Duration.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.duration.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for Duration.
func (d *Duration) UnmarshalText(text []byte) error {
	parsed, err := ParseDuration(string(text))
	if err != nil {
		return err
	}

	*d = parsed

	return nil
}

// MarshalText implements encoding.TextMarshaler for LogLevel.
func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for LogLevel.
func (l *LogLevel) UnmarshalText(text []byte) error {
	parsed, err := ParseLogLevel(string(text))
	if err != nil {
		return err
	}

	*l = parsed

	return nil
}

// MarshalText implements encoding.TextMarshaler for LogFormat.
func (f LogFormat) MarshalText() ([]byte, error) {
	return []byte(f.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for LogFormat.
func (f *LogFormat) UnmarshalText(text []byte) error {
	parsed, err := ParseLogFormat(string(text))
	if err != nil {
		return err
	}

	*f = parsed

	return nil
}

// Option represents an optional value: either Some(T) or None.
// This is similar to Rust's Option type and provides type-safe optional values.
// Use Some(v) to create an Option with a value, None[T]() for an empty Option.
//
// Example:
//
//	opt := Some(42)
//	if v, ok := opt.Get(); ok {
//	    fmt.Println(v) // 42
//	}
//
//	// Or use Unwrap with a default:
//	value := opt.UnwrapOr(0) // 42
//
//	// Chain operations:
//	result := opt.Map(func(v int) int { return v * 2 }).UnwrapOr(0) // 84
type Option[T any] struct {
	value T
	ok    bool
}

// Some creates an Option containing a value.
func Some[T any](v T) Option[T] {
	return Option[T]{value: v, ok: true}
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
		return json.Marshal(any(o.value))
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler for Option.
// Deserializes null as None, any other value as Some.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.ok = false

		return nil
	}

	var v T

	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
func (o Option[T]) String() string {
	if o.ok {
		return fmt.Sprintf("Some(%v)", o.value)
	}

	return "None"
}
