package v2

import (
	"strings"
	"testing"

	"github.com/samber/do/v2"
)

func TestProvide(t *testing.T) {
	t.Run("registers service provider", func(t *testing.T) {
		scope := NewScope("test")

		err := Provide(scope, func(i do.Injector) (string, error) {
			return "test-value", nil
		})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		value, err := Invoke[string](scope)
		if err != nil {
			t.Fatalf("expected no error invoking, got: %v", err)
		}

		if value != "test-value" {
			t.Errorf("expected value 'test-value', got %q", value)
		}
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		err := Provide(nil, func(i do.Injector) (string, error) {
			return "value", nil
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "scope is nil") {
			t.Errorf("expected error to contain 'scope is nil', got: %v", err)
		}
	})

	t.Run("provider can use dependencies", func(t *testing.T) {
		scope := NewScope("test")

		// Register a dependency
		type Dep string

		if err := ProvideValue(scope, Dep("dependency")); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		// Register a service that uses the dependency
		type Service string

		err := Provide(scope, func(i do.Injector) (Service, error) {
			dep, invokeErr := do.Invoke[Dep](i)
			if invokeErr != nil {
				return "", invokeErr
			}

			return Service(dep + "-enhanced"), nil
		})
		if err != nil {
			t.Fatalf("expected no error providing, got: %v", err)
		}

		value, err := Invoke[Service](scope)
		if err != nil {
			t.Fatalf("expected no error invoking, got: %v", err)
		}

		if value != Service("dependency-enhanced") {
			t.Errorf("expected value 'dependency-enhanced', got %q", value)
		}
	})
}

func TestProvideValue(t *testing.T) {
	t.Run("registers value directly", func(t *testing.T) {
		scope := NewScope("test")

		err := ProvideValue(scope, 42)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		value, err := Invoke[int](scope)
		if err != nil {
			t.Fatalf("expected no error invoking, got: %v", err)
		}

		if value != 42 {
			t.Errorf("expected value 42, got %d", value)
		}
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		err := ProvideValue(nil, 42)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "scope is nil") {
			t.Errorf("expected error to contain 'scope is nil', got: %v", err)
		}
	})

	t.Run("can register struct values", func(t *testing.T) {
		type Config struct {
			Name string
			Port int
		}

		scope := NewScope("test")

		cfg := Config{Name: "app", Port: 8080}
		if err := ProvideValue(scope, cfg); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		value, err := Invoke[Config](scope)
		if err != nil {
			t.Fatalf("expected no error invoking, got: %v", err)
		}

		if value.Name != "app" {
			t.Errorf("expected name 'app', got %q", value.Name)
		}

		if value.Port != 8080 {
			t.Errorf("expected port 8080, got %d", value.Port)
		}
	})
}

func TestInvoke(t *testing.T) {
	t.Run("returns registered service", func(t *testing.T) {
		scope := NewScope("test")
		if err := ProvideValue(scope, "hello"); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		value, err := Invoke[string](scope)
		if err != nil {
			t.Fatalf("expected no error invoking, got: %v", err)
		}

		if value != "hello" {
			t.Errorf("expected value 'hello', got %q", value)
		}
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		value, err := Invoke[string](nil)
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

	t.Run("returns error for unregistered service", func(t *testing.T) {
		scope := NewScope("test")

		value, err := Invoke[string](scope)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if value != "" {
			t.Errorf("expected empty value, got %q", value)
		}
	})

	t.Run("can invoke different types", func(t *testing.T) {
		scope := NewScope("test")
		if err := ProvideValue(scope, 123); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		if err := ProvideValue(scope, "text"); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		if err := ProvideValue(scope, true); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		intVal, err := Invoke[int](scope)
		if err != nil {
			t.Fatalf("expected no error invoking int, got: %v", err)
		}

		if intVal != 123 {
			t.Errorf("expected int value 123, got %d", intVal)
		}

		strVal, err := Invoke[string](scope)
		if err != nil {
			t.Fatalf("expected no error invoking string, got: %v", err)
		}

		if strVal != "text" {
			t.Errorf("expected string value 'text', got %q", strVal)
		}

		boolVal, err := Invoke[bool](scope)
		if err != nil {
			t.Fatalf("expected no error invoking bool, got: %v", err)
		}

		if !boolVal {
			t.Error("expected bool value true")
		}
	})
}
