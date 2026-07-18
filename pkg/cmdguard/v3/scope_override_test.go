package v3_test

import (
	"context"
	"testing"

	"github.com/samber/do/v2"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

type realService struct {
	Name string
}

// provideRealService registers a realService factory that constructs an instance
// with the given name. Used by override and clone-scope tests.
func provideRealService(t *testing.T, scope *v3.Scope, name string) {
	t.Helper()

	err := v3.Provide(scope, func(i do.Injector) (*realService, error) {
		return &realService{Name: name}, nil
	})
	if err != nil {
		t.Fatalf("Provide failed: %v", err)
	}
}

func TestOverride(t *testing.T) {
	t.Parallel()

	t.Run("replaces service provider with mock", func(t *testing.T) {
		t.Parallel()

		scope := v3.NewScope("test")

		provideRealService(t, scope, "real")

		err := v3.Override(scope, func(i do.Injector) (*realService, error) {
			return &realService{Name: "mock"}, nil
		})
		if err != nil {
			t.Fatalf("Override failed: %v", err)
		}

		svc, err := v3.Invoke[*realService](scope)
		if err != nil {
			t.Fatalf("Invoke failed: %v", err)
		}

		if svc.Name != "mock" {
			t.Errorf("Name = %q, want %q", svc.Name, "mock")
		}
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		t.Parallel()

		err := v3.Override(nil, func(i do.Injector) (*realService, error) {
			return &realService{}, nil
		})
		if err == nil {
			t.Fatal("expected error for nil scope")
		}
	})
}

func TestOverrideValue(t *testing.T) {
	t.Parallel()

	t.Run("replaces pre-constructed value", func(t *testing.T) {
		t.Parallel()

		scope := v3.NewScope("test")

		err := v3.ProvideValue(scope, &realService{Name: "original"})
		if err != nil {
			t.Fatalf("ProvideValue failed: %v", err)
		}

		err = v3.OverrideValue(scope, &realService{Name: "overridden"})
		if err != nil {
			t.Fatalf("OverrideValue failed: %v", err)
		}

		svc, err := v3.Invoke[*realService](scope)
		if err != nil {
			t.Fatalf("Invoke failed: %v", err)
		}

		if svc.Name != "overridden" {
			t.Errorf("Name = %q, want %q", svc.Name, "overridden")
		}
	})

	t.Run("returns error for nil scope", func(t *testing.T) {
		t.Parallel()

		err := v3.OverrideValue(nil, &realService{})
		if err == nil {
			t.Fatal("expected error for nil scope")
		}
	})
}

func TestCloneScope(t *testing.T) {
	t.Parallel()

	t.Run("creates independent scope with same registrations", func(t *testing.T) {
		t.Parallel()

		scope := v3.NewScope("test")

		err := v3.ProvideValue(scope, &realService{Name: "original"})
		if err != nil {
			t.Fatalf("ProvideValue failed: %v", err)
		}

		cloned := v3.CloneScope(scope)

		svc, err := v3.Invoke[*realService](cloned)
		if err != nil {
			t.Fatalf("Invoke on cloned failed: %v", err)
		}

		if svc.Name != "original" {
			t.Errorf("cloned Name = %q, want %q", svc.Name, "original")
		}
	})

	t.Run("override in clone does not affect original", func(t *testing.T) {
		t.Parallel()

		scope := v3.NewScope("test")

		err := v3.ProvideValue(scope, &realService{Name: "original"})
		if err != nil {
			t.Fatalf("ProvideValue failed: %v", err)
		}

		cloned := v3.CloneScope(scope)

		err = v3.OverrideValue(cloned, &realService{Name: "mocked"})
		if err != nil {
			t.Fatalf("OverrideValue failed: %v", err)
		}

		clonedSvc, err := v3.Invoke[*realService](cloned)
		if err != nil {
			t.Fatalf("Invoke on cloned failed: %v", err)
		}

		originalSvc, err := v3.Invoke[*realService](scope)
		if err != nil {
			t.Fatalf("Invoke on original failed: %v", err)
		}

		if clonedSvc.Name != "mocked" {
			t.Errorf("cloned Name = %q, want %q", clonedSvc.Name, "mocked")
		}

		if originalSvc.Name != "original" {
			t.Errorf("original Name = %q, want %q", originalSvc.Name, "original")
		}
	})

	t.Run("cloned scope invokes original providers", func(t *testing.T) {
		t.Parallel()

		scope := v3.NewScope("test")

		provideRealService(t, scope, "lazy-original")

		cloned := v3.CloneScope(scope)

		svc, err := v3.Invoke[*realService](cloned)
		if err != nil {
			t.Fatalf("Invoke on cloned failed: %v", err)
		}

		if svc.Name != "lazy-original" {
			t.Errorf("Name = %q, want %q", svc.Name, "lazy-original")
		}
	})

	t.Run("cloned scope supports shutdown", func(t *testing.T) {
		t.Parallel()

		scope := v3.NewScope("test")
		cloned := v3.CloneScope(scope)

		err := cloned.Shutdown(context.Background())
		if err != nil {
			t.Errorf("Shutdown failed: %v", err)
		}
	})
}
