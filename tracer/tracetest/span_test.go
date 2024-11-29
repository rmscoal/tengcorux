package tracetest

import (
	"context"
	"errors"
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"testing"
)

func TestSpan_End(t *testing.T) {
	tracer := NewTracer()

	span := &Span{tracer: tracer}
	span.End()

	if span.EndTime.IsZero() {
		t.Error("end time should not be zero")
	}
	if len(span.tracer.recorder.EndedSpans()) == 0 {
		t.Error("ended spans should not be empty")
	}
}

func TestSpan_SetAttributes(t *testing.T) {
	span := &Span{}
	span.SetAttributes(
		attribute.KeyValuePair("Hello1", "World1"),
		attribute.KeyValuePair("Hello2", "World2"),
	)

	if len(span.Attributes) != 2 {
		t.Errorf("expected span to have 2 attributes, but got %d", len(span.Attributes))
	}

	for _, kv := range span.Attributes {
		switch kv.Key {
		case "Hello1":
			if kv.Value != "World1" {
				t.Errorf("expected attribute value to be \"World1\", but got %s", kv.Value)
			}
		case "Hello2":
			if kv.Value != "World2" {
				t.Errorf("expected attribute value to be \"World2\", but got %s", kv.Value)
			}
		default:
			t.Errorf("unknown attribute key: %s", kv.Key)
		}
	}
}

func TestSpan_AddEvent(t *testing.T) {
	span := &Span{}
	span.AddEvent("some event is happening here")
	span.AddEvent("another event is happening here")

	if len(span.Events) != 2 {
		t.Errorf("expected span to have 2 events, but got %d", len(span.Events))
	}

	for _, ev := range span.Events {
		switch ev {
		case "some event is happening here", "another event is happening here":
		default:
			t.Errorf("unknown event: %s", ev)
		}
	}
}

func TestSpan_RecordError(t *testing.T) {
	span := &Span{}
	err := errors.New("some error")
	span.RecordError(err)
	if span.Error == nil {
		t.Error("expected error to be non-nil")
	} else if !errors.Is(span.Error, err) {
		t.Errorf("expected error to be %v, got %v", err, span.Error)
	}
}

func TestSpan_Context(t *testing.T) {
	span := &Span{
		spanContext: &SpanContext{
			ctx: context.Background(),
		},
	}
	ctx := span.Context()
	if ctx == nil {
		t.Error("expected context to be non-nil")
	}
}

func TestSpanContext_Context(t *testing.T) {
	tracer := NewTracer()
	_, span := tracer.StartSpan(context.Background(), "hello")

	ctx := span.Context().Context()
	if ctx == nil {
		t.Error("expected context to be non-nil")
	}
}

func TestSpanContext_TraceID(t *testing.T) {
	tracer := NewTracer()
	_, span := tracer.StartSpan(context.Background(), "hello")

	sc := span.Context()
	if sc.TraceID() == "" {
		t.Error("expected trace ID to be non-empty string")
	}
}

func TestSpanContext_SpanID(t *testing.T) {
	tracer := NewTracer()
	_, span := tracer.StartSpan(context.Background(), "hello")

	sc := span.Context()
	if sc.SpanID() == "" {
		t.Error("expected trace ID to be non-empty string")
	}
}
