package v2

import (
	"context"
	"testing"
)

func TestNewBranchingFlowContext(t *testing.T) {
	t.Parallel()
	t.Run("with context.TODO", func(t *testing.T) {
		t.Parallel()
		bfc := NewBranchingFlowContext(context.TODO())
		if bfc == nil {
			t.Fatal("expected non-nil BranchingFlowContext")
		}
		if bfc.Context == nil {
			t.Error("expected embedded context to be non-nil")
		}
		if !bfc.IsRoot() {
			t.Error("expected IsRoot to be true")
		}
		if bfc.PathString() != "" {
			t.Errorf("expected empty path, got %q", bfc.PathString())
		}
	})

	t.Run("with valid context", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		bfc := NewBranchingFlowContext(ctx)
		if bfc.Context != ctx {
			t.Error("expected embedded context to match input")
		}
	})
}
