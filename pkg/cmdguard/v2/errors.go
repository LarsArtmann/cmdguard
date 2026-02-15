// Package v2 provides a type-safe, DI-powered CLI framework with no panics.
// All functions return errors instead of panicking.
package v2

import (
	"errors"
	"fmt"
)

// Sentinel errors for type checking with errors.Is()
var (
	// ErrInvalidCommand indicates a command is malformed or missing required fields.
	ErrInvalidCommand = errors.New("invalid command")

	// ErrMissingHandler indicates a command has no RunE handler and no subcommands.
	ErrMissingHandler = errors.New("command has no handler")

	// ErrMissingName indicates a command has no name (Use field is empty).
	ErrMissingName = errors.New("command has no name")

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

	// ErrInvalidEnum indicates an invalid enum value.
	ErrInvalidEnum = errors.New("invalid enum value")

	// ErrInvalidDuration indicates an invalid duration format.
	ErrInvalidDuration = errors.New("invalid duration")

	// ErrInvalidFlagType indicates the flag type F is not a struct or pointer to struct.
	ErrInvalidFlagType = errors.New("invalid flag type")
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
	FlagName string
	Err      error
}

func (e *FlagError) Error() string {
	return fmt.Sprintf("flag %q: %v", e.FlagName, e.Err)
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
		Err:    err,
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
