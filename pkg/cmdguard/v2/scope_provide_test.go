package v2

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvide(t *testing.T) {
	t.Run("registers service provider", func(t *testing.T) {
		scope := NewScope("test")
		err := Provide(scope, func(i do.Injector) (string, error) {
			return "test-value", nil
		})
		require.NoError(t, err)

		value, err := Invoke[string](scope)
		require.NoError(t, err)
		assert.Equal(t, "test-value", value)
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		err := Provide[string](nil, func(i do.Injector) (string, error) {
			return "value", nil
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "scope is nil")
	})

	t.Run("provider can use dependencies", func(t *testing.T) {
		scope := NewScope("test")

		// Register a dependency
		type Dep string
		require.NoError(t, ProvideValue(scope, Dep("dependency")))

		// Register a service that uses the dependency
		type Service string
		err := Provide(scope, func(i do.Injector) (Service, error) {
			dep, err := do.Invoke[Dep](i)
			if err != nil {
				return "", err
			}
			return Service(dep + "-enhanced"), nil
		})
		require.NoError(t, err)

		value, err := Invoke[Service](scope)
		require.NoError(t, err)
		assert.Equal(t, Service("dependency-enhanced"), value)
	})
}

func TestProvideValue(t *testing.T) {
	t.Run("registers value directly", func(t *testing.T) {
		scope := NewScope("test")
		err := ProvideValue(scope, 42)
		require.NoError(t, err)

		value, err := Invoke[int](scope)
		require.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		err := ProvideValue[int](nil, 42)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "scope is nil")
	})

	t.Run("can register struct values", func(t *testing.T) {
		type Config struct {
			Name string
			Port int
		}

		scope := NewScope("test")
		cfg := Config{Name: "app", Port: 8080}
		require.NoError(t, ProvideValue(scope, cfg))

		value, err := Invoke[Config](scope)
		require.NoError(t, err)
		assert.Equal(t, "app", value.Name)
		assert.Equal(t, 8080, value.Port)
	})
}

func TestInvoke(t *testing.T) {
	t.Run("returns registered service", func(t *testing.T) {
		scope := NewScope("test")
		require.NoError(t, ProvideValue(scope, "hello"))

		value, err := Invoke[string](scope)
		require.NoError(t, err)
		assert.Equal(t, "hello", value)
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		value, err := Invoke[string](nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "scope is nil")
		assert.Equal(t, "", value)
	})

	t.Run("returns error for unregistered service", func(t *testing.T) {
		scope := NewScope("test")

		value, err := Invoke[string](scope)
		require.Error(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("can invoke different types", func(t *testing.T) {
		scope := NewScope("test")
		require.NoError(t, ProvideValue(scope, 123))
		require.NoError(t, ProvideValue(scope, "text"))
		require.NoError(t, ProvideValue(scope, true))

		intVal, err := Invoke[int](scope)
		require.NoError(t, err)
		assert.Equal(t, 123, intVal)

		strVal, err := Invoke[string](scope)
		require.NoError(t, err)
		assert.Equal(t, "text", strVal)

		boolVal, err := Invoke[bool](scope)
		require.NoError(t, err)
		assert.True(t, boolVal)
	})
}
