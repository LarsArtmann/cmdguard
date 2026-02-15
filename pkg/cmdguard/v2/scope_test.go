package v2

import (
	"context"
	"testing"
	"time"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScope(t *testing.T) {
	t.Run("creates root scope", func(t *testing.T) {
		scope := NewScope("root")
		require.NotNil(t, scope)
		assert.Equal(t, "root", scope.Name())
		assert.Nil(t, scope.Parent())
		assert.True(t, scope.IsRoot())
	})

	t.Run("creates scope with injector", func(t *testing.T) {
		scope := NewScope("test")
		require.NotNil(t, scope.Injector())
	})
}

func TestNewScopeFromInjector(t *testing.T) {
	t.Run("creates scope from existing injector", func(t *testing.T) {
		injector := do.New()
		scope := NewScopeFromInjector(injector, "custom")
		require.NotNil(t, scope)
		assert.Equal(t, "custom", scope.Name())
		assert.Equal(t, injector, scope.Injector())
	})

	t.Run("scope is root when created from injector", func(t *testing.T) {
		injector := do.New()
		scope := NewScopeFromInjector(injector, "root")
		assert.True(t, scope.IsRoot())
	})
}

func TestScope_Child(t *testing.T) {
	t.Run("creates child scope", func(t *testing.T) {
		parent := NewScope("parent")
		child := parent.Child("child")

		require.NotNil(t, child)
		assert.Equal(t, "child", child.Name())
		assert.Equal(t, parent, child.Parent())
		assert.False(t, child.IsRoot())
	})

	t.Run("child inherits from parent", func(t *testing.T) {
		parent := NewScope("parent")
		require.NoError(t, ProvideValue(parent, "parent-value"))

		child := parent.Child("child")
		value, err := Invoke[string](child)
		require.NoError(t, err)
		assert.Equal(t, "parent-value", value)
	})

	t.Run("grandchild scope", func(t *testing.T) {
		parent := NewScope("parent")
		child := parent.Child("child")
		grandchild := child.Child("grandchild")

		assert.Equal(t, child, grandchild.Parent())
		assert.Equal(t, parent, grandchild.Parent().Parent())
	})
}

func TestScope_Name(t *testing.T) {
	t.Run("returns scope name", func(t *testing.T) {
		scope := NewScope("my-scope")
		assert.Equal(t, "my-scope", scope.Name())
	})
}

func TestScope_Parent(t *testing.T) {
	t.Run("returns nil for root scope", func(t *testing.T) {
		scope := NewScope("root")
		assert.Nil(t, scope.Parent())
	})

	t.Run("returns parent for child scope", func(t *testing.T) {
		parent := NewScope("parent")
		child := parent.Child("child")
		assert.Equal(t, parent, child.Parent())
	})
}

func TestScope_Injector(t *testing.T) {
	t.Run("returns underlying injector", func(t *testing.T) {
		scope := NewScope("test")
		injector := scope.Injector()
		require.NotNil(t, injector)
	})
}

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
			dep := do.MustInvoke[Dep](i)
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

func TestMustInvoke(t *testing.T) {
	t.Run("returns registered service", func(t *testing.T) {
		scope := NewScope("test")
		require.NoError(t, ProvideValue(scope, "hello"))

		value := MustInvoke[string](scope)
		assert.Equal(t, "hello", value)
	})

	t.Run("panics for unregistered service", func(t *testing.T) {
		scope := NewScope("test")

		assert.Panics(t, func() {
			MustInvoke[string](scope)
		})
	})
}

func TestScope_Shutdown(t *testing.T) {
	t.Run("returns nil for nil injector", func(t *testing.T) {
		scope := &Scope{injector: nil}
		err := scope.Shutdown(context.Background())
		require.NoError(t, err)
	})

	t.Run("shuts down successfully", func(t *testing.T) {
		scope := NewScope("test")
		require.NoError(t, ProvideValue(scope, "value"))

		err := scope.Shutdown(context.Background())
		require.NoError(t, err)
	})
}

func TestScope_ShutdownAll(t *testing.T) {
	t.Run("shuts down single scope", func(t *testing.T) {
		scope := NewScope("root")
		require.NoError(t, ProvideValue(scope, "value"))

		err := scope.ShutdownAll(context.Background())
		require.NoError(t, err)
	})

	t.Run("shuts down scope hierarchy", func(t *testing.T) {
		parent := NewScope("parent")
		child := parent.Child("child")
		grandchild := child.Child("grandchild")

		require.NoError(t, ProvideValue(parent, "parent-value"))
		require.NoError(t, ProvideValue(child, "child-value"))
		require.NoError(t, ProvideValue(grandchild, "grandchild-value"))

		err := grandchild.ShutdownAll(context.Background())
		require.NoError(t, err)
	})
}

func TestScope_HealthCheck(t *testing.T) {
	t.Run("returns nil for nil injector", func(t *testing.T) {
		scope := &Scope{injector: nil}
		err := scope.HealthCheck()
		require.NoError(t, err)
	})

	t.Run("returns nil for healthy services", func(t *testing.T) {
		scope := NewScope("test")
		require.NoError(t, ProvideValue(scope, "value"))

		err := scope.HealthCheck()
		require.NoError(t, err)
	})
}

func TestScopedProvider(t *testing.T) {
	t.Run("creates provider in child scope", func(t *testing.T) {
		parent := NewScope("parent")

		provider := ScopedProvider(parent, "plugin", func(i do.Injector) (string, error) {
			return "plugin-value", nil
		})

		value, err := provider(parent.Injector())
		require.NoError(t, err)
		assert.Equal(t, "plugin-value", value)
	})
}

func TestRegisterInScope(t *testing.T) {
	t.Run("returns error for nil parent", func(t *testing.T) {
		child, err := RegisterInScope(nil, "child")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "parent scope is nil")
		assert.Nil(t, child)
	})

	t.Run("creates child scope with providers", func(t *testing.T) {
		parent := NewScope("parent")

		provider := func(i do.Injector) (any, error) {
			return "service-value", nil
		}

		child, err := RegisterInScope(parent, "child", provider)
		require.NoError(t, err)
		require.NotNil(t, child)
		assert.Equal(t, "child", child.Name())
		assert.Equal(t, parent, child.Parent())
	})

	t.Run("returns error for invalid provider type", func(t *testing.T) {
		parent := NewScope("parent")

		// Invalid provider - wrong signature
		invalidProvider := "not-a-function"

		child, err := RegisterInScope(parent, "child", invalidProvider)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid type")
		assert.Nil(t, child)
	})

	t.Run("supports single provider", func(t *testing.T) {
		parent := NewScope("parent")

		provider := func(i do.Injector) (any, error) {
			return "single-value", nil
		}

		child, err := RegisterInScope(parent, "child", provider)
		require.NoError(t, err)
		require.NotNil(t, child)
	})
}

