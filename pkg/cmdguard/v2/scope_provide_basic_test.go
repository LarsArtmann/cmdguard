package v2

import (
	"testing"

	"github.com/samber/do/v2"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestProvide(t *testing.T) {
	t.Parallel()
	t.Run("registers service provider", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("test")

		err := Provide(scope, func(_ do.Injector) (string, error) {
			return "test-value", nil
		})
		testutil.AssertNoError(t, err)

		value := mustInvoke[string](t, scope)

		testutil.AssertFieldEqString(t, value, "test-value", "value")
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		t.Parallel()

		err := Provide(nil, func(_ do.Injector) (string, error) {
			return "value", nil
		})
		testutil.AssertExpectedError(t, err)

		assertErrorContains(t, err, "scope is nil")
	})

	t.Run("provider can use dependencies", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("test")

		// Register a dependency
		type Dep string

		mustProvideValue(t, scope, Dep("dependency"))

		type Service string

		err := Provide(scope, func(i do.Injector) (Service, error) {
			dep, invokeErr := do.Invoke[Dep](i)
			if invokeErr != nil {
				return "", invokeErr
			}

			return Service(dep + "-enhanced"), nil
		})
		testutil.AssertNoError(t, err)

		value := mustInvoke[Service](t, scope)

		testutil.AssertFieldEqString(t, string(value), "dependency-enhanced", "value")
	})
}

func TestProvideValue(t *testing.T) {
	t.Parallel()
	t.Run("registers value directly", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("test")

		err := ProvideValue(scope, 42)
		testutil.AssertNoError(t, err)

		value := mustInvoke[int](t, scope)

		testutil.AssertFieldEq(t, value, 42, "value")
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		t.Parallel()

		err := ProvideValue(nil, 42)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "scope is nil")
	})

	t.Run("can register struct values", func(t *testing.T) {
		t.Parallel()

		type Config struct {
			Name string
			Port int
		}

		scope := NewScope("test")

		cfg := Config{Name: "app", Port: 8080}
		mustProvideValue(t, scope, cfg)

		value := mustInvoke[Config](t, scope)

		testutil.AssertFieldEqString(t, value.Name, "app", "Name")
		testutil.AssertFieldEq(t, value.Port, 8080, "Port")
	})
}

func TestInvoke(t *testing.T) {
	t.Parallel()
	t.Run("returns registered service", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("test")
		mustProvideValue(t, scope, "hello")

		value := mustInvoke[string](t, scope)

		testutil.AssertFieldEqString(t, value, "hello", "value")
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		t.Parallel()

		value, err := Invoke[string](nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "scope is nil")

		if value != "" {
			t.Errorf("expected empty value, got %q", value)
		}
	})

	t.Run("returns error for unregistered service", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		scope := NewScope("test")
		mustProvideValue(t, scope, 123)

		mustProvideValue(t, scope, "text")

		mustProvideValue(t, scope, true)

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
