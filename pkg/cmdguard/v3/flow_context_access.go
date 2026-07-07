package v3

import (
	"context"
	"maps"
)

// FlowContextOption configures a new BranchingFlowContext branch.
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
