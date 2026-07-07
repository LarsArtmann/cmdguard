package v3

import (
	"context"
	"testing"

	"github.com/samber/do/v2"
)

func TestScope_Shutdown(t *testing.T) {
	t.Parallel()
	t.Run("returns nil for nil injector", func(t *testing.T) {
		t.Parallel()

		scope := &Scope{injector: nil}

		err := scope.Shutdown(t.Context())
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("shuts down successfully", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("test")
		mustProvideValue(t, scope, "value")

		err := scope.Shutdown(t.Context())
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestScope_ShutdownAll(t *testing.T) {
	t.Parallel()
	t.Run("shuts down single scope", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("root")
		mustProvideValue(t, scope, "value")

		err := scope.ShutdownAll(t.Context())
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("shuts down scope hierarchy", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")
		child := parent.Child("child")
		grandchild := child.Child("grandchild")

		mustProvideValue(t, parent, "parent-value")

		mustProvideValue(t, child, "child-value")

		mustProvideValue(t, grandchild, "grandchild-value")

		err := grandchild.ShutdownAll(t.Context())
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestScope_HealthCheck(t *testing.T) {
	t.Parallel()
	t.Run("returns nil for nil injector", func(t *testing.T) {
		t.Parallel()

		scope := &Scope{injector: nil}

		err := scope.HealthCheck()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("returns nil for healthy services", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("test")
		mustProvideValue(t, scope, "value")

		err := scope.HealthCheck()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

type failingShutdownService struct{}

func (f *failingShutdownService) Shutdown(_ context.Context) error {
	return context.DeadlineExceeded
}

func TestScope_ShutdownAll_WithError(t *testing.T) {
	t.Parallel()
	t.Run("accumulates errors from shutdown failures", func(t *testing.T) {
		t.Parallel()

		scope := NewScope("test")

		do.Provide(scope.Injector(), func(_ do.Injector) (*failingShutdownService, error) {
			return &failingShutdownService{}, nil
		})

		_, err := Invoke[*failingShutdownService](scope)
		if err != nil {
			t.Fatalf("expected no error invoking, got: %v", err)
		}

		err = scope.ShutdownAll(t.Context())
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "scope \"test\"")
	})
}
