package v2

import (
	"strings"
	"testing"

	"github.com/samber/do/v2"
)

func TestScopedProvider(t *testing.T) {
	t.Run("creates provider in child scope", func(t *testing.T) {
		parent := NewScope("parent")

		provider := ScopedProvider(parent, "plugin", func(i do.Injector) (string, error) {
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
	t.Run("returns error for nil parent", func(t *testing.T) {
		child, err := RegisterInScope(nil, "child")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "parent scope is nil") {
			t.Errorf("expected error to contain 'parent scope is nil', got: %v", err)
		}

		if child != nil {
			t.Errorf("expected child to be nil, got %v", child)
		}
	})

	t.Run("creates child scope with providers", func(t *testing.T) {
		parent := NewScope("parent")

		provider := func(i do.Injector) (any, error) {
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
	})

	t.Run("returns error for invalid provider type", func(t *testing.T) {
		parent := NewScope("parent")

		// Invalid provider - wrong signature
		invalidProvider := "not-a-function"

		child, err := RegisterInScope(parent, "child", invalidProvider)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "provider type=") {
			t.Errorf("expected error to contain 'provider type=', got: %v", err)
		}

		if child != nil {
			t.Errorf("expected child to be nil, got %v", child)
		}
	})

	t.Run("supports single provider", func(t *testing.T) {
		parent := NewScope("parent")

		provider := func(i do.Injector) (any, error) {
			return "single-value", nil
		}

		child, err := RegisterInScope(parent, "child", provider)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if child == nil {
			t.Fatal("expected child to not be nil")
		}
	})
}
