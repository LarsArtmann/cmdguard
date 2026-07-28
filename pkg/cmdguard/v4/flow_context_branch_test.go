package v3

import (
	"context"
	"testing"
	"time"
)

func TestBranchingFlowContext_Branch(t *testing.T) {
	t.Parallel()
	t.Run("creates child context", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())

		child, cancel := root.Branch("subcommand")
		defer cancel()

		if child == nil {
			t.Fatal("expected non-nil child context")
		}

		if child.Parent() != root {
			t.Error("expected child to reference parent")
		}

		children := root.Children()
		if len(children) != 1 {
			t.Errorf("expected 1 child, got %d", len(children))
		}

		if child.PathString() != "subcommand" {
			t.Errorf("expected path 'subcommand', got %q", child.PathString())
		}
	})

	t.Run("inherits values", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())
		root.SetValueLocal("key", "root-value")

		child, cancel := root.Branch("child")
		defer cancel()

		val := child.Value("key")
		if val != "root-value" {
			t.Errorf("expected 'root-value', got %v", val)
		}
	})

	t.Run("supports multiple branches", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())

		child1, cancel1 := root.Branch("cmd1")
		defer cancel1()

		child2, cancel2 := root.Branch("cmd2")
		defer cancel2()

		kids := root.Children()
		if len(kids) != 2 {
			t.Errorf("expected 2 children, got %d", len(kids))
		}

		if child1.PathString() != "cmd1" {
			t.Errorf("expected 'cmd1', got %q", child1.PathString())
		}

		if child2.PathString() != "cmd2" {
			t.Errorf("expected 'cmd2', got %q", child2.PathString())
		}
	})
}

func TestBranchingFlowContext_BranchWithDuration(t *testing.T) {
	t.Parallel()
	t.Run("valid duration", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())

		child, cancel := root.BranchWithDuration("cmd", 100*time.Millisecond)
		defer cancel()

		if child == nil {
			t.Fatal("expected non-nil child context")
		}

		if child.PathString() != "cmd" {
			t.Errorf("expected 'cmd', got %q", child.PathString())
		}
	})

	t.Run("context expires", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())

		child, cancel := root.BranchWithDuration("cmd", 1*time.Nanosecond)
		defer cancel()

		time.Sleep(10 * time.Millisecond)

		if child.Err() == nil {
			t.Error("expected context to be expired")
		}
	})
}

func TestBranchingFlowContext_BranchWithDeadlineTime(t *testing.T) {
	t.Parallel()
	t.Run("valid deadline", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())
		deadline := time.Now().Add(time.Hour)

		child, cancel := root.BranchWithDeadlineTime("cmd", deadline)
		defer cancel()

		if child == nil {
			t.Fatal("expected non-nil child context")
		}
	})

	t.Run("past deadline expires immediately", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())
		past := time.Now().Add(-time.Hour)

		child, cancel := root.BranchWithDeadlineTime("cmd", past)
		defer cancel()

		if child.Err() == nil {
			t.Error("expected context to be expired with past deadline")
		}
	})
}

func TestBranchingFlowContext_Path(t *testing.T) {
	t.Parallel()

	root := NewBranchingFlowContext(context.Background())

	if len(root.Path()) != 0 {
		t.Errorf("expected empty path, got %v", root.Path())
	}

	if root.Depth() != 0 {
		t.Errorf("expected depth 0, got %d", root.Depth())
	}

	child, cancel := root.Branch("level1")
	defer cancel()

	grandchild, cancel2 := child.Branch("level2")
	defer cancel2()

	if len(child.Path()) != 1 {
		t.Errorf("expected path length 1, got %d", len(child.Path()))
	}

	if child.Depth() != 1 {
		t.Errorf("expected depth 1, got %d", child.Depth())
	}

	if len(grandchild.Path()) != 2 {
		t.Errorf("expected path length 2, got %d", len(grandchild.Path()))
	}

	if grandchild.Depth() != 2 {
		t.Errorf("expected depth 2, got %d", grandchild.Depth())
	}

	if grandchild.PathString() != "level1.level2" {
		t.Errorf("expected 'level1.level2', got %q", grandchild.PathString())
	}
}
