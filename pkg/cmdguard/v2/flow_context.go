package v2

import (
	"context"
	"maps"
	"slices"
	"strings"
	"time"
)

// BranchingFlowContext provides context that branches through the command hierarchy.
// It tracks the path through command execution and allows context values to flow
// through the command tree while supporting independent branching for subcommands.
//
// BranchingFlowContext is not goroutine-safe. Concurrent calls to SetValue/Branch
// require external synchronization. In typical CLI usage this is not needed since
// command execution is sequential.
//
//nolint:containedctx // Intentional: this struct IS a context wrapper for flow control
type BranchingFlowContext struct {
	context.Context

	path       []string
	values     map[any]any
	selfCancel context.CancelFunc
	parent     *BranchingFlowContext
	children   []*BranchingFlowContext
}

// NewBranchingFlowContext creates a new root branching flow context.
//
//nolint:contextcheck // Intentional: factory function wraps provided context
func NewBranchingFlowContext(ctx context.Context) *BranchingFlowContext {
	if ctx == nil {
		ctx = context.Background()
	}

	return &BranchingFlowContext{
		Context:  ctx,
		path:     []string{},
		values:   make(map[any]any),
		parent:   nil,
		children: nil,
	}
}

// branchWithCtx creates a child context, registers it, and applies options.
func (b *BranchingFlowContext) branchWithCtx(
	branchCtx context.Context,
	cancel context.CancelFunc,
	commandName string,
	opts []FlowContextOption,
) *BranchingFlowContext {
	child := b.newChild(branchCtx, commandName)
	child.selfCancel = cancel
	applyOptions(child, opts)
	b.children = append(b.children, child)

	return child
}

// Branch creates a new child context for a subcommand.
// The child inherits all values from the parent but has its own cancellation.
func (b *BranchingFlowContext) Branch(
	commandName string,
	opts ...FlowContextOption,
) (*BranchingFlowContext, func()) {
	branchCtx, cancel := context.WithCancel(b.Context)

	return b.branchWithCtx(branchCtx, cancel, commandName, opts), cancel
}

// BranchWithDuration creates a child context with a timeout for a subcommand.
// Prefer this over BranchWithTimeout — it accepts time.Duration directly
// instead of parsing a string at runtime.
func (b *BranchingFlowContext) BranchWithDuration(
	commandName string,
	d time.Duration,
	opts ...FlowContextOption,
) (*BranchingFlowContext, func()) {
	branchCtx, cancel := context.WithTimeout(b.Context, d)

	return b.branchWithCtx(branchCtx, cancel, commandName, opts), cancel
}

// BranchWithDeadlineTime creates a child context with a deadline for a subcommand.
// Prefer this over BranchWithDeadline — it accepts time.Time directly
// instead of parsing a string at runtime.
func (b *BranchingFlowContext) BranchWithDeadlineTime(
	commandName string,
	deadline time.Time,
	opts ...FlowContextOption,
) (*BranchingFlowContext, func()) {
	branchCtx, cancel := context.WithDeadline(b.Context, deadline)

	return b.branchWithCtx(branchCtx, cancel, commandName, opts), cancel
}

// newChild creates a child BranchingFlowContext with the given context and command name.
func (b *BranchingFlowContext) newChild(
	ctx context.Context,
	commandName string,
) *BranchingFlowContext {
	return &BranchingFlowContext{
		Context:  ctx,
		path:     append(slices.Clone(b.path), commandName),
		values:   maps.Clone(b.values),
		parent:   b,
		children: nil,
	}
}

// applyOptions applies FlowContextOptions to a BranchingFlowContext.
func applyOptions(child *BranchingFlowContext, opts []FlowContextOption) {
	for _, opt := range opts {
		opt(child)
	}
}

// Path returns the command path from root to this context.
func (b *BranchingFlowContext) Path() []string {
	return b.path
}

// PathString returns the command path as a dot-separated string.
func (b *BranchingFlowContext) PathString() string {
	return strings.Join(b.path, ".")
}

// Depth returns the depth of this context in the tree (root = 0).
func (b *BranchingFlowContext) Depth() int {
	return len(b.path)
}

// IsRoot returns true if this is the root context.
func (b *BranchingFlowContext) IsRoot() bool {
	return b.parent == nil
}

// IsLeaf returns true if this context has no children.
func (b *BranchingFlowContext) IsLeaf() bool {
	return len(b.children) == 0
}

// Parent returns the parent context, or nil if this is the root.
func (b *BranchingFlowContext) Parent() *BranchingFlowContext {
	return b.parent
}

// Children returns the child contexts.
func (b *BranchingFlowContext) Children() []*BranchingFlowContext {
	return b.children
}

// Root returns the root context of this tree.
func (b *BranchingFlowContext) Root() *BranchingFlowContext {
	current := b
	for current.parent != nil {
		current = current.parent
	}

	return current
}

// SetValue sets a value in this context and all descendants.
// Use this to propagate values down the command tree.
func (b *BranchingFlowContext) SetValue(key, value any) {
	b.values[key] = value

	for _, child := range b.children {
		child.SetValue(key, value)
	}
}

// SetValueLocal sets a value only in this context (not propagated to children).
func (b *BranchingFlowContext) SetValueLocal(key, value any) {
	b.values[key] = value
}

// GetValue retrieves a value from this context or any ancestor.
func (b *BranchingFlowContext) GetValue(key any) (any, bool) {
	if v, ok := b.values[key]; ok {
		return v, true
	}

	if b.parent != nil {
		return b.parent.GetValue(key)
	}

	return nil, false
}

// Value retrieves a value, returning nil if not found.
func (b *BranchingFlowContext) Value(key any) any {
	v, _ := b.GetValue(key)

	return v
}

// Cancel cancels this context and all its children.
// Each node cancels its own context via selfCancel, then recursively cancels children.
func (b *BranchingFlowContext) Cancel() {
	if b.selfCancel != nil {
		b.selfCancel()
	}

	for _, child := range b.children {
		child.Cancel()
	}
}

// CancelChildren cancels all child contexts but not this one.
func (b *BranchingFlowContext) CancelChildren() {
	for _, child := range b.children {
		child.Cancel()
	}
}

// CancelSiblings cancels all sibling contexts but not this one.
func (b *BranchingFlowContext) CancelSiblings() {
	if b.parent == nil {
		return
	}

	for _, sibling := range b.parent.children {
		if sibling != b {
			sibling.Cancel()
		}
	}
}
