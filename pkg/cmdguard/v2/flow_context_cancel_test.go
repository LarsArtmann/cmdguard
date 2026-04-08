package v2

import (
	"context"
	"testing"
)

func TestBranchingFlowContext_IsLeaf(t *testing.T) {
	t.Parallel()
	root := NewBranchingFlowContext(context.Background())

	if !root.IsLeaf() {
		t.Error("expected root to be a leaf")
	}

	child, cancel := root.Branch("child")
	defer cancel()

	if root.IsLeaf() {
		t.Error("expected root to not be a leaf after adding child")
	}
	if !child.IsLeaf() {
		t.Error("expected child to be a leaf")
	}
}

func TestBranchingFlowContext_Root(t *testing.T) {
	t.Parallel()
	root := NewBranchingFlowContext(context.Background())
	child, cancel := root.Branch("child")
	defer cancel()
	grandchild, cancel2 := child.Branch("grandchild")
	defer cancel2()

	if grandchild.Root() != root {
		t.Error("expected grandchild.Root() to return root")
	}
	if child.Root() != root {
		t.Error("expected child.Root() to return root")
	}
	if root.Root() != root {
		t.Error("expected root.Root() to return root")
	}
}

func TestBranchingFlowContext_Cancel(t *testing.T) {
	t.Parallel()
	root := NewBranchingFlowContext(context.Background())
	child, cancelChild := root.Branch("child")
	grandchild, cancelGrandchild := child.Branch("grandchild")

	root.Cancel()

	select {
	case <-root.Done():
		t.Error("root context should NOT be cancelled by root.Cancel()")
	default:
	}

	select {
	case <-child.Done():
	default:
		t.Error("expected child context to be cancelled")
	}

	select {
	case <-grandchild.Done():
	default:
		t.Error("expected grandchild context to be cancelled")
	}

	cancelChild()
	cancelGrandchild()
}

func TestBranchingFlowContext_CancelChildren(t *testing.T) {
	t.Parallel()
	root := NewBranchingFlowContext(context.Background())
	child, cancelChild := root.Branch("child")
	_ = cancelChild

	root.CancelChildren()

	select {
	case <-root.Done():
		t.Error("root context should not be cancelled")
	default:
	}

	select {
	case <-child.Done():
	default:
		t.Error("child context should be cancelled")
	}
}

func TestBranchingFlowContext_CancelSiblings(t *testing.T) {
	t.Parallel()
	root := NewBranchingFlowContext(context.Background())
	sibling1, cancel1 := root.Branch("sibling1")
	sibling2, cancel2 := root.Branch("sibling2")
	defer cancel1()
	defer cancel2()

	sibling1.CancelSiblings()

	select {
	case <-sibling2.Done():
	default:
		t.Error("sibling2 should be cancelled")
	}

	select {
	case <-sibling1.Done():
		t.Error("sibling1 should not be cancelled")
	default:
	}
}
