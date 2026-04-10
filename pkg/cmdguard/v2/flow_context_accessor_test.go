package v2

import (
	"context"
	"testing"
)

func TestFlowContextAccessor(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())
	child, cancel := root.Branch("level1")
	grandchild, _ := child.Branch("level2")

	cancel()

	accessor := NewFlowContextAccessor(grandchild)

	if len(accessor.Path()) != 2 {
		t.Errorf("expected path length 2, got %d", len(accessor.Path()))
	}

	if accessor.PathString() != "level1.level2" {
		t.Errorf("expected 'level1.level2', got %q", accessor.PathString())
	}

	if accessor.Depth() != 2 {
		t.Errorf("expected depth 2, got %d", accessor.Depth())
	}
}

func TestGet_Typed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	bfc := NewBranchingFlowContext(ctx)
	bfc.SetValue("string-key", "string-value")
	bfc.SetValue("int-key", 42)

	wrapped := WithBranchingFlowContext(ctx, bfc)

	t.Run("string value", func(t *testing.T) {
		t.Parallel()

		val, ok := Get[string](wrapped, "string-key")
		if !ok {
			t.Fatal("expected to find string value")
		}

		if val != "string-value" {
			t.Errorf("expected 'string-value', got %q", val)
		}
	})

	t.Run("int value", func(t *testing.T) {
		t.Parallel()

		val, ok := Get[int](wrapped, "int-key")
		if !ok {
			t.Fatal("expected to find int value")
		}

		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("missing value", func(t *testing.T) {
		t.Parallel()

		_, ok := Get[string](wrapped, "nonexistent")
		if ok {
			t.Error("expected not to find nonexistent value")
		}
	})
}

func TestMustGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	bfc := NewBranchingFlowContext(ctx)
	bfc.SetValue("key", "value")

	wrapped := WithBranchingFlowContext(ctx, bfc)

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		val := MustGet[string](wrapped, "key")
		if val != "value" {
			t.Errorf("expected 'value', got %q", val)
		}
	})

	t.Run("missing panics", func(t *testing.T) {
		t.Parallel()

		assertPanics(t, func() {
			MustGet[string](wrapped, "nonexistent")
		})
	})
}
