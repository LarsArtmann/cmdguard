package v2

import (
	"context"
	"strings"
	"testing"

	"github.com/samber/do/v2"
)

func TestScope_Shutdown(t *testing.T) {
	t.Run("returns nil for nil injector", func(t *testing.T) {
		scope := &Scope{injector: nil}

		err := scope.Shutdown(context.Background())
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("shuts down successfully", func(t *testing.T) {
		scope := NewScope("test")
		if err := ProvideValue(scope, "value"); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		err := scope.Shutdown(context.Background())
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestScope_ShutdownAll(t *testing.T) {
	t.Run("shuts down single scope", func(t *testing.T) {
		scope := NewScope("root")
		if err := ProvideValue(scope, "value"); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		err := scope.ShutdownAll(context.Background())
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("shuts down scope hierarchy", func(t *testing.T) {
		parent := NewScope("parent")
		child := parent.Child("child")
		grandchild := child.Child("grandchild")

		if err := ProvideValue(parent, "parent-value"); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		if err := ProvideValue(child, "child-value"); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		if err := ProvideValue(grandchild, "grandchild-value"); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		err := grandchild.ShutdownAll(context.Background())
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestScope_HealthCheck(t *testing.T) {
	t.Run("returns nil for nil injector", func(t *testing.T) {
		scope := &Scope{injector: nil}

		err := scope.HealthCheck()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("returns nil for healthy services", func(t *testing.T) {
		scope := NewScope("test")
		if err := ProvideValue(scope, "value"); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

		err := scope.HealthCheck()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

type failingShutdownService struct{}

func (f *failingShutdownService) Shutdown(ctx context.Context) error {
	return context.DeadlineExceeded
}

func TestScope_ShutdownAll_WithError(t *testing.T) {
	t.Run("accumulates errors from shutdown failures", func(t *testing.T) {
		scope := NewScope("test")

		do.Provide(scope.Injector(), func(i do.Injector) (*failingShutdownService, error) {
			return &failingShutdownService{}, nil
		})

		_, err := Invoke[*failingShutdownService](scope)
		if err != nil {
			t.Fatalf("expected no error invoking, got: %v", err)
		}

		err = scope.ShutdownAll(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "scope \"test\"") {
			t.Errorf("expected error to contain scope name, got: %v", err)
		}
	})
}
