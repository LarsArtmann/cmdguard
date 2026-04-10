package v2

import (
	"testing"
)

func TestScope_Child(t *testing.T) {
	t.Parallel()
	t.Run("creates child scope", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")
		child := parent.Child("child")

		if child == nil {
			t.Fatal("expected child to not be nil")
		}

		if child.Name() != "child" {
			t.Errorf("expected name to be 'child', got %q", child.Name())
		}

		if child.Parent() != parent {
			t.Error("expected parent to be the same")
		}

		if child.IsRoot() {
			t.Error("expected IsRoot to be false")
		}
	})

	t.Run("child inherits from parent", func(t *testing.T) {
		t.Parallel()
		assertChildInheritsParent(t)
	})

	t.Run("grandchild scope", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")
		child := parent.Child("child")
		grandchild := child.Child("grandchild")

		if grandchild.Parent() != child {
			t.Error("expected grandchild parent to be child")
		}

		if grandchild.Parent().Parent() != parent {
			t.Error("expected grandchild grandparent to be parent")
		}
	})
}

func TestScope_Name(t *testing.T) {
	t.Parallel()
	t.Run("returns scope name", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("my-scope")
		if scope.Name() != "my-scope" {
			t.Errorf("expected name to be 'my-scope', got %q", scope.Name())
		}
	})
}

func TestScope_Parent(t *testing.T) {
	t.Parallel()
	t.Run("returns nil for root scope", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("root")
		if scope.Parent() != nil {
			t.Errorf("expected parent to be nil, got %v", scope.Parent())
		}
	})

	t.Run("returns parent for child scope", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")

		child := parent.Child("child")
		if child.Parent() != parent {
			t.Error("expected parent to be the same")
		}
	})
}

func TestScope_Injector(t *testing.T) {
	t.Parallel()
	t.Run("returns underlying injector", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("test")

		injector := scope.Injector()
		if injector == nil {
			t.Fatal("expected injector to not be nil")
		}
	})
}

func TestScope_IsRoot(t *testing.T) {
	t.Parallel()
	t.Run("returns true for root scope", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("root")
		if !scope.IsRoot() {
			t.Error("expected IsRoot to be true")
		}
	})

	t.Run("returns false for child scope", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")

		child := parent.Child("child")
		if child.IsRoot() {
			t.Error("expected IsRoot to be false")
		}
	})

	t.Run("returns false for nested child", func(t *testing.T) {
		t.Parallel()

		root := NewScope("root")
		level1 := root.Child("level1")
		level2 := level1.Child("level2")
		level3 := level2.Child("level3")

		if !root.IsRoot() {
			t.Error("expected root IsRoot to be true")
		}

		if level1.IsRoot() {
			t.Error("expected level1 IsRoot to be false")
		}

		if level2.IsRoot() {
			t.Error("expected level2 IsRoot to be false")
		}

		if level3.IsRoot() {
			t.Error("expected level3 IsRoot to be false")
		}
	})
}

func TestScope_Path(t *testing.T) {
	t.Parallel()
	t.Run("returns single element for root scope", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("root")
		path := scope.Path()

		expected := []string{"root"}
		if !slicesEqual(path, expected) {
			t.Errorf("expected path %v, got %v", expected, path)
		}
	})

	t.Run("returns path for child scope", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")
		child := parent.Child("child")
		path := child.Path()

		expected := []string{"parent", "child"}
		if !slicesEqual(path, expected) {
			t.Errorf("expected path %v, got %v", expected, path)
		}
	})

	t.Run("returns full path for nested scopes", func(t *testing.T) {
		t.Parallel()

		root := NewScope("root")
		level1 := root.Child("level1")
		level2 := level1.Child("level2")
		level3 := level2.Child("level3")

		if !slicesEqual(root.Path(), []string{"root"}) {
			t.Errorf("expected root path [root], got %v", root.Path())
		}

		if !slicesEqual(level1.Path(), []string{"root", "level1"}) {
			t.Errorf("expected level1 path [root level1], got %v", level1.Path())
		}

		if !slicesEqual(level2.Path(), []string{"root", "level1", "level2"}) {
			t.Errorf("expected level2 path [root level1 level2], got %v", level2.Path())
		}

		if !slicesEqual(level3.Path(), []string{"root", "level1", "level2", "level3"}) {
			t.Errorf("expected level3 path [root level1 level2 level3], got %v", level3.Path())
		}
	})
}

func assertChildInheritsParent(t *testing.T) {
	t.Helper()

	parent := NewScope("parent")
	if err := ProvideValue(parent, "parent-value"); err != nil {
		t.Fatalf("expected no error providing value, got: %v", err)
	}

	child := parent.Child("child")

	value, err := Invoke[string](child)
	if err != nil {
		t.Fatalf("expected no error invoking, got: %v", err)
	}

	if value != "parent-value" {
		t.Errorf("expected value 'parent-value', got %q", value)
	}
}
