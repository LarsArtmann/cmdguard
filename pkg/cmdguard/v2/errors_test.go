package v2

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSentinelErrors(t *testing.T) {
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
			assert.NotNil(t, tt.err)
			assert.Contains(t, tt.err.Error(), tt.msg)
		})
	}
}

func TestCommandError(t *testing.T) {
	t.Run("Error message", func(t *testing.T) {
		err := NewCommandError("test-cmd", ErrMissingHandler)
		assert.Contains(t, err.Error(), "test-cmd")
		assert.Contains(t, err.Error(), "command has no handler")
	})

	t.Run("Unwrap", func(t *testing.T) {
		err := NewCommandError("test-cmd", ErrMissingHandler)
		unwrapped := err.Unwrap()
		assert.True(t, errors.Is(unwrapped, ErrMissingHandler))
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewCommandError("test-cmd", ErrInvalidCommand)
		assert.True(t, errors.Is(err, ErrInvalidCommand))
	})
}

func TestFlagError(t *testing.T) {
	t.Run("Error message", func(t *testing.T) {
		err := NewFlagError("test-flag", errors.New("invalid value"))
		assert.Contains(t, err.Error(), "test-flag")
		assert.Contains(t, err.Error(), "invalid value")
	})

	t.Run("Unwrap", func(t *testing.T) {
		inner := errors.New("inner error")
		err := NewFlagError("test-flag", inner)
		unwrapped := err.Unwrap()
		assert.Equal(t, inner, unwrapped)
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewFlagError("test-flag", ErrFlagParseFailed)
		assert.True(t, errors.Is(err, ErrFlagParseFailed))
	})
}

func TestConfigError(t *testing.T) {
	t.Run("Error message", func(t *testing.T) {
		err := NewConfigError("LogLevel", errors.New("must be one of debug,info,warn,error"))
		assert.Contains(t, err.Error(), "LogLevel")
		assert.Contains(t, err.Error(), "must be one of")
	})

	t.Run("Unwrap", func(t *testing.T) {
		inner := errors.New("inner error")
		err := NewConfigError("field", inner)
		unwrapped := err.Unwrap()
		assert.Equal(t, inner, unwrapped)
	})
}

func TestEnumError(t *testing.T) {
	t.Run("Error message", func(t *testing.T) {
		err := NewEnumError("invalid", []string{"valid1", "valid2"})
		assert.Contains(t, err.Error(), "invalid")
		assert.Contains(t, err.Error(), "valid1")
		assert.Contains(t, err.Error(), "valid2")
	})

	t.Run("Unwrap", func(t *testing.T) {
		err := NewEnumError("invalid", []string{"valid"})
		unwrapped := err.Unwrap()
		assert.True(t, errors.Is(unwrapped, ErrInvalidEnum))
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewEnumError("invalid", []string{"valid"})
		assert.True(t, errors.Is(err, ErrInvalidEnum))
	})
}

func TestDurationError(t *testing.T) {
	t.Run("Error message", func(t *testing.T) {
		err := NewDurationError("not-a-duration", errors.New("time: invalid duration"))
		assert.Contains(t, err.Error(), "not-a-duration")
		assert.Contains(t, err.Error(), "time: invalid duration")
	})

	t.Run("Unwrap", func(t *testing.T) {
		inner := errors.New("inner error")
		err := NewDurationError("bad", inner)
		unwrapped := err.Unwrap()
		assert.True(t, errors.Is(unwrapped, ErrInvalidDuration))
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewDurationError("bad", errors.New("inner"))
		assert.True(t, errors.Is(err, ErrInvalidDuration))
	})
}

func TestErrorChaining(t *testing.T) {
	t.Run("nested errors", func(t *testing.T) {
		enumErr := NewEnumError("bad", []string{"good"})
		flagErr := NewFlagError("my-flag", enumErr)
		cmdErr := NewCommandError("my-cmd", flagErr)

		// All levels should be accessible via errors.Is
		assert.True(t, errors.Is(cmdErr, ErrInvalidEnum))
		require.NotNil(t, cmdErr)
	})
}

func TestServiceError(t *testing.T) {
	t.Run("Error message", func(t *testing.T) {
		inner := errors.New("service not initialized")
		err := NewServiceError("*DatabaseService", inner)
		assert.Contains(t, err.Error(), "*DatabaseService")
		assert.Contains(t, err.Error(), "service not initialized")
	})

	t.Run("Unwrap", func(t *testing.T) {
		inner := errors.New("inner error")
		err := NewServiceError("*LoggerService", inner)
		unwrapped := err.Unwrap()
		assert.Equal(t, inner, unwrapped)
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewServiceError("*DatabaseService", ErrServiceNotFound)
		assert.True(t, errors.Is(err, ErrServiceNotFound))
	})
}
