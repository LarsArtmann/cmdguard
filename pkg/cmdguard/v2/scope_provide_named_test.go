package v2

import (
	"strings"
	"testing"

	"github.com/samber/do/v2"
)

// provideTestNamed is a helper that reduces boilerplate in tests by wrapping
// ProvideNamed with a simple value-returning provider function.
func provideTestNamed[T any](scope *Scope, name string, value T) error {
	return ProvideNamed(scope, name, func(_ do.Injector) (T, error) {
		return value, nil
	})
}

func TestProvideNamed(t *testing.T) {
	t.Run("registers named service provider", func(t *testing.T) {
		scope := NewScope("test")

		err := ProvideNamed(scope, "cache-memory", func(_ do.Injector) (string, error) {
			return "memory-cache", nil
		})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		value, err := InvokeNamed[string](scope, "cache-memory")
		if err != nil {
			t.Fatalf("expected no error invoking named, got: %v", err)
		}

		if value != "memory-cache" {
			t.Errorf("expected value 'memory-cache', got %q", value)
		}
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		err := ProvideNamed(nil, "name", func(_ do.Injector) (string, error) {
			return "value", nil
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "scope is nil") {
			t.Errorf("expected error to contain 'scope is nil', got: %v", err)
		}
	})

	t.Run("can register multiple named implementations", func(t *testing.T) {
		scope := NewScope("test")

		if err := provideTestNamed(scope, "impl1", "implementation-1"); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if err := provideTestNamed(scope, "impl2", "implementation-2"); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		val1, err := InvokeNamed[string](scope, "impl1")
		if err != nil {
			t.Fatalf("expected no error invoking impl1, got: %v", err)
		}

		val2, err := InvokeNamed[string](scope, "impl2")
		if err != nil {
			t.Fatalf("expected no error invoking impl2, got: %v", err)
		}

		if val1 != "implementation-1" {
			t.Errorf("expected 'implementation-1', got %q", val1)
		}

		if val2 != "implementation-2" {
			t.Errorf("expected 'implementation-2', got %q", val2)
		}
	})
}

func TestInvokeNamed(t *testing.T) {
	t.Run("returns named service", func(t *testing.T) {
		scope := NewScope("test")

		if err := provideTestNamed(scope, "my-service", 42); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		value, err := InvokeNamed[int](scope, "my-service")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if value != 42 {
			t.Errorf("expected value 42, got %d", value)
		}
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		value, err := InvokeNamed[string](nil, "name")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "scope is nil") {
			t.Errorf("expected error to contain 'scope is nil', got: %v", err)
		}

		if value != "" {
			t.Errorf("expected empty value, got %q", value)
		}
	})
}

func TestMustInvoke(t *testing.T) {
	t.Run("returns registered service", func(t *testing.T) {
		scope := NewScope("test")
		if err := ProvideValue(scope, "hello"); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		value := MustInvoke[string](scope)
		if value != "hello" {
			t.Errorf("expected value 'hello', got %q", value)
		}
	})

	t.Run("panics for nil scope", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()

		MustInvoke[string](nil)
	})

	t.Run("panics for unregistered service", func(t *testing.T) {
		scope := NewScope("test")

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()

		MustInvoke[string](scope)
	})
}

func TestMustInvokeNamed(t *testing.T) {
	t.Run("returns named service", func(t *testing.T) {
		scope := NewScope("test")

		if err := provideTestNamed(scope, "my-service", 99); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		value := MustInvokeNamed[int](scope, "my-service")
		if value != 99 {
			t.Errorf("expected value 99, got %d", value)
		}
	})

	t.Run("panics for nil scope", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()

		MustInvokeNamed[string](nil, "name")
	})

	t.Run("panics for unregistered named service", func(t *testing.T) {
		scope := NewScope("test")

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()

		MustInvokeNamed[string](scope, "nonexistent")
	})
}
