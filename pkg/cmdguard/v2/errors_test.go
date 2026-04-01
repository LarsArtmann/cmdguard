package v2

import (
	"errors"
	"strings"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		err   error
		msg   string
		isNil bool
	}{
		{"ErrInvalidCommand", ErrInvalidCommand, "invalid command", false},
		{"ErrMissingHandler", ErrMissingHandler, "command has no handler", false},
		{"ErrMissingName", ErrMissingName, "command has no name", false},
		{"ErrFlagParseFailed", ErrFlagParseFailed, "failed to parse flags", false},
		{"ErrConfigValidation", ErrConfigValidation, "config validation failed", false},
		{"ErrDuplicateCommand", ErrDuplicateCommand, "duplicate command", false},
		{"ErrInvalidScope", ErrInvalidScope, "invalid scope", false},
		{"ErrServiceNotFound", ErrServiceNotFound, "service not found", false},
		{"ErrInvalidEnum", ErrInvalidEnum, "invalid enum value", false},
		{"ErrInvalidDuration", ErrInvalidDuration, "invalid duration", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatal("expected non-nil error")
			}

			if !strings.Contains(tt.err.Error(), tt.msg) {
				t.Errorf("error %q does not contain %q", tt.err.Error(), tt.msg)
			}
		})
	}
}

func TestCommandError(t *testing.T) {
	t.Run("Error message", func(t *testing.T) {
		err := NewCommandError("test-cmd", ErrMissingHandler)
		if !strings.Contains(err.Error(), "test-cmd") {
			t.Errorf("expected error to contain 'test-cmd', got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "command has no handler") {
			t.Errorf("expected error to contain 'command has no handler', got %q", err.Error())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		err := NewCommandError("test-cmd", ErrMissingHandler)

		unwrapped := err.Unwrap()
		if !errors.Is(unwrapped, ErrMissingHandler) {
			t.Errorf("expected unwrapped error to be ErrMissingHandler, got %v", unwrapped)
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewCommandError("test-cmd", ErrInvalidCommand)
		if !errors.Is(err, ErrInvalidCommand) {
			t.Errorf("expected error to match ErrInvalidCommand")
		}
	})
}

func TestFlagError(t *testing.T) {
	innerErr := errors.New("invalid value")

	t.Run("Error message", func(t *testing.T) {
		err := NewFlagError("test-flag", innerErr)
		if !strings.Contains(err.Error(), "test-flag") {
			t.Errorf("expected error to contain 'test-flag', got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "invalid value") {
			t.Errorf("expected error to contain 'invalid value', got %q", err.Error())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		unwrapped := NewFlagError("test-flag", innerErr).Unwrap()
		if !errors.Is(unwrapped, innerErr) {
			t.Errorf("expected unwrapped error to be %v, got %v", innerErr, unwrapped)
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewFlagError("test-flag", ErrFlagParseFailed)
		if !errors.Is(err, ErrFlagParseFailed) {
			t.Errorf("expected error to match ErrFlagParseFailed")
		}
	})
}

func TestConfigError(t *testing.T) {
	innerErr := errors.New("must be one of debug,info,warn,error")

	t.Run("Error message", func(t *testing.T) {
		err := NewConfigError("LogLevel", innerErr)
		if !strings.Contains(err.Error(), "LogLevel") {
			t.Errorf("expected error to contain 'LogLevel', got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "must be one of") {
			t.Errorf("expected error to contain 'must be one of', got %q", err.Error())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		unwrapped := NewConfigError("field", innerErr).Unwrap()
		if !errors.Is(unwrapped, innerErr) {
			t.Errorf("expected unwrapped error to be %v, got %v", innerErr, unwrapped)
		}
	})
}

func TestEnumError(t *testing.T) {
	t.Run("Error message", func(t *testing.T) {
		err := NewEnumError("invalid", []string{"valid1", "valid2"})

		errMsg := err.Error()
		if !strings.Contains(errMsg, "invalid") {
			t.Errorf("expected error to contain 'invalid', got %q", errMsg)
		}

		if !strings.Contains(errMsg, "valid1") {
			t.Errorf("expected error to contain 'valid1', got %q", errMsg)
		}

		if !strings.Contains(errMsg, "valid2") {
			t.Errorf("expected error to contain 'valid2', got %q", errMsg)
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		err := NewEnumError("invalid", []string{"valid"})

		unwrapped := err.Unwrap()
		if !errors.Is(unwrapped, ErrInvalidEnum) {
			t.Errorf("expected unwrapped error to be ErrInvalidEnum")
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewEnumError("invalid", []string{"valid"})
		if !errors.Is(err, ErrInvalidEnum) {
			t.Errorf("expected error to match ErrInvalidEnum")
		}
	})
}

func TestDurationError(t *testing.T) {
	innerErr := errors.New("time: invalid duration")

	t.Run("Error message", func(t *testing.T) {
		err := NewDurationError("not-a-duration", innerErr)
		if !strings.Contains(err.Error(), "not-a-duration") {
			t.Errorf("expected error to contain 'not-a-duration', got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "time: invalid duration") {
			t.Errorf("expected error to contain 'time: invalid duration', got %q", err.Error())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		unwrapped := NewDurationError("bad", innerErr).Unwrap()
		if !errors.Is(unwrapped, ErrInvalidDuration) {
			t.Errorf("expected unwrapped error to be ErrInvalidDuration")
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewDurationError("bad", innerErr)
		if !errors.Is(err, ErrInvalidDuration) {
			t.Errorf("expected error to match ErrInvalidDuration")
		}
	})
}

func TestErrorChaining(t *testing.T) {
	t.Run("nested errors", func(t *testing.T) {
		enumErr := NewEnumError("bad", []string{"good"})
		flagErr := NewFlagError("my-flag", enumErr)
		cmdErr := NewCommandError("my-cmd", flagErr)

		if !errors.Is(cmdErr, ErrInvalidEnum) {
			t.Errorf("expected error chain to contain ErrInvalidEnum")
		}

		if cmdErr == nil {
			t.Error("expected non-nil command error")
		}
	})
}

func TestServiceError(t *testing.T) {
	innerErr := errors.New("service not initialized")

	t.Run("Error message", func(t *testing.T) {
		err := NewServiceError("*DatabaseService", innerErr)

		errMsg := err.Error()
		if !strings.Contains(errMsg, "*DatabaseService") {
			t.Errorf("expected error to contain '*DatabaseService', got %q", errMsg)
		}

		if !strings.Contains(errMsg, "service not initialized") {
			t.Errorf("expected error to contain 'service not initialized', got %q", errMsg)
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		unwrapped := NewServiceError("*LoggerService", innerErr).Unwrap()
		if !errors.Is(unwrapped, innerErr) {
			t.Errorf("expected unwrapped error to be %v, got %v", innerErr, unwrapped)
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewServiceError("*DatabaseService", ErrServiceNotFound)
		if !errors.Is(err, ErrServiceNotFound) {
			t.Errorf("expected error to match ErrServiceNotFound")
		}
	})
}
