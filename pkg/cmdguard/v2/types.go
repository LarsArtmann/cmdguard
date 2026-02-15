package v2

import (
	"fmt"
	"time"
)

// Enum provides type-safe enum values with validation.
// Use this for config fields that must be one of a set of allowed values.
type Enum struct {
	value   string
	allowed []string
}

// ParseEnum creates a new Enum from a string value.
// Returns an error if the value is not in the allowed list.
func ParseEnum(value string, allowed []string) (Enum, error) {
	for _, a := range allowed {
		if value == a {
			return Enum{value: value, allowed: allowed}, nil
		}
	}
	return Enum{}, NewEnumError(value, allowed)
}

// MustEnum creates a new Enum or panics.
// Use only in tests or when the value is guaranteed to be valid.
func MustEnum(value string, allowed []string) Enum {
	e, err := ParseEnum(value, allowed)
	if err != nil {
		panic(err)
	}
	return e
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

// MustDuration creates a new Duration or panics.
// Use only in tests or when the value is guaranteed to be valid.
func MustDuration(s string) Duration {
	d, err := ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
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
type LogLevel Enum

// Log level constants.
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

// LogFormat is a type-safe log format enum.
type LogFormat Enum

// Log format constants.
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
func (e *Enum) UnmarshalText(text []byte) error {
	value := string(text)
	for _, a := range e.allowed {
		if value == a {
			e.value = value
			return nil
		}
	}
	if len(e.allowed) == 0 {
		// If no allowed values set yet, just accept any value
		e.value = value
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

// Ptr returns a pointer to any value.
// Useful for optional config fields.
func Ptr[T any](v T) *T {
	return &v
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
		return fmt.Errorf("%s: must not be nil", name)
	}
	return nil
}
