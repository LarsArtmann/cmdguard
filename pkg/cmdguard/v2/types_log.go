package v2

import "log/slog"

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
