package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

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
