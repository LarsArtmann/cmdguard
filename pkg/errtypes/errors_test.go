package errtypes

import (
	"testing"
)

func TestNewCodedError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		code    string
	}{
		{name: "both fields", message: "something failed", code: "E001"},
		{name: "empty message", message: "", code: "E002"},
		{name: "empty code", message: "something failed", code: ""},
		{name: "both empty", message: "", code: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewCodedError(tt.message, tt.code)
			if err == nil {
				t.Fatal("expected non-nil error")
			}

			if err.Message != tt.message {
				t.Errorf("Message = %q, want %q", err.Message, tt.message)
			}

			if err.Code != tt.code {
				t.Errorf("Code = %q, want %q", err.Code, tt.code)
			}
		})
	}
}

func TestCodedError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		want    string
	}{
		{name: "non-empty message", message: "something failed", want: "something failed"},
		{name: "empty message", message: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewCodedError(tt.message, "E001")
			if got := err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCodedError_ImplementsError(t *testing.T) {
	t.Parallel()

	var err error = NewCodedError("test", "E001")
	if err.Error() != "test" {
		t.Errorf("expected %q, got %q", "test", err.Error())
	}
}
