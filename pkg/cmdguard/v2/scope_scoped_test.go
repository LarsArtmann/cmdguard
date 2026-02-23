package v2

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
