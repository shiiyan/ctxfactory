package ctxfactory

import (
	"context"
	"testing"
	"time"
)

func TestBuildAppliesDefaults(t *testing.T) {
	t.Parallel()

	f := NewContextFactory(map[any]any{
		"user":  1,
		"trace": "t1",
	})

	ctx := f.Build()
	if v := ctx.Value("user"); v != 1 {
		t.Fatalf("expected user=1 got %#v", v)
	}
	if v := ctx.Value("trace"); v != "t1" {
		t.Fatalf("expected trace=t1 got %#v", v)
	}
}

func TestWithOverridesReplaceDefaults(t *testing.T) {
	t.Parallel()

	f := NewContextFactory(map[any]any{"user": 1})
	f.With(map[any]any{"user": 42})

	ctx := f.Build()
	if v := ctx.Value("user"); v != 42 {
		t.Fatalf("expected user=42 got %#v", v)
	}
}

func TestSkipOmitsDefaultKeys(t *testing.T) {
	t.Parallel()

	f := NewContextFactory(map[any]any{
		"user":  1,
		"trace": "t1",
	})
	f.Skip("trace")

	ctx := f.Build()
	if v := ctx.Value("trace"); v != nil {
		t.Fatalf("expected trace to be skipped, got %#v", v)
	}
	// other keys still present
	if v := ctx.Value("user"); v != 1 {
		t.Fatalf("expected user=1 got %#v", v)
	}
}

type baseKeyType struct{}

func TestBuildWithBaseContext(t *testing.T) {
	t.Parallel()

	baseKey := baseKeyType{}
	base := context.WithValue(context.Background(), baseKey, "ok")
	f := NewContextFactory(map[any]any{"user": 1})

	ctx := f.BuildWith(base)
	if v := ctx.Value(baseKey); v != "ok" {
		t.Fatalf("expected base value propagated, got %#v", v)
	}
	if v := ctx.Value("user"); v != 1 {
		t.Fatalf("expected user=1 got %#v", v)
	}
}

func TestBuildWithCancelReturnsCancelableContext(t *testing.T) {
	t.Parallel()

	f := NewContextFactory(nil)

	ctx, cancel := f.BuildWithCancel(context.TODO())
	// cancel and ensure Done is closed
	cancel()

	select {
	case <-ctx.Done():
		// ok
	case <-time.After(50 * time.Millisecond):
		t.Fatal("context not canceled after cancel()")
	}
}

func TestBuildWithTimeoutSetsDeadline(t *testing.T) {
	t.Parallel()

	f := NewContextFactory(nil)

	ctx, cancel := f.BuildWithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		t.Fatal("expected a deadline to be set")
	}

	// ensure context is done after the timeout
	select {
	case <-ctx.Done():
		// ok (may happen quickly)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("context did not expire after timeout")
	}
}

func TestBuildWithDeadlineSetsDeadline(t *testing.T) {
    t.Parallel()

    f := NewContextFactory(nil)

    deadline := time.Now().Add(10 * time.Millisecond)
    ctx, cancel := f.BuildWithDeadline(context.TODO(), deadline)
    defer cancel()

    if _, ok := ctx.Deadline(); !ok {
        t.Fatal("expected a deadline to be set")
    }

    // ensure context is done after the deadline
    select {
    case <-ctx.Done():
        // ok
    case <-time.After(200 * time.Millisecond):
        t.Fatal("context did not expire after deadline")
    }
}
