package tracer

import (
	"context"
)

type Tracer interface {
	// StartSpan creates a new parent span.
	StartSpan(ctx context.Context, name string, opts ...StartSpanOption) (context.Context, Span)

	// Shutdown stops the tracer from receiving and exporting tracer.
	Shutdown(ctx context.Context) error

	// SpanFromContext retrieves the Span from a given context.Context.
	SpanFromContext(ctx context.Context) Span
}

func StartSpan(ctx context.Context, name string, opts ...StartSpanOption) (context.Context, Span) {
	return GetGlobalTracer().StartSpan(ctx, name, opts...)
}

func Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return GetGlobalTracer().Shutdown(ctx)
	}
}

func SpanFromContext(ctx context.Context) Span {
	return GetGlobalTracer().SpanFromContext(ctx)
}
