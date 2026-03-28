package v2

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"
)

// BranchingFlowContext provides context that branches through the command hierarchy.
// It tracks the path through command execution and allows context values to flow
// through the command tree while supporting independent branching for subcommands.
//
//nolint:containedctx // Intentional: this struct IS a context wrapper for flow control
type BranchingFlowContext struct {
	context.Context

	path       []string
	values     map[any]any
	cancels    []context.CancelFunc
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
		cancels:  nil,
		parent:   nil,
		children: nil,
	}
}

// Branch creates a new child context for a subcommand.
// The child inherits all values from the parent but has its own cancellation.
func (b *BranchingFlowContext) Branch(
	commandName string,
	opts ...FlowContextOption,
) (*BranchingFlowContext, func()) {
	branchCtx, cancel := context.WithCancel(b.Context)

	child := b.newChild(branchCtx, commandName)
	child.selfCancel = cancel
	applyOptions(child, opts)
	b.children = append(b.children, child)
	b.cancels = append(b.cancels, cancel)

	return child, cancel
}

// BranchWithTimeout creates a child context with a timeout for a subcommand.
func (b *BranchingFlowContext) BranchWithTimeout(
	commandName string,
	timeout string,
	opts ...FlowContextOption,
) (*BranchingFlowContext, func(), error) {
	d, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid timeout %q: %w", timeout, err)
	}

	branchCtx, cancel := context.WithTimeout(b.Context, d)

	child := b.newChild(branchCtx, commandName)
	child.selfCancel = cancel
	applyOptions(child, opts)
	b.children = append(b.children, child)
	b.cancels = append(b.cancels, cancel)

	return child, cancel, nil
}

// BranchWithDeadline creates a child context with a deadline for a subcommand.
func (b *BranchingFlowContext) BranchWithDeadline(
	commandName string,
	deadline string,
	opts ...FlowContextOption,
) (*BranchingFlowContext, func(), error) {
	t, err := time.Parse(time.RFC3339, deadline)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid deadline %q: %w", deadline, err)
	}

	branchCtx, cancel := context.WithDeadline(b.Context, t)

	child := b.newChild(branchCtx, commandName)
	child.selfCancel = cancel
	applyOptions(child, opts)
	b.children = append(b.children, child)
	b.cancels = append(b.cancels, cancel)

	return child, cancel, nil
}

// newChild creates a child BranchingFlowContext with the given context and command name.
func (b *BranchingFlowContext) newChild(
	ctx context.Context,
	commandName string,
) *BranchingFlowContext {
	return &BranchingFlowContext{
		Context:  ctx,
		path:     append(slices.Clone(b.path), commandName),
		values:   b.values,
		cancels:  nil,
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
func (b *BranchingFlowContext) Cancel() {
	if b.selfCancel != nil {
		b.selfCancel()
	}

	for _, cancel := range b.cancels {
		cancel()
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

// FlowContextOption configures a branching flow context.
type FlowContextOption func(*BranchingFlowContext)

// WithFlowContextValue adds a value to the new context branch.
func WithFlowContextValue(key, value any) FlowContextOption {
	return func(b *BranchingFlowContext) {
		b.values[key] = value
	}
}

// WithFlowContextValues adds multiple values to the new context branch.
func WithFlowContextValues(values map[any]any) FlowContextOption {
	return func(b *BranchingFlowContext) {
		maps.Copy(b.values, values)
	}
}

// FlowContext is an interface for context-aware DI scopes.
// Scopes implementing this interface can participate in branching flow context.
type FlowContext interface {
	// FlowContext returns the branching flow context for this scope.
	FlowContext() *BranchingFlowContext

	// Branch creates a new scope with a branched context for a subcommand.
	Branch(commandName string) (*Scope, *BranchingFlowContext, func())
}

// flowContextKey is the context key for branching flow context.
type flowContextKey struct{}

// WithBranchingFlowContext adds a branching flow context to a context.Context.
func WithBranchingFlowContext(ctx context.Context, bfc *BranchingFlowContext) context.Context {
	return context.WithValue(ctx, flowContextKey{}, bfc)
}

// GetBranchingFlowContext retrieves the branching flow context from a context.Context.
func GetBranchingFlowContext(ctx context.Context) (*BranchingFlowContext, bool) {
	if ctx == nil {
		return nil, false
	}

	val := ctx.Value(flowContextKey{})
	if val == nil {
		return nil, false
	}

	bfc, ok := val.(*BranchingFlowContext)
	return bfc, ok
}

// RequireBranchingFlowContext retrieves the branching flow context or panics.
// Use this in handlers where branching flow context is required.
func RequireBranchingFlowContext(ctx context.Context) *BranchingFlowContext {
	bfc, ok := GetBranchingFlowContext(ctx)
	if !ok {
		panic("RequireBranchingFlowContext: no branching flow context in context")
	}

	return bfc
}

// FlowContextAccessor provides convenient access to flow context values.
type FlowContextAccessor struct {
	bfc *BranchingFlowContext
}

// NewFlowContextAccessor creates an accessor for the branching flow context.
func NewFlowContextAccessor(bfc *BranchingFlowContext) *FlowContextAccessor {
	return &FlowContextAccessor{bfc: bfc}
}

// Path returns the command path.
func (a *FlowContextAccessor) Path() []string {
	return a.bfc.Path()
}

// PathString returns the command path as a string.
func (a *FlowContextAccessor) PathString() string {
	return a.bfc.PathString()
}

// Depth returns the context depth.
func (a *FlowContextAccessor) Depth() int {
	return a.bfc.Depth()
}

// Get retrieves a typed value from the flow context.
func Get[T any](ctx context.Context, key any) (T, bool) {
	bfc, ok := GetBranchingFlowContext(ctx)
	if !ok {
		var zero T
		return zero, false
	}

	val := bfc.Value(key)
	if val == nil {
		var zero T
		return zero, false
	}

	typed, ok := val.(T)
	return typed, ok
}

// MustGet retrieves a typed value from the flow context, panicking if not found.
func MustGet[T any](ctx context.Context, key any) T {
	val, ok := Get[T](ctx, key)
	if !ok {
		panic(fmt.Sprintf("MustGet: key %v not found in flow context", key))
	}

	return val
}
