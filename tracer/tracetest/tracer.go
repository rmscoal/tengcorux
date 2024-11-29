package tracetest

import (
	"context"
	"time"

	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
)

type Tracer struct {
	recorder *SpanRecorder
}

// Checks if our test tracer implements tengcorux tracer interface.
var _ tengcoruxTracer.Tracer = (*Tracer)(nil)

// NewTracer returns a test trace instance with a new span recorder.
func NewTracer() *Tracer {
	return &Tracer{
		recorder: NewSpanRecorder(),
	}
}

// StartSpan starts a test span and insert the span value into the context.
// It also inserts the span into the recorder slice of spans.
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...tengcoruxTracer.StartSpanOption) (context.Context, tengcoruxTracer.Span) {
	spanConfig := tengcoruxTracer.DefaultStartSpanConfig()
	for _, opt := range opts {
		opt(spanConfig)
	}

	span := &Span{
		StartTime: time.Now(),
		Name:      name,
		TraceID:   newRandomIntegerID(),
		SpanID:    newRandomIntegerID(),
		Layer:     spanConfig.SpanLayer,
		Type:      spanConfig.SpanType,
	}

	// Search for the previous span in the context and adjust values
	// for current span if found.
	prevSpan, exists := ctx.Value(prevSpanKey).(*Span)
	if exists && prevSpan != nil {
		span.TraceID = prevSpan.TraceID
		span.ParentSpanID = prevSpan.SpanID
	}

	// Replaces the context's prevSpanKey with the current span.
	ctx = context.WithValue(ctx, prevSpanKey, span)
	span.spanContext = &SpanContext{ctx: ctx}
	t.recorder.OnStart(span)

	return ctx, span
}

// Shutdown does nothing and returns nil.
func (t *Tracer) Shutdown(_ context.Context) error { return nil }

// SpanFromContext searches for the previous span value in the context
// and returns the span. If non-existing, it will return nil.
func (t *Tracer) SpanFromContext(ctx context.Context) tengcoruxTracer.Span {
	if ctx == nil {
		return nil
	}

	span, ok := ctx.Value(prevSpanKey).(*Span)
	if !ok {
		return nil
	}

	return span
}
