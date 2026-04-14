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
			t.Parallel()

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
	t.Parallel()
	t.Run("Error message", func(t *testing.T) {
		t.Parallel()

		err := NewCommandError("test-cmd", ErrMissingHandler)
		assertErrorContains(t, err, "test-cmd", "command has no handler")
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		err := NewCommandError("test-cmd", ErrMissingHandler)

		unwrapped := err.Unwrap()
		if !errors.Is(unwrapped, ErrMissingHandler) {
			t.Errorf("expected unwrapped error to be ErrMissingHandler, got %v", unwrapped)
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		t.Parallel()

		err := NewCommandError("test-cmd", ErrInvalidCommand)
		if !errors.Is(err, ErrInvalidCommand) {
			t.Errorf("expected error to match ErrInvalidCommand")
		}
	})
}

func TestFlagError(t *testing.T) {
	t.Parallel()

	innerErr := errors.New("invalid value")

	t.Run("Error message", func(t *testing.T) {
		t.Parallel()

		err := NewFlagError("test-flag", innerErr)
		assertErrorContains(t, err, "test-flag", "invalid value")
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		unwrapped := NewFlagError("test-flag", innerErr).Unwrap()
		if !errors.Is(unwrapped, innerErr) {
			t.Errorf("expected unwrapped error to be %v, got %v", innerErr, unwrapped)
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		t.Parallel()

		err := NewFlagError("test-flag", ErrFlagParseFailed)
		if !errors.Is(err, ErrFlagParseFailed) {
			t.Errorf("expected error to match ErrFlagParseFailed")
		}
	})
}

func TestConfigError(t *testing.T) {
	t.Parallel()

	innerErr := errors.New("must be one of debug,info,warn,error")

	t.Run("Error message", func(t *testing.T) {
		t.Parallel()

		err := NewConfigError("LogLevel", innerErr)
		assertErrorContains(t, err, "LogLevel", "must be one of")
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		unwrapped := NewConfigError("field", innerErr).Unwrap()
		if !errors.Is(unwrapped, innerErr) {
			t.Errorf("expected unwrapped error to be %v, got %v", innerErr, unwrapped)
		}
	})
}

func TestEnumError(t *testing.T) {
	t.Parallel()
	t.Run("Error message", func(t *testing.T) {
		t.Parallel()

		err := NewEnumError("invalid", []string{"valid1", "valid2"})
		assertStringContains(t, err.Error(), "invalid", "valid1", "valid2")
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		err := NewEnumError("invalid", []string{"valid"})

		unwrapped := err.Unwrap()
		if !errors.Is(unwrapped, ErrInvalidEnum) {
			t.Errorf("expected unwrapped error to be ErrInvalidEnum")
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		t.Parallel()

		err := NewEnumError("invalid", []string{"valid"})
		if !errors.Is(err, ErrInvalidEnum) {
			t.Errorf("expected error to match ErrInvalidEnum")
		}
	})
}

func TestDurationError(t *testing.T) {
	t.Parallel()

	innerErr := errors.New("time: invalid duration")

	t.Run("Error message", func(t *testing.T) {
		t.Parallel()

		err := NewDurationError("not-a-duration", innerErr)
		assertErrorContains(t, err, "not-a-duration", "time: invalid duration")
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		unwrapped := NewDurationError("bad", innerErr).Unwrap()
		if !errors.Is(unwrapped, ErrInvalidDuration) {
			t.Errorf("expected unwrapped error to be ErrInvalidDuration")
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		t.Parallel()

		err := NewDurationError("bad", innerErr)
		if !errors.Is(err, ErrInvalidDuration) {
			t.Errorf("expected error to match ErrInvalidDuration")
		}
	})
}

func TestErrorChaining(t *testing.T) {
	t.Parallel()
	t.Run("nested errors", func(t *testing.T) {
		t.Parallel()

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
	t.Parallel()

	innerErr := errors.New("service not initialized")

	t.Run("Error message", func(t *testing.T) {
		t.Parallel()

		err := NewServiceError("*DatabaseService", innerErr)
		assertStringContains(t, err.Error(), "*DatabaseService", "service not initialized")
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		unwrapped := NewServiceError("*LoggerService", innerErr).Unwrap()
		if !errors.Is(unwrapped, innerErr) {
			t.Errorf("expected unwrapped error to be %v, got %v", innerErr, unwrapped)
		}
	})

	t.Run("errors.Is support", func(t *testing.T) {
		t.Parallel()

		err := NewServiceError("*DatabaseService", ErrServiceNotFound)
		if !errors.Is(err, ErrServiceNotFound) {
			t.Errorf("expected error to match ErrServiceNotFound")
		}
	})
}
