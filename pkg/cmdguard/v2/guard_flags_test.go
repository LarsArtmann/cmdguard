package v2

import (
	"errors"
	"strings"
	"testing"
)

func TestCloneFlags(t *testing.T) {
	type testFlags struct {
		Name  string
		Count int
	}

	t.Run("clones struct", func(t *testing.T) {
		original := testFlags{Name: "test", Count: 42}
		cloned := cloneFlags(original)

		if cloned.Name != original.Name {
			t.Errorf("Name = %q, want %q", cloned.Name, original.Name)
		}

		if cloned.Count != original.Count {
			t.Errorf("Count = %d, want %d", cloned.Count, original.Count)
		}

		cloned.Name = "modified"
		if original.Name != "test" {
			t.Errorf(
				"original.Name = %q, want %q (should be unaffected by clone modification)",
				original.Name,
				"test",
			)
		}
	})

	t.Run("clones pointer to struct", func(t *testing.T) {
		original := &testFlags{Name: "test", Count: 42}
		cloned := cloneFlags(original)

		if cloned == nil {
			t.Fatal("expected non-nil cloned value")
		}

		if cloned.Name != original.Name {
			t.Errorf("Name = %q, want %q", cloned.Name, original.Name)
		}

		if cloned.Count != original.Count {
			t.Errorf("Count = %d, want %d", cloned.Count, original.Count)
		}

		if cloned == original {
			t.Error("cloned should be a different pointer from original")
		}
	})

	t.Run("returns nil for nil pointer", func(t *testing.T) {
		var original *testFlags

		cloned := cloneFlags(original)
		if cloned != nil {
			t.Errorf("cloneFlags(nil) = %v, want nil", cloned)
		}
	})

	t.Run("returns as-is for non-struct", func(t *testing.T) {
		original := "string value"

		cloned := cloneFlags(original)
		if cloned != original {
			t.Errorf("cloneFlags(%q) = %q, want %q", original, cloned, original)
		}
	})
}

func TestFlagTypeConstraint(t *testing.T) {
	type testFlags struct {
		Name  string
		Count int
	}

	t.Run("accepts NoFlags (struct{})", func(t *testing.T) {
		err := FlagTypeConstraint[NoFlags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("accepts pointer to struct", func(t *testing.T) {
		err := FlagTypeConstraint[*testFlags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("accepts empty struct", func(t *testing.T) {
		type emptyFlags struct{}

		err := FlagTypeConstraint[emptyFlags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("accepts struct with fields", func(t *testing.T) {
		err := FlagTypeConstraint[testFlags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects pointer to non-struct", func(t *testing.T) {
		err := FlagTypeConstraint[*string]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		if !strings.Contains(err.Error(), "*string") {
			t.Errorf("error should contain '*string', got %q", err.Error())
		}
	})

	t.Run("rejects int", func(t *testing.T) {
		err := FlagTypeConstraint[int]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		if !strings.Contains(err.Error(), "int") {
			t.Errorf("error should contain 'int', got %q", err.Error())
		}
	})

	t.Run("rejects string", func(t *testing.T) {
		err := FlagTypeConstraint[string]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		if !strings.Contains(err.Error(), "string") {
			t.Errorf("error should contain 'string', got %q", err.Error())
		}
	})

	t.Run("rejects slice", func(t *testing.T) {
		err := FlagTypeConstraint[[]string]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		if !strings.Contains(err.Error(), "[]string") {
			t.Errorf("error should contain '[]string', got %q", err.Error())
		}
	})

	t.Run("rejects map", func(t *testing.T) {
		err := FlagTypeConstraint[map[string]string]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		if !strings.Contains(err.Error(), "map[string]string") {
			t.Errorf("error should contain 'map[string]string', got %q", err.Error())
		}
	})
}
