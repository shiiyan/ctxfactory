package ctxfactory

import (
	"context"
	"maps"
	"time"
)

type ContextFactory struct {
	defaults  map[any]any
	overrides map[any]any
	skipKeys  map[any]struct{}
}

// New creates a context factory with arbitrary default values.
func New(defaults map[any]any) *ContextFactory {
	if defaults == nil {
		defaults = map[any]any{}
	}
	return &ContextFactory{
		defaults:  defaults,
		overrides: map[any]any{},
		skipKeys:  map[any]struct{}{},
	}
}

// With merges overrides into the factory and returns the same factory for fluent use.
func (f *ContextFactory) With(overrides map[any]any) *ContextFactory {
	if overrides == nil {
		return f
	}
	if f.overrides == nil {
		f.overrides = map[any]any{}
	}
	maps.Copy(f.overrides, overrides)
	return f
}

// Skip marks keys to omit from applying defaults.
func (f *ContextFactory) Skip(keys ...any) *ContextFactory {
	if f.skipKeys == nil {
		f.skipKeys = map[any]struct{}{}
	}
	for _, k := range keys {
		f.skipKeys[k] = struct{}{}
	}
	return f
}

// Build constructs a context from defaults and overrides, honoring skips.
func (f *ContextFactory) Build() context.Context {
	return f.BuildWith(context.TODO())
}

// BuildWith constructs a context using the provided base context.
func (f *ContextFactory) BuildWith(base context.Context) context.Context {
	if base == nil || base == context.TODO() {
		base = context.Background()
	}
	ctx := base

	// apply defaults unless skipped
	for k, v := range f.defaults {
		if _, skip := f.skipKeys[k]; skip {
			continue
		}
		ctx = context.WithValue(ctx, k, v)
	}

	// apply overrides (they replace defaults)
	for k, v := range f.overrides {
		ctx = context.WithValue(ctx, k, v)
	}

	return ctx
}

// BuildWithCancel wraps the built context with a cancel function.
func (f *ContextFactory) BuildWithCancel(base context.Context) (context.Context, context.CancelFunc) {
	ctx := f.BuildWith(base)
	return context.WithCancel(ctx)
}

// BuildWithTimeout wraps the built context with a timeout and returns the cancel func.
func (f *ContextFactory) BuildWithTimeout(base context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	ctx := f.BuildWith(base)
	return context.WithTimeout(ctx, d)
}

// BuildWithDeadline wraps the built context with a deadline and returns the cancel func.
func (f *ContextFactory) BuildWithDeadline(base context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	ctx := f.BuildWith(base)
	return context.WithDeadline(ctx, deadline)
}
