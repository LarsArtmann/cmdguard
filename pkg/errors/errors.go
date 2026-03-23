// Package errors provides custom error types and error handling utilities.
//
// This package defines application-specific error types that can be used
// throughout the codebase for consistent error handling and reporting.
package errors

// CodedError is an error with an associated code for categorization.
type CodedError struct {
	Message string
	Code    string
}

// Error implements the error interface.
func (e *CodedError) Error() string {
	return e.Message
}

// New creates a new CodedError with the given message and code.
func New(message, code string) *CodedError {
	return &CodedError{
		Message: message,
		Code:    code,
	}
}
