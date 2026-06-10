package v2

import "time"

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

// MustParseDuration creates a Duration from a string, panicking if invalid.
// Use only when you know the duration is valid (e.g., for constants).
func MustParseDuration(s string) Duration {
	return MustParse("MustParseDuration", s, ParseDuration)
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

// IsEmpty returns true if the duration is zero.
func (d Duration) IsEmpty() bool {
	return d.IsZero()
}

// Milliseconds returns the duration in milliseconds.
func (d Duration) Milliseconds() int64 {
	return d.duration.Milliseconds()
}

// Seconds returns the duration in seconds as a float64.
func (d Duration) Seconds() float64 {
	return d.duration.Seconds()
}

// MarshalText implements encoding.TextMarshaler for Duration.
func (d Duration) MarshalText() ([]byte, error) {
	return textMarshal(d, Duration.String)
}

func (d *Duration) UnmarshalText(text []byte) error {
	return textUnmarshal(d, text, ParseDuration)
}
