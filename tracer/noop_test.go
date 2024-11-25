package tracer

import (
	"context"
	"errors"
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"testing"
)

func TestNoopTracer_StartSpan(t *testing.T) {
	noop := &NoopTracer{}
	ctx, span := noop.StartSpan(context.Background(), "testing_noop",
		WithSpanType(SpanTypeEntry),
		WithSpanLayer(SpanLayerUnknown),
	)
	if ctx == nil {
		t.Error("ctx should not be nil")
	} else if span == nil {
		t.Error("span should not be nil")
	}

	_, ok := span.(*NoopSpan)
	if !ok {
		t.Error("span should be of type (*NoopSpan)")
	}
}

func TestNoopTracer_Shutdown(t *testing.T) {
	noop := &NoopTracer{}

	err := noop.Shutdown(context.Background())
	if err != nil {
		t.Error("error should be nil")
	}
}

func TestNoopTracer_SpanFromContext(t *testing.T) {
	noop := &NoopTracer{}

	span := noop.SpanFromContext(context.Background())
	_, ok := span.(*NoopSpan)
	if !ok {
		t.Error("span should be of type *NoopSpan")
	}
}

func TestNoopSpan_SetAttributes(t *testing.T) {
	defer t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	span := &NoopSpan{}
	span.SetAttributes(attribute.DBName("some_db_name"))
	span.SetAttributes(attribute.DBOperation("some_operation"))
}

func TestNoopSpan_AddEvent(t *testing.T) {
	defer t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	span := &NoopSpan{}
	span.AddEvent("event_1")
	span.AddEvent("event_2")
}

func TestNoopSpan_RecordError(t *testing.T) {
	defer t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	span := &NoopSpan{}
	span.RecordError(errors.New("some_error"))
}

func TestNoopSpan_Context(t *testing.T) {
	defer t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	span := &NoopSpan{}
	_, ok := span.Context().(*NoopSpanContext)
	if !ok {
		t.Error("should be of type (*NoopSpanContext)")
	}
}

func TestNoopSpanContext_TraceID(t *testing.T) {
	defer t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	sc := &NoopSpanContext{}
	traceId := sc.TraceID()
	if traceId != "" {
		t.Error("traceId should be empty")
	}
}

func TestNoopSpanContext_SpanID(t *testing.T) {
	defer t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	sc := &NoopSpanContext{}
	spanId := sc.SpanID()
	if spanId != "" {
		t.Error("spanId should be empty")
	}
}

func TestNoopSpanContext_Context(t *testing.T) {
	defer t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	})

	sc := &NoopSpanContext{}
	ctx := sc.Context()
	if ctx == nil {
		t.Error("ctx should not be nil")
	}
}
