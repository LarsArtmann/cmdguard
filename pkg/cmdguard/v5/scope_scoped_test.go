package v5

import (
	"testing"

	"github.com/samber/do/v2"
)

func TestScopedProvider(t *testing.T) {
	t.Parallel()
	t.Run("creates provider in child scope", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")

		provider := ScopedProvider(parent, "plugin", func(_ do.Injector) (string, error) {
			return "plugin-value", nil
		})

		value, err := provider(parent.Injector())
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if value != "plugin-value" {
			t.Errorf("expected value 'plugin-value', got %q", value)
		}
	})
}

func TestRegisterInScope(t *testing.T) {
	t.Parallel()
	t.Run("returns error for nil parent", func(t *testing.T) {
		t.Parallel()

		child, err := RegisterInScope[string](nil, "child", nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "parent scope is nil")

		if child != nil {
			t.Errorf("expected child to be nil, got %v", child)
		}
	})

	t.Run("returns error for nil provider", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")

		child, err := RegisterInScope[string](parent, "child", nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "provider is nil")

		if child != nil {
			t.Errorf("expected child to be nil, got %v", child)
		}
	})

	t.Run("registers a typed provider and resolves it", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")

		provider := func(_ do.Injector) (string, error) {
			return "service-value", nil
		}

		child, err := RegisterInScope(parent, "child", provider)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if child == nil {
			t.Fatal("expected child to not be nil")
		}

		if child.Name() != "child" {
			t.Errorf("expected name 'child', got %q", child.Name())
		}

		if child.Parent() != parent {
			t.Error("expected parent to be the same")
		}

		value, err := do.Invoke[string](child.Injector())
		if err != nil {
			t.Fatalf("expected service to resolve, got: %v", err)
		}

		if value != "service-value" {
			t.Errorf("expected value 'service-value', got %q", value)
		}
	})
}
