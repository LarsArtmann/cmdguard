package v2

import "log/slog"

// logLevelAllowed is the set of valid log levels.
//
//nolint:gochecknoglobals // Shared constant for DRY
var logLevelAllowed = []string{"debug", "info", "warn", "error"}

// logFormatAllowed is the set of valid log formats.
//
//nolint:gochecknoglobals // Shared constant for DRY
var logFormatAllowed = []string{"text", "json"}

// LogLevel is a type-safe log level enum.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type LogLevel Enum

// Log level constants.
//
//nolint:gochecknoglobals // Predefined constants for type-safe defaults
var (
	LogLevelDebug = LogLevel{value: "debug", allowed: logLevelAllowed}
	LogLevelInfo  = LogLevel{value: "info", allowed: logLevelAllowed}
	LogLevelWarn  = LogLevel{value: "warn", allowed: logLevelAllowed}
	LogLevelError = LogLevel{value: "error", allowed: logLevelAllowed}
)

// ParseLogLevel creates a LogLevel from a string.
func ParseLogLevel(s string) (LogLevel, error) {
	e, err := ParseEnum(s, logLevelAllowed)
	if err != nil {
		return LogLevel{}, err
	}

	return LogLevel(e), nil
}

// MustParseLogLevel creates a LogLevel from a string, panicking if invalid.
// Use only when you know the level is valid (e.g., for constants).
func MustParseLogLevel(s string) LogLevel {
	return MustParse("MustParseLogLevel", s, ParseLogLevel)
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

// IsEmpty returns true if the log level is unset.
func (l LogLevel) IsEmpty() bool {
	return l.value == ""
}

// LogFormat is a type-safe log format enum.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type LogFormat Enum

// Log format constants.
//
//nolint:gochecknoglobals // Predefined constants for type-safe defaults
var (
	LogFormatText = LogFormat{value: "text", allowed: logFormatAllowed}
	LogFormatJSON = LogFormat{value: "json", allowed: logFormatAllowed}
)

// ParseLogFormat creates a LogFormat from a string.
func ParseLogFormat(s string) (LogFormat, error) {
	e, err := ParseEnum(s, logFormatAllowed)
	if err != nil {
		return LogFormat{}, err
	}

	return LogFormat(e), nil
}

// MustParseLogFormat creates a LogFormat from a string, panicking if invalid.
// Use only when you know the format is valid (e.g., for constants).
func MustParseLogFormat(s string) LogFormat {
	return MustParse("MustParseLogFormat", s, ParseLogFormat)
}

// String returns the log format as a string.
func (f LogFormat) String() string {
	return f.value
}

// IsEmpty returns true if the log format is unset.
func (f LogFormat) IsEmpty() bool {
	return f.value == ""
}

// MarshalText implements encoding.TextMarshaler for LogLevel.
func (l LogLevel) MarshalText() ([]byte, error) {
	return textMarshal(l, LogLevel.String)
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	return textUnmarshal(l, text, ParseLogLevel)
}

// MarshalText implements encoding.TextMarshaler for LogFormat.
func (f LogFormat) MarshalText() ([]byte, error) {
	return textMarshal(f, LogFormat.String)
}

func (f *LogFormat) UnmarshalText(text []byte) error {
	return textUnmarshal(f, text, ParseLogFormat)
}