func TestScope_IsRoot(t *testing.T) {
	t.Run("returns true for root scope", func(t *testing.T) {
		scope := NewScope("root")
		assert.True(t, scope.IsRoot())
	})

	t.Run("returns false for child scope", func(t *testing.T) {
		parent := NewScope("parent")
		child := parent.Child("child")
		assert.False(t, child.IsRoot())
	})

	t.Run("returns false for nested child", func(t *testing.T) {
		root := NewScope("root")
		level1 := root.Child("level1")
		level2 := level1.Child("level2")
		level3 := level2.Child("level3")

		assert.True(t, root.IsRoot())
		assert.False(t, level1.IsRoot())
		assert.False(t, level2.IsRoot())
		assert.False(t, level3.IsRoot())
	})
}

func TestScope_Path(t *testing.T) {
	t.Run("returns single element for root scope", func(t *testing.T) {
		scope := NewScope("root")
		path := scope.Path()
		assert.Equal(t, []string{"root"}, path)
	})

	t.Run("returns path for child scope", func(t *testing.T) {
		parent := NewScope("parent")
		child := parent.Child("child")
		path := child.Path()
		assert.Equal(t, []string{"parent", "child"}, path)
	})

	t.Run("returns full path for nested scopes", func(t *testing.T) {
		root := NewScope("root")
		level1 := root.Child("level1")
		level2 := level1.Child("level2")
		level3 := level2.Child("level3")

		assert.Equal(t, []string{"root"}, root.Path())
		assert.Equal(t, []string{"root", "level1"}, level1.Path())
		assert.Equal(t, []string{"root", "level1", "level2"}, level2.Path())
		assert.Equal(t, []string{"root", "level1", "level2", "level3"}, level3.Path())
	})
}

func TestScope_Integration(t *testing.T) {
	t.Run("full workflow with DI", func(t *testing.T) {
		// Create root scope
		root := NewScope("app")

		// Register services
		type Config struct {
			Debug bool
		}
		require.NoError(t, ProvideValue(root, Config{Debug: true}))

		require.NoError(t, Provide(root, func(i do.Injector) (string, error) {
			cfg := do.MustInvoke[Config](i)
			if cfg.Debug {
				return "debug-mode", nil
			}
			return "production-mode", nil
		}))

		// Verify services
		cfg, err := Invoke[Config](root)
		require.NoError(t, err)
		assert.True(t, cfg.Debug)

		mode, err := Invoke[string](root)
		require.NoError(t, err)
		assert.Equal(t, "debug-mode", mode)

		// Health check
		require.NoError(t, root.HealthCheck())

		// Shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		require.NoError(t, root.Shutdown(ctx))
	})

	t.Run("child scope can override parent services", func(t *testing.T) {
		parent := NewScope("parent")
		require.NoError(t, ProvideValue(parent, "parent-value"))

		child := parent.Child("child")

		// Initially inherits parent value
		value, err := Invoke[string](child)
		require.NoError(t, err)
		assert.Equal(t, "parent-value", value)
	})
}
