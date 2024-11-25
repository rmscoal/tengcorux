package tracer

import (
	"context"
	"testing"
)

func TestStartSpan(t *testing.T) {
	// Cleanup function to catch if there are panics as it should not panic.
	defer t.Cleanup(func() {
		if m := recover(); m != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	t.Run("WithoutOptions", func(t *testing.T) {
		ctx, span := StartSpan(context.Background(), "testing")
		if ctx == nil {
			t.Error("ctx should not be nil")
		} else if span == nil {
			t.Error("span should not be nil")
		}

		spanCtx := span.Context()
		if spanCtx == nil {
			t.Error("span context should not be nil")
		}
	})

	t.Run("WithOptions", func(t *testing.T) {
		ctx, span := StartSpan(context.Background(), "testing",
			WithSpanType(SpanTypeLocal),
			WithSpanLayer(SpanLayerDatabase),
		)
		if ctx == nil {
			t.Error("ctx should not be nil")
		} else if span == nil {
			t.Error("span should not be nil")
		}

		spanCtx := span.Context()
		if spanCtx == nil {
			t.Error("span context should not be nil")
		}
	})
}

func TestShutdown(t *testing.T) {
	defer t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	t.Run("With Non Empty Context", func(t *testing.T) {
		ctx := context.Background()
		err := Shutdown(ctx)
		if err != nil {
			t.Error("should not return an error")
		}
	})

	t.Run("With Empty Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := Shutdown(ctx)
		if err == nil {
			t.Error("should return an error of context timeout")
		}
	})
}

func TestSpanFromContext(t *testing.T) {
	defer t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	span := SpanFromContext(context.Background())
	if span == nil {
		t.Error("span should not be nil")
	}
}
