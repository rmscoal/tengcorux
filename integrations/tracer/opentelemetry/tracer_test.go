package opentelemetry

import (
	"context"
	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	"testing"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestTracer_StartSpan(t *testing.T) {
	tt := tracetest.NewNoopExporter()
	tracer := NewTracer("testing", WithExporter(tt))

	t.Run("Normal(Noop)", func(t *testing.T) {
		ctx, span := tracer.StartSpan(context.TODO(), "test")
		if ctx == nil {
			t.Fatal("expected non-nil context")
		} else if span == nil {
			t.Fatal("expected non-nil span")
		}

		traceSpan := span.(*Span).span
		if traceSpan == nil {
			t.Fatal("expected a trace span")
		} else if !traceSpan.IsRecording() {
			t.Skipf("trce span is not recording")
		}

		if !traceSpan.SpanContext().HasTraceID() {
			t.Fatal("span does not have a traceID")
		} else if !traceSpan.SpanContext().HasSpanID() {
			t.Fatal("span does not have a spanID")
		}

		otelSpan := traceSpan.(sdktrace.ReadWriteSpan)
		if otelSpan.Name() != "test" {
			t.Fatal("span name is not equal to test")
		} else if otelSpan.SpanKind() != trace.SpanKindInternal {
			t.Fatalf("span kind is not internal, instead %v", otelSpan.SpanKind())
		}
	})

	t.Run("WithTraceID", func(t *testing.T) {
		givenTraceID := "5b8aa5a2d2c872e8321cf37308d69df2"
		ctx, span := tracer.StartSpan(context.TODO(), "test",
			tengcoruxTracer.WithTraceID(givenTraceID))
		if ctx == nil {
			t.Fatal("expected non-nil context")
		} else if span == nil {
			t.Fatal("expected non-nil span")
		}

		otelSpan := span.(*Span).span.(sdktrace.ReadWriteSpan)
		if otelSpanTraceID := otelSpan.SpanContext().
			TraceID().
			String(); otelSpanTraceID != givenTraceID {
			t.Fatalf("expected traceID to be %s but got %s",
				givenTraceID, otelSpanTraceID)
		}
	})

	t.Run("WithInvalidTraceID", func(t *testing.T) {
		invalidTraceID := "invalid-trace-id"
		_, span := tracer.StartSpan(context.TODO(), "test",
			tengcoruxTracer.WithTraceID(invalidTraceID))
		if span == nil {
			t.Fatal("expected non-nil span even with invalid trace ID")
		}
		// Verify that an invalid trace ID results in a new valid trace ID being generated
		otelSpan := span.(*Span).span.(sdktrace.ReadWriteSpan)
		if !otelSpan.SpanContext().TraceID().IsValid() {
			t.Fatal("expected a valid trace ID to be generated for invalid input")
		}
	})
}

func TestTracer_Shutdown(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("test panicked: %v", r)
		}
	}()

	t.Run("ContextBackground", func(t *testing.T) {
		tracer := NewTracer("testing")
		err := tracer.Shutdown(context.Background())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("ContextWithCancel", func(t *testing.T) {
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatal(err)
		}

		t.Run("Cancelled", func(t *testing.T) {
			tracer := NewTracer("testing", WithExporter(exporter))
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err = tracer.Shutdown(ctx)
			if err == nil {
				t.Errorf("expected a context cancelled error")
			}
		})

		t.Run("Normal", func(t *testing.T) {
			tracer := NewTracer("testing", WithExporter(exporter))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			err := tracer.Shutdown(ctx)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	})

	t.Run("ContextWithTimeout", func(t *testing.T) {
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatal(err)
		}
		tracer := NewTracer("testing", WithExporter(exporter))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		err = tracer.Shutdown(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestTracer_SpanFromContext(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("test panicked: %v", r)
		}
	}()

	tracer := NewTracer("testing")

	// NOTE: For now all type of context always return a noop span by otel.

	t.Run("Nil", func(t *testing.T) {
		span := tracer.SpanFromContext(nil)
		if span == nil {
			t.Fatal("expected non-nil span")
		}
	})

	t.Run("EmptyContext", func(t *testing.T) {
		span := tracer.SpanFromContext(context.Background())
		if span == nil {
			t.Fatal("expected non-nil span")
		}
	})

	t.Run("ActiveSpan", func(t *testing.T) {
		ctx, _ := tracer.StartSpan(context.Background(), "test")
		span := tracer.SpanFromContext(ctx)
		if span == nil {
			t.Fatal("expected non-nil span")
		}
	})
}
