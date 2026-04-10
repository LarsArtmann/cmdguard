package v2

import (
	"context"
	"testing"
)

func TestWithFlowContextValue(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())

	child, cancel := root.Branch("cmd", WithFlowContextValue("key", "value"))
	defer cancel()

	if child.Value("key") != "value" {
		t.Error("expected child to have value set via option")
	}
}

func TestWithFlowContextValues(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())
	values := map[any]any{
		"key1": "value1",
		"key2": "value2",
	}

	child, cancel := root.Branch("cmd", WithFlowContextValues(values))
	defer cancel()

	if child.Value("key1") != "value1" {
		t.Error("expected key1 to be set")
	}

	if child.Value("key2") != "value2" {
		t.Error("expected key2 to be set")
	}
}

func TestWithBranchingFlowContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	bfc := NewBranchingFlowContext(ctx)

	wrapped := WithBranchingFlowContext(ctx, bfc)

	retrieved, ok := GetBranchingFlowContext(wrapped)
	if !ok {
		t.Fatal("expected to retrieve branching flow context")
	}

	if retrieved != bfc {
		t.Error("expected retrieved context to match original")
	}
}

func TestGetBranchingFlowContext_NilContext(t *testing.T) {
	t.Parallel()

	_, ok := GetBranchingFlowContext(context.TODO())
	if ok {
		t.Error("expected ok to be false for context.TODO()")
	}
}

func TestGetBranchingFlowContext_NotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	_, ok := GetBranchingFlowContext(ctx)
	if ok {
		t.Error("expected ok to be false when no branching flow context")
	}
}

func TestRequireBranchingFlowContext(t *testing.T) {
	t.Parallel()
	t.Run("found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		bfc := NewBranchingFlowContext(ctx)
		wrapped := WithBranchingFlowContext(ctx, bfc)

		retrieved := RequireBranchingFlowContext(wrapped)
		if retrieved != bfc {
			t.Error("expected retrieved context to match original")
		}
	})

	t.Run("not found panics", func(t *testing.T) {
		t.Parallel()

		assertPanics(t, func() {
			ctx := context.Background()
			RequireBranchingFlowContext(ctx)
		})
	})
}
