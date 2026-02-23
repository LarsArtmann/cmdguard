package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
