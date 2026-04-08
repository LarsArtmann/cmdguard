package v2

import (
	"testing"

	"github.com/samber/do/v2"
)

func TestNewScope(t *testing.T) {
	t.Parallel()
	t.Run("creates root scope", func(t *testing.T) {
		t.Parallel()
		scope := NewScope("root")
		if scope == nil {
			t.Fatal("expected scope to not be nil")
		}

		if scope.Name() != "root" {
			t.Errorf("expected name to be 'root', got %q", scope.Name())
		}

		if scope.Parent() != nil {
			t.Errorf("expected parent to be nil, got %v", scope.Parent())
		}

		if !scope.IsRoot() {
			t.Error("expected IsRoot to be true")
		}
	})

	t.Run("creates scope with injector", func(t *testing.T) {
		t.Parallel()
		scope := NewScope("test")
		if scope.Injector() == nil {
			t.Fatal("expected injector to not be nil")
		}
	})
}

func TestNewScopeFromInjector(t *testing.T) {
	t.Parallel()
	t.Run("creates scope from existing injector", func(t *testing.T) {
		t.Parallel()
		injector := do.New()

		scope := NewScopeFromInjector(injector, "custom")
		if scope == nil {
			t.Fatal("expected scope to not be nil")
		}

		if scope.Name() != "custom" {
			t.Errorf("expected name to be 'custom', got %q", scope.Name())
		}

		if scope.Injector() != injector {
			t.Error("expected injector to be the same")
		}
	})

	t.Run("scope is root when created from injector", func(t *testing.T) {
		t.Parallel()
		injector := do.New()

		scope := NewScopeFromInjector(injector, "root")
		if !scope.IsRoot() {
			t.Error("expected IsRoot to be true")
		}
	})
}
