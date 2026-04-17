package v2

import (
	"context"
	"testing"
)

func TestBranchingFlowContext_SetValue(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())

	child, cancel := root.Branch("child")
	defer cancel()

	root.SetValue("propagate", "yes")

	assertFlowValue(t, root, "propagate", "yes", "root should have value")

	assertFlowValue(t, child, "propagate", "yes", "child should inherit value from root")
}

func TestBranchingFlowContext_SetValueLocal(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())

	child, cancel := root.Branch("child")
	defer cancel()

	root.SetValueLocal("local", "only-root")

	assertFlowValue(t, root, "local", "only-root", "root should have local value")

	assertFlowValue(t, child, "local", "only-root", "child should see parent value via GetValue fallback")
}

func TestBranchingFlowContext_ChildValueIsolation(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())
	root.SetValueLocal("shared", "root-val")

	child, cancel := root.Branch("child")
	defer cancel()

	child.SetValueLocal("shared", "child-val")

	assertFlowValue(t, root, "shared", "root-val", "root value should not be affected by child's SetValueLocal")

	assertFlowValue(t, child, "shared", "child-val", "child should see its own local value")
}

func TestBranchingFlowContext_GetValue(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())

	child, cancel := root.Branch("child")
	t.Cleanup(cancel)

	root.SetValue("inherited", "from-root")
	child.SetValueLocal("local", "child-only")

	testCases := []struct {
		name     string
		getter   func() (any, bool)
		expected any
	}{
		{
			name: "inherited value",
			getter: func() (any, bool) {
				return child.GetValue("inherited")
			},
			expected: "from-root",
		},
		{
			name:     "local value",
			getter:   func() (any, bool) { return child.GetValue("local") },
			expected: "child-only",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			val, ok := tc.getter()
			if !ok {
				t.Error("expected to find value")
			}

			if val != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, val)
			}
		})
	}

	t.Run("missing value", func(t *testing.T) {
		t.Parallel()

		_, ok := root.GetValue("nonexistent")
		if ok {
			t.Error("expected not to find nonexistent value")
		}
	})
}
