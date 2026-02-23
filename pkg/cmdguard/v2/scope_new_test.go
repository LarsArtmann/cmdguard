package v2

import (
	"testing"

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
