package skywalking

import (
	"context"
	"errors"
	"testing"

	tengcoruxAttribute "github.com/rmscoal/tengcorux/tracer/attribute"
)

func TestSkyWalkingSpan_End(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, _ := startTestingTracer()
	ctx, span := tracer.StartSpan(context.Background(), "testing")
	if ctx == nil {
		t.Error("should not throw error")
	}

	span.End()
}

func TestSkyWalkingSpan_SetAttributes(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, _ := startTestingTracer()
	ctx, span := tracer.StartSpan(context.Background(), "testing")
	if ctx == nil {
		t.Error("should not throw error")
	}

	span.SetAttributes(tengcoruxAttribute.DBSystem("some_system"))
}

func TestSkyWalkingSpan_RecordError(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, _ := startTestingTracer()
	ctx, span := tracer.StartSpan(context.Background(), "testing")
	if ctx == nil {
		t.Error("should not throw error")
	}

	span.RecordError(errors.New("some_error"))
}

func TestSkyWalkingSpan_AddEvent(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, _ := startTestingTracer()
	ctx, span := tracer.StartSpan(context.Background(), "testing")
	if ctx == nil {
		t.Error("should not throw error")
	}

	span.AddEvent("some event is happening now")
}

func TestSkyWalkingSpan_Context(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, _ := startTestingTracer()
	ctx, span := tracer.StartSpan(context.Background(), "testing")
	if ctx == nil {
		t.Error("should not throw error")
	}

	_, ok := span.Context().(*SpanContext)
	if !ok {
		t.Error("context should be of type SpanContext")
	}
}

func TestSkyWalkingSpanContext(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, _ := startTestingTracer()
	ctx, span := tracer.StartSpan(context.Background(), "testing")
	if ctx == nil {
		t.Error("should not throw error")
	}

	sc, ok := span.Context().(*SpanContext)
	if !ok {
		t.Error("context should be of type SpanContext")
	}

	t.Run("Empty", func(t *testing.T) {
		sc := &SpanContext{
			ctx: context.Background(),
		}

		t.Run("TraceID", func(t *testing.T) {
			traceId := sc.TraceID()
			if traceId != "" {
				t.Errorf("traceId should be empty, but got %s", traceId)
			}
		})
		t.Run("SpanID", func(t *testing.T) {
			spanId := sc.SpanID()
			if spanId != "" {
				t.Error("spanId should be empty")
			}
		})
		t.Run("Context", func(t *testing.T) {
			ctx := sc.Context()
			if ctx == nil {
				t.Error("ctx should not be nil")
			}
		})
	})

	// This fails, TraceID is always empty
	t.Run("TraceID", func(t *testing.T) {
		traceId := sc.TraceID()
		if traceId == "" {
			t.Error("traceId should not be empty")
		}
	})
	t.Run("SpanID", func(t *testing.T) {
		spanId := sc.SpanID()
		if spanId == "" {
			t.Error("spanId should not be empty")
		}
	})
	t.Run("Context", func(t *testing.T) {
		ctx := sc.Context()
		if ctx == nil {
			t.Error("ctx should not be nil")
		}
	})
}
