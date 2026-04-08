// Package errtypes provides custom error types and error handling utilities.
//
// This package defines application-specific error types that can be used
// throughout the codebase for consistent error handling and reporting.
package errtypes

// CodedError is a foundation for custom errors with structured metadata.
type CodedError struct {
	Message string
	Code    string
}

// Error implements the error interface.
func (e *CodedError) Error() string {
	return e.Message
}

// NewCodedError creates a new CodedError with the given message and code.
func NewCodedError(message, code string) *CodedError {
	return &CodedError{
		Message: message,
		Code:    code,
	}
}
