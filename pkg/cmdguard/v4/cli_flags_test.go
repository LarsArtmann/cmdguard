package v4

import (
	"errors"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

func TestCloneFlags(t *testing.T) {
	t.Parallel()

	type testFlags struct {
		Name  string
		Count int
	}

	t.Run("clones struct", func(t *testing.T) {
		t.Parallel()

		original := testFlags{Name: "test", Count: 42}
		cloned := cloneFlags(original)

		testutil.AssertFieldEqString(t, cloned.Name, original.Name, "Name")
		testutil.AssertFieldEq(t, cloned.Count, original.Count, "Count")

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
		t.Parallel()

		original := &testFlags{Name: "test", Count: 42}
		cloned := cloneFlags(original)

		if cloned == nil {
			t.Fatal("expected non-nil cloned value")
		}

		testutil.AssertFieldEqString(t, cloned.Name, original.Name, "Name")
		testutil.AssertFieldEq(t, cloned.Count, original.Count, "Count")

		if cloned == original {
			t.Error("cloned should be a different pointer from original")
		}
	})

	t.Run("returns nil for nil pointer", func(t *testing.T) {
		t.Parallel()

		var original *testFlags

		cloned := cloneFlags(original)
		if cloned != nil {
			t.Errorf("cloneFlags(nil) = %v, want nil", cloned)
		}
	})

	t.Run("returns as-is for non-struct", func(t *testing.T) {
		t.Parallel()

		original := "string value"

		cloned := cloneFlags(original)
		if cloned != original {
			t.Errorf("cloneFlags(%q) = %q, want %q", original, cloned, original)
		}
	})

	t.Run("deep copies nested slices in struct", func(t *testing.T) {
		t.Parallel()

		type nestedFlags struct {
			Tags []string
		}

		original := nestedFlags{Tags: []string{"a", "b", "c"}}
		cloned := cloneFlags(original)

		testutil.AssertFieldEq(t, len(cloned.Tags), 3, "len(Tags)")

		cloned.Tags[0] = "modified"
		if original.Tags[0] != "a" {
			t.Errorf(
				"original.Tags[0] = %q, want %q (deep copy should isolate slices)",
				original.Tags[0],
				"a",
			)
		}
	})

	t.Run("deep copies nested pointer in struct", func(t *testing.T) {
		t.Parallel()

		type ptrFlags struct {
			Label *string
		}

		label := "original"
		original := ptrFlags{Label: &label}
		cloned := cloneFlags(original)

		if cloned.Label == nil {
			t.Fatal("expected non-nil Label")
		}

		*cloned.Label = "modified"
		if *original.Label != "original" {
			t.Errorf(
				"original.Label = %q, want %q (deep copy should isolate pointers)",
				*original.Label,
				"original",
			)
		}
	})

	t.Run("deep copies nested slices in pointer to struct", func(t *testing.T) {
		t.Parallel()

		type nestedPtr struct {
			Items []string
		}

		original := &nestedPtr{Items: []string{"x", "y"}}
		cloned := cloneFlags(original)

		if cloned == original {
			t.Error("cloned should be a different pointer")
		}

		cloned.Items[0] = "modified"
		if original.Items[0] != "x" {
			t.Errorf(
				"original.Items[0] = %q, want %q (deep copy should isolate nested slices)",
				original.Items[0],
				"x",
			)
		}
	})
}

func TestFlagTypeConstraint(t *testing.T) {
	t.Parallel()

	type testFlags struct {
		Name  string
		Count int
	}

	t.Run("accepts NoFlags (struct{})", func(t *testing.T) {
		t.Parallel()

		err := FlagTypeConstraint[NoFlags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("accepts pointer to struct", func(t *testing.T) {
		t.Parallel()

		err := FlagTypeConstraint[*testFlags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("accepts empty struct", func(t *testing.T) {
		t.Parallel()

		type emptyFlags struct{}

		err := FlagTypeConstraint[emptyFlags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("accepts struct with fields", func(t *testing.T) {
		t.Parallel()

		err := FlagTypeConstraint[testFlags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects pointer to non-struct", func(t *testing.T) {
		t.Parallel()

		err := FlagTypeConstraint[*string]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		assertErrorContains(t, err, "*string")
	})

	t.Run("rejects int", func(t *testing.T) {
		t.Parallel()

		err := FlagTypeConstraint[int]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		assertErrorContains(t, err, "int")
	})

	t.Run("rejects string", func(t *testing.T) {
		t.Parallel()

		err := FlagTypeConstraint[string]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		assertErrorContains(t, err, "string")
	})

	t.Run("rejects slice", func(t *testing.T) {
		t.Parallel()

		err := FlagTypeConstraint[[]string]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		assertErrorContains(t, err, "[]string")
	})

	t.Run("rejects map", func(t *testing.T) {
		t.Parallel()

		err := FlagTypeConstraint[map[string]string]()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		assertErrorContains(t, err, "map[string]string")
	})
}
