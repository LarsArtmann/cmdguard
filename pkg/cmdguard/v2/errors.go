// Package v2 provides a type-safe, DI-powered CLI framework with no panics.
// All functions return errors instead of panicking.
package v2

import (
	"errors"
	"fmt"
)

// Sentinel errors for type checking with errors.Is().
var (
	// ErrInvalidCommand indicates a command is malformed or missing required fields.
	ErrInvalidCommand = errors.New("invalid command")

	// ErrMissingHandler indicates a command has no RunE handler and no subcommands.
	ErrMissingHandler = errors.New("command has no handler")

	// ErrMissingName indicates a command has no name (Use field is empty).
	ErrMissingName = errors.New("command has no name")

	// ErrMissingLong indicates a parent command has no long description.
	ErrMissingLong = errors.New("command has no long description")

	// ErrMissingShort indicates a command has no short description.
	ErrMissingShort = errors.New("command has no short description")

	// ErrTooFewArgs indicates a command received fewer positional arguments than required.
	ErrTooFewArgs = errors.New("too few arguments")

	// ErrTooManyArgs indicates a command received more positional arguments than allowed.
	ErrTooManyArgs = errors.New("too many arguments")

	// ErrFlagParseFailed indicates flag parsing failed.
	ErrFlagParseFailed = errors.New("failed to parse flags")

	// ErrConfigValidation indicates config validation failed.
	ErrConfigValidation = errors.New("config validation failed")

	// ErrDuplicateCommand indicates a command with the same name already exists.
	ErrDuplicateCommand = errors.New("duplicate command")

	// ErrInvalidScope indicates an invalid DI scope operation.
	ErrInvalidScope = errors.New("invalid scope")

	// ErrServiceNotFound indicates a service was not found in the DI container.
	ErrServiceNotFound = errors.New("service not found")

	// ErrServiceConstruction indicates a service provider failed during construction.
	ErrServiceConstruction = errors.New("service construction failed")

	// ErrServiceRegistration indicates a service failed to register in the DI container.
	ErrServiceRegistration = errors.New("service registration failed")

	// ErrInvalidEnum indicates an invalid enum value.
	ErrInvalidEnum = errors.New("invalid enum value")

	// ErrInvalidDuration indicates an invalid duration format.
	ErrInvalidDuration = errors.New("invalid duration")

	// ErrInvalidFlagType indicates the flag type F is not a struct or pointer to struct.
	ErrInvalidFlagType = errors.New("invalid flag type")

	// ErrConfigNil indicates a nil config was passed where a config struct is required.
	ErrConfigNil = errors.New("config must not be nil")

	// ErrFlagNotFound indicates a flag was not found in the command.
	ErrFlagNotFound = errors.New("flag not found")

	// ErrRequiredFlag indicates a required flag was not set.
	ErrRequiredFlag = errors.New("required flag not set")

	// ErrConfigNotPointer indicates config is not a pointer to struct.
	ErrConfigNotPointer = errors.New("config must be a pointer to struct")

	// ErrNoFlags indicates the command has no flags to register.
	ErrNoFlags = errors.New("no flags to register")

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

	// ErrFieldNotFound indicates a struct field was not found by name.
	ErrFieldNotFound = errors.New("field not found")

	// ErrFieldNotSettable indicates a struct field is unexported or otherwise not settable.
	ErrFieldNotSettable = errors.New("field is not settable")

	// ErrTypeConversion indicates a type conversion or assertion failed.
	ErrTypeConversion = errors.New("type conversion failed")

	// ErrUnsupportedConversion indicates a string-to-type conversion is not supported.
	ErrUnsupportedConversion = errors.New("unsupported string conversion")

	// ErrNilValue indicates a required value was nil.
	ErrNilValue = errors.New("value must not be nil")

	// ErrFlagInstance indicates a flag struct instance could not be created.
	ErrFlagInstance = errors.New("failed to create flag instance")

	// ErrLogLevel indicates an invalid log level value.
	ErrLogLevel = errors.New("invalid log level")

	// ErrLogFormat indicates an invalid log format value.
	ErrLogFormat = errors.New("invalid log format")

	// ErrValueTooShort indicates a string value is shorter than required.
	ErrValueTooShort = errors.New("value too short")

	// ErrValueTooLong indicates a string value is longer than allowed.
	ErrValueTooLong = errors.New("value too long")

	// ErrValueTooSmall indicates a numeric value is smaller than minimum.
	ErrValueTooSmall = errors.New("value too small")

	// ErrValueTooLarge indicates a numeric value is larger than maximum.
	ErrValueTooLarge = errors.New("value too large")

	// ErrValuePatternMismatch indicates a value does not match the required pattern.
	ErrValuePatternMismatch = errors.New("value does not match pattern")

	// ErrValueEmpty indicates a value is empty but is required to be non-empty.
	ErrValueEmpty = errors.New("value is empty")

	// ErrUnknownValidator indicates a validator name in a validate tag is not registered.
	ErrUnknownValidator = errors.New("unknown validator")

	// ErrInvalidValidatorParam indicates a validator parameter could not be parsed.
	ErrInvalidValidatorParam = errors.New("invalid validator parameter")

	// ErrCommandPanic indicates a command handler panicked during execution.
	ErrCommandPanic = errors.New("command panicked")

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

	// ErrNegativeArgCount indicates a negative argument count was provided.
	ErrNegativeArgCount = errors.New("argument count must not be negative")

	// ErrInvalidArgRange indicates min > max in a range argument validator.
	ErrInvalidArgRange = errors.New("minimum argument count must not exceed maximum")

	// ErrInvalidExitCode indicates an exit code outside the valid 0–255 range.
	ErrInvalidExitCode = errors.New("exit code must be between 0 and 255")

	// ErrMissingExample indicates a leaf command has no example in draconian mode.
	ErrMissingExample = errors.New("command has no example")
)

// CommandError wraps an error with command context.
type CommandError struct {
	CommandName string
	Err         error
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("command %q: %v", e.CommandName, e.Err)
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
	msg := fmt.Sprintf("flag %q: %v", e.FlagName, e.Err)
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
	return fmt.Sprintf("config field %q: %v", e.Field, e.Err)
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
	return fmt.Sprintf("invalid duration %q: %v", e.Value, e.Err)
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
	return fmt.Sprintf("service %q: %v", e.ServiceType, e.Err)
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
	ExitCode() int
}

// ExitError wraps an error with a specific exit code.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
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
