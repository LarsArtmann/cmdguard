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

	if root.Value("propagate") != "yes" {
		t.Error("root should have value")
	}

	if child.Value("propagate") != "yes" {
		t.Error("child should inherit value from root")
	}
}

func TestBranchingFlowContext_SetValueLocal(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())

	child, cancel := root.Branch("child")
	defer cancel()

	root.SetValueLocal("local", "only-root")

	if root.Value("local") != "only-root" {
		t.Error("root should have local value")
	}

	if child.Value("local") != "only-root" {
		t.Error("child should see parent value via GetValue fallback")
	}
}

func TestBranchingFlowContext_ChildValueIsolation(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())
	root.SetValueLocal("shared", "root-val")

	child, cancel := root.Branch("child")
	defer cancel()

	child.SetValueLocal("shared", "child-val")

	if root.Value("shared") != "root-val" {
		t.Error("root value should not be affected by child's SetValueLocal")
	}

	if child.Value("shared") != "child-val" {
		t.Error("child should see its own local value")
	}
}

func TestBranchingFlowContext_GetValue(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())

	child, cancel := root.Branch("child")
	defer cancel()

	root.SetValue("inherited", "from-root")
	child.SetValueLocal("local", "child-only")

	testCases := []struct {
		name     string
		getter   func() (any, bool)
		expected any
	}{
		{
			name: "inherited value",
			getter: func() (v any, ok bool) {
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
