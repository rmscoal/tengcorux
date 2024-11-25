package tracer

import (
	"context"

	"github.com/rmscoal/tengcorux/tracer/attribute"
)

// NoopTracer a no operation tracer that implements [Tracer].
type NoopTracer struct{}

// Make sure that NoopTracer implements [Tracer] during compile time.
var _ Tracer = (*NoopTracer)(nil)

// StartSpan does returns the context itself and a [NoopSpan].
func (t *NoopTracer) StartSpan(ctx context.Context, _ string, _ ...StartSpanOption) (context.Context, Span) {
	return ctx, &NoopSpan{}
}

func (t *NoopTracer) Shutdown(_ context.Context) error {
	return nil
}

func (t *NoopTracer) SpanFromContext(ctx context.Context) Span {
	return &NoopSpan{}
}

// NoopSpan n no operation span that implements [Span].
type NoopSpan struct{}

// Make sure that NoopSpan implements [Span] during compile time.
var _ Span = (*NoopSpan)(nil)

// End does nothing.
func (s *NoopSpan) End() {}

// SetAttributes does nothing.
func (s *NoopSpan) SetAttributes(_ ...attribute.KeyValue) {}

// RecordError does nothing.
func (s *NoopSpan) RecordError(_ error) {}

// AddEvent does nothing.
func (s *NoopSpan) AddEvent(_ ...string) {}

// Context returns an empty context.
func (s *NoopSpan) Context() SpanContext {
	return &NoopSpanContext{}
}

// NoopSpanContext a no operation span context that implements [SpanContext].
type NoopSpanContext struct{}

// Make sure that NoopSpanContext implements [SpanContext] during compile time.
var _ SpanContext = (*NoopSpanContext)(nil)

// TraceID returns an empty string.
func (s *NoopSpanContext) TraceID() string { return "" }

// SpanID returns an empty string.
func (s *NoopSpanContext) SpanID() string { return "" }

// Context returns an empty context.
func (s *NoopSpanContext) Context() context.Context {
	return context.TODO()
}
