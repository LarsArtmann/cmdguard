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
