package tracetest

import (
	"context"
	"testing"

	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
)

func TestTracer_StartSpan(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("should not panic, but occurred: ", err)
		}
	}()

	tr := NewTracer()
	ctx, span := tr.StartSpan(context.Background(), "test")
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
	} else if testSpan.Layer != tengcoruxTracer.SpanLayerUnknown {
		t.Fatal("expected unknown layer")
	} else if testSpan.Type != tengcoruxTracer.SpanTypeLocal {
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
}

func TestTracer_Shutdown(t *testing.T) {
	tr := NewTracer()
	err := tr.Shutdown(context.Background())
	if err != nil {
		t.Fatal("should not error, but occurred: ", err)
	}
}

func TestTracer_SpanFromContext(t *testing.T) {
	tr := NewTracer()

	ctx, _ := tr.StartSpan(context.Background(), "test")
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
}
