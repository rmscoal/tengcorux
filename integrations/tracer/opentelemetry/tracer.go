package opentelemetry

import (
	"context"

	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tracer      trace.Tracer
	shutdowns   []func(context.Context) error
	serviceName string
	environment string
}

// StartSpan starts a new Span with the given name and option.
func (t *Tracer) StartSpan(ctx context.Context, name string,
	opts ...tengcoruxTracer.StartSpanOption,
) (context.Context, tengcoruxTracer.Span) {
	startSpanConfig := tengcoruxTracer.DefaultStartSpanConfig()
	for _, opt := range opts {
		opt(startSpanConfig)
	}

	ctx, span := t.tracer.Start(
		generateContextFromStartSpanConfig(ctx, startSpanConfig),
		name,
		trace.WithSpanKind(mapSpanKind(startSpanConfig.SpanType,
			startSpanConfig.SpanLayer)),
	)

	return ctx, &Span{
		tracer: t,
		span:   span,
		spanContext: &SpanContext{
			ctx: ctx,
		},
	}
}

// Shutdown closes all the exporter shutdowns with the given context.
func (t *Tracer) Shutdown(ctx context.Context) error {
	for _, shutdown := range t.shutdowns {
		if err := shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}

// SpanFromContext return a span from a given context. Currently, whether the
// context has an active span or empty (even nil), it will always return a span
// although a noop span by otel.
func (t *Tracer) SpanFromContext(ctx context.Context) tengcoruxTracer.Span {
	span := trace.SpanFromContext(ctx)

	// NOTE: otel trace always return a span, however when the context is empty
	// it returns a noop span instead. Should we return nil if it is a noop or
	// just return the span?
	//
	// if !span.IsRecording() { return nil }

	return &Span{
		tracer: t,
		span:   span,
		spanContext: &SpanContext{
			ctx: ctx,
		},
	}
}

// generateContextFromStartSpanConfig generates a new context only if
// the start span config given includes a TraceID and/or ParentSpanID.
func generateContextFromStartSpanConfig(ctx context.Context,
	tengcoruxStartSpanConfig *tengcoruxTracer.StartSpanConfig,
) context.Context {
	if tengcoruxStartSpanConfig.TraceID != "" {
		traceSpanConfig := trace.SpanContextConfig{
			Remote:     true,
			TraceFlags: trace.FlagsSampled,
		}

		traceID, err := trace.TraceIDFromHex(tengcoruxStartSpanConfig.TraceID)
		if err == nil {
			traceSpanConfig.TraceID = traceID
		}

		if tengcoruxStartSpanConfig.ParentSpanID != "" {
			spanID, err := trace.SpanIDFromHex(tengcoruxStartSpanConfig.ParentSpanID)
			if err == nil {
				traceSpanConfig.SpanID = spanID
			}
		}

		// At the end, we check whether the traceID is valid or not
		// before we generate/replace the given context.
		if traceSpanConfig.TraceID.IsValid() {
			traceSpanCtx := trace.NewSpanContext(traceSpanConfig)
			ctx = trace.ContextWithSpanContext(ctx, traceSpanCtx)
		}
	}

	return ctx
}
