package opentelemetry

import (
	"context"
	"errors"
	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

func TestSpan(t *testing.T) {
	defer func(t *testing.T) {
		if r := recover(); r != nil {
			t.Errorf("test panicked: %v", r)
		}
	}(t)

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		t.Fatal(err)
	}
	tr := NewTracer("testing", WithExporter(exporter))
	_, span := tr.StartSpan(context.Background(), "test")

	t.Run("End", func(t *testing.T) {
		span.End()
		if span.(*Span).span.IsRecording() {
			t.Error("expected the span to stop recording")
		}
	})

	t.Run("SetAttributes", func(t *testing.T) {
		span.SetAttributes(
			attribute.KeyValuePair("key", "value"),
			attribute.KeyValuePair("key1", 1),
			attribute.KeyValuePair("key2", 4.35443),
			attribute.KeyValuePair("key3", map[string]interface{}{"hello": "world"}),
		)
	})

	t.Run("RecordError", func(t *testing.T) {
		span.RecordError(errors.New("error"))
		span.RecordError(nil)
	})

	t.Run("AddEvent", func(t *testing.T) {
		span.AddEvent("hello 1", "hello 2", "hello 3")
		span.AddEvent()
	})

	t.Run("Context", func(t *testing.T) {
		sc := span.Context()
		if sc == nil {
			t.Error("expected a span context")
		}
	})
}

func TestSpanContext(t *testing.T) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		t.Fatal(err)
	}
	tr := NewTracer("testing", WithExporter(exporter))
	_, span := tr.StartSpan(context.Background(), "test")

	t.Run("TraceID", func(t *testing.T) {
		if trID := span.Context().TraceID(); trID == "" {
			t.Error("expected a non-empty trace ID")
		}
	})
	t.Run("SpanID", func(t *testing.T) {
		if trID := span.Context().SpanID(); trID == "" {
			t.Error("expected a non-empty span ID")
		}
	})
	t.Run("Context", func(t *testing.T) {
		if ctx := span.Context().Context(); ctx == nil {
			t.Error("expected a non-empty context.Context")
		}
	})
}

func TestMapSpanKind(t *testing.T) {
	t.Run("SpanTypeLocal", func(t *testing.T) {
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeLocal,
			tengcoruxTracer.SpanLayerUnknown); kind != trace.SpanKindInternal {
			t.Errorf("expected %s but got %s", trace.SpanKindInternal, kind)
		}
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeLocal,
			tengcoruxTracer.SpanLayerDatabase); kind != trace.SpanKindInternal {
			t.Errorf("expected %s but got %s", trace.SpanKindInternal, kind)
		}
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeLocal,
			tengcoruxTracer.SpanLayerHttp); kind != trace.SpanKindInternal {
			t.Errorf("expected %s but got %s", trace.SpanKindInternal, kind)
		}
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeLocal,
			tengcoruxTracer.SpanLayerMQ); kind != trace.SpanKindInternal {
			t.Errorf("expected %s but got %s", trace.SpanKindInternal, kind)
		}
	})
	t.Run("SpanTypeEntry", func(t *testing.T) {
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeEntry,
			tengcoruxTracer.SpanLayerUnknown); kind != trace.SpanKindServer {
			t.Errorf("expected %s but got %s", trace.SpanKindServer, kind)
		}
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeEntry,
			tengcoruxTracer.SpanLayerHttp); kind != trace.SpanKindServer {
			t.Errorf("expected %s but got %s", trace.SpanKindServer, kind)
		}
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeEntry,
			tengcoruxTracer.SpanLayerDatabase); kind != trace.SpanKindServer {
			t.Errorf("expected %s but got %s", trace.SpanKindServer, kind)
		}
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeEntry,
			tengcoruxTracer.SpanLayerMQ); kind != trace.SpanKindConsumer {
			t.Errorf("expected %s but got %s", trace.SpanKindConsumer, kind)
		}
	})
	t.Run("SpanTypeExit", func(t *testing.T) {
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeExit,
			tengcoruxTracer.SpanLayerUnknown); kind != trace.SpanKindClient {
			t.Errorf("expected %s but got %s", trace.SpanKindClient, kind)
		}
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeExit,
			tengcoruxTracer.SpanLayerHttp); kind != trace.SpanKindClient {
			t.Errorf("expected %s but got %s", trace.SpanKindClient, kind)
		}
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeExit,
			tengcoruxTracer.SpanLayerDatabase); kind != trace.SpanKindClient {
			t.Errorf("expected %s but got %s", trace.SpanKindClient, kind)
		}
		if kind := mapSpanKind(tengcoruxTracer.SpanTypeExit,
			tengcoruxTracer.SpanLayerMQ); kind != trace.SpanKindProducer {
			t.Errorf("expected %s but got %s", trace.SpanKindProducer, kind)
		}
	})
}
