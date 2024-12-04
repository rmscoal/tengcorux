package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"testing"
	"time"
)

func TestTracer_StartSpan(t *testing.T) {
	tracer := NewTracer("testing")
	ctx, span := tracer.StartSpan(context.TODO(), "test")
	if ctx == nil {
		t.Fatal("expected non-nil context")
	} else if span == nil {
		t.Fatal("expected non-nil span")
	}
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
