package v2

import (
	"errors"
	"fmt"
)

// Type-specific sentinel errors (remaining from the original errors.go split).
// Command errors → errors_command.go, Flag errors → errors_flags.go,
// Config errors → errors_config.go, DI errors → errors_di.go.
var (
	// ErrInvalidDuration indicates an invalid duration format.
	ErrInvalidDuration = errors.New("invalid duration")

	// ErrInvalidURL indicates an invalid URL format.
	ErrInvalidURL = errors.New("invalid URL")

	// ErrInvalidEmail indicates an invalid email address format.
	ErrInvalidEmail = errors.New("invalid email address")

	// ErrInvalidPort indicates an invalid port number.
	ErrInvalidPort = errors.New("invalid port")

	// ErrInvalidFilePath indicates an invalid file path.
	ErrInvalidFilePath = errors.New("invalid file path")

	// ErrInvalidHostPort indicates an invalid host:port format.
	ErrInvalidHostPort = errors.New("invalid host:port")

	// ErrLogLevel indicates an invalid log level value.
	ErrLogLevel = errors.New("invalid log level")

	// ErrLogFormat indicates an invalid log format value.
	ErrLogFormat = errors.New("invalid log format")

	// ErrUnsupportedFormat indicates the requested output format is not supported.
	ErrUnsupportedFormat = errors.New("unsupported output format")

	// ErrFormatRequiresTypedData indicates the format requires structured data.
	ErrFormatRequiresTypedData = errors.New("format requires typed data")

	// ErrMissingVersion indicates a version command was requested but no version is set.
	ErrMissingVersion = errors.New("version is required but not set")

	// ErrEditorTempFile indicates a temporary file could not be created for editing.
	ErrEditorTempFile = errors.New("failed to create temp file for editor")

	// ErrEditorWrite indicates writing content to the temp file failed.
	ErrEditorWrite = errors.New("failed to write to temp file")

	// ErrEditorRun indicates the editor process failed.
	ErrEditorRun = errors.New("editor execution failed")

	// ErrEditorRead indicates reading the edited file failed.
	ErrEditorRead = errors.New("failed to read edited file")

	// ErrInvalidExitCode indicates an exit code outside the valid 0–255 range.
	ErrInvalidExitCode = errors.New("exit code must be between 0 and 255")

	// ErrDoctorFailed indicates one or more doctor checks failed.
	ErrDoctorFailed = errors.New("doctor checks failed")

)

// labeledError formats an error with a labeled context for consistent error messages.
func labeledError(label, value string, err error) string {
	return fmt.Sprintf("%s %q: %v", label, value, err)
}

// CommandError wraps an error with command context.
type CommandError struct {
	CommandName string
	Err         error
}

func (e *CommandError) Error() string {
	return labeledError("command", e.CommandName, e.Err)
}

func (e *CommandError) Unwrap() error {
	return e.Err
}

// NewCommandError creates a new CommandError.
func NewCommandError(name string, err error) *CommandError {
	return &CommandError{
		CommandName: name,
		Err:         err,
	}
}

// FlagError wraps an error with flag context.
type FlagError struct {
	FlagName   string
	Err        error
	Suggestion string
}

func (e *FlagError) Error() string {
	msg := labeledError("flag", e.FlagName, e.Err)
	if e.Suggestion != "" {
		msg += fmt.Sprintf(" (did you mean --%s?)", e.Suggestion)
	}

	return msg
}

func (e *FlagError) Unwrap() error {
	return e.Err
}

// NewFlagError creates a new FlagError.
func NewFlagError(name string, err error) *FlagError {
	return &FlagError{
		FlagName: name,
		Err:      err,
	}
}

// NewFlagErrorWithSuggestion creates a new FlagError with a suggestion.
func NewFlagErrorWithSuggestion(name string, err error, suggestion string) *FlagError {
	return &FlagError{
		FlagName:   name,
		Err:        err,
		Suggestion: suggestion,
	}
}

// ConfigError wraps an error with config field context.
type ConfigError struct {
	Field string
	Err   error
}

func (e *ConfigError) Error() string {
	return labeledError("config field", e.Field, e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new ConfigError.
func NewConfigError(field string, err error) *ConfigError {
	return &ConfigError{
		Field: field,
		Err:   err,
	}
}

// EnumError indicates an invalid enum value was provided.
type EnumError struct {
	Value   string
	Allowed []string
}

func (e *EnumError) Error() string {
	return fmt.Sprintf("invalid value %q, must be one of: %v", e.Value, e.Allowed)
}

func (e *EnumError) Unwrap() error {
	return ErrInvalidEnum
}

// NewEnumError creates a new EnumError.
func NewEnumError(value string, allowed []string) *EnumError {
	return &EnumError{
		Value:   value,
		Allowed: allowed,
	}
}

// DurationError indicates an invalid duration format.
type DurationError struct {
	Value string
	Err   error
}

func (e *DurationError) Error() string {
	return labeledError("invalid duration", e.Value, e.Err)
}

func (e *DurationError) Unwrap() error {
	return ErrInvalidDuration
}

// NewDurationError creates a new DurationError.
func NewDurationError(value string, err error) *DurationError {
	return &DurationError{
		Value: value,
		Err:   err,
	}
}

// ServiceError wraps a DI service error with type context.
// Use this when service invocation or construction fails.
type ServiceError struct {
	ServiceType string
	Err         error
}

func (e *ServiceError) Error() string {
	return labeledError("service", e.ServiceType, e.Err)
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError creates a new ServiceError.
func NewServiceError(serviceType string, err error) *ServiceError {
	return &ServiceError{ServiceType: serviceType, Err: err}
}

// ExitCoder is an interface that errors can implement to provide a specific exit code.
// When ExecuteAndExit encounters an error implementing ExitCoder, it uses the returned
// code instead of the default 1.
type ExitCoder interface {
	error
	ExitCode() int
}

// ExitError wraps an error with a specific exit code.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exit code %d", e.Code)
	}

	return e.Err.Error()
}

func (e *ExitError) Unwrap() error {
	return e.Err
}

func (e *ExitError) ExitCode() int {
	return e.Code
}

// NewExitError creates a new ExitError with the given code and cause.
// Returns an error if the code is outside the valid range 0-255.
func NewExitError(code int, err error) (*ExitError, error) {
	if code < 0 || code > 255 {
		return nil, fmt.Errorf("%w: got %d", ErrInvalidExitCode, code)
	}

	return &ExitError{Code: code, Err: err}, nil
}
