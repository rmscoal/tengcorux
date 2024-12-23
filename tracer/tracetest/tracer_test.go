package tracetest

import (
	"context"
	"testing"

	"github.com/rmscoal/tengcorux/tracer"
)

func TestTracer_StartSpan(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("should not panic, but occurred: ", err)
		}
	}()

	tr := NewTracer()

	t.Run("From Empty Context", func(t *testing.T) {
		ctx, span := tr.StartSpan(context.TODO(), "test")
		if ctx == nil {
			t.Fatal("expected non-nil context")
		} else if span == nil {
			t.Fatal("expected non-nil span")
		}

		testSpan, ok := span.(*Span)
		if !ok {
			t.Fatal("span is not a Span")
		}

		// Checks testSpan each fields.
		if testSpan.SpanID == 0 {
			t.Fatal("expected non-zero SpanID")
		} else if testSpan.TraceID == 0 {
			t.Fatal("expected non-zero TraceID")
		} else if testSpan.ParentSpanID != 0 {
			t.Fatal("expected zero ParentSpanID")
		} else if testSpan.Name != "test" {
			t.Fatal("expected test span")
		} else if testSpan.StartTime.IsZero() {
			t.Fatal("expected non-zero StartTime")
		} else if !testSpan.EndTime.IsZero() {
			t.Fatal("expected a zero EndTime")
		} else if testSpan.Layer != tracer.SpanLayerUnknown {
			t.Fatal("expected unknown layer")
		} else if testSpan.Type != tracer.SpanTypeLocal {
			t.Fatal("expected local span")
		} else if testSpan.tracer == nil {
			t.Fatal("expected non-nil tracer")
		}

		// Checks the context whether it contains prevSpanKey.
		spanInCtx, ok := ctx.Value(prevSpanKey).(*Span)
		if !ok {
			t.Fatal("expected context to contain prevSpanKey")
		} else if spanInCtx == nil {
			t.Fatal("expected span to not be nil in context")
		} else if spanInCtx.SpanID != testSpan.SpanID {
			t.Fatal("expected span to be the same in context and the one returned")
		}
	})

	t.Run("From Existing Span", func(t *testing.T) {
		prevCtx, _ := tr.StartSpan(context.TODO(), "span_number_one")
		ctx, span := tr.StartSpan(prevCtx, "span_number_two")

		testSpan, ok := span.(*Span)
		if !ok {
			t.Fatal("span is not a Span")
		} else if testSpan.Name != "span_number_two" {
			t.Fatal("expected span_number_two span")
		}

		// Checks the context whether it contains prevSpanKey.
		spanInCtx, ok := ctx.Value(prevSpanKey).(*Span)
		if !ok {
			t.Fatal("expected context to contain prevSpanKey")
		} else if spanInCtx == nil {
			t.Fatal("expected span to not be nil in context")
		} else if spanInCtx.SpanID != testSpan.SpanID {
			t.Fatal("expected span to be the same in context and the one returned")
		}
	})

	t.Run("WithOptions", func(t *testing.T) {
		_, span := tr.StartSpan(context.TODO(), "span_number_one",
			tracer.WithSpanLayer(tracer.SpanLayerDatabase),
			tracer.WithSpanType(tracer.SpanTypeExit))
		testSpan := span.(*Span)
		if testSpan.Layer != tracer.SpanLayerDatabase {
			t.Errorf("expected layer to be database but got %d", testSpan.Layer)
		} else if testSpan.Type != tracer.SpanTypeExit {
			t.Errorf("expected type to be database but got %d", testSpan.Type)
		}
	})
}

func TestTracer_Shutdown(t *testing.T) {
	tr := NewTracer()
	err := tr.Shutdown(context.Background())
	if err != nil {
		t.Fatal("should not error, but occurred: ", err)
	}
}

func TestTracer_SpanFromContext(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("should not panic, but occurred: ", err)
		}
	}()

	tr := NewTracer()

	t.Run("Nil", func(t *testing.T) {
		span := tr.SpanFromContext(nil)
		if span != nil {
			t.Fatal("expected nil span")
		}
	})

	t.Run("Empty Context", func(t *testing.T) {
		span := tr.SpanFromContext(context.TODO())
		if span != nil {
			t.Fatal("expected nil span")
		}
	})

	t.Run("From Span Context", func(t *testing.T) {
		ctx, _ := tr.StartSpan(context.TODO(), "test")
		span := tr.SpanFromContext(ctx)
		if span == nil {
			t.Fatal("expected non-nil span")
		}

		testSpan, ok := span.(*Span)
		if !ok {
			t.Fatal("span is not a Span")
		} else if testSpan.SpanID == 0 {
			t.Fatal("expected span id to not be zero")
		} else if testSpan.TraceID == 0 {
			t.Fatal("expected trace id to not be zero")
		} else if testSpan.Name != "test" {
			t.Fatal("expected test span name to be 'test'")
		}
	})
}

func TestTracer_Recorder(t *testing.T) {
	tr := NewTracer()
	rec := tr.Recorder()
	if rec == nil {
		t.Fatal("expected non-nil recorder")
	}
}
