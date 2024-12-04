package opentelemetry

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"time"

	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	tengcoruxAttribute "github.com/rmscoal/tengcorux/tracer/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	tracer      *Tracer
	span        trace.Span
	spanContext *SpanContext
}

// End ends the current span.
func (s *Span) End() {
	s.span.End()
}

// SetAttributes sets attributes to the current span.
func (s *Span) SetAttributes(tengcoruxAttributes ...tengcoruxAttribute.KeyValue) {
	var attributes []attribute.KeyValue

	for _, attr := range tengcoruxAttributes {
		key := string(attr.Key)
		switch val := attr.Value.(type) {
		case bool:
			attributes = append(attributes, attribute.Bool(key, val))
		case string:
			attributes = append(attributes, attribute.String(key, val))
		case int:
			attributes = append(attributes, attribute.Int(key, val))
		case int32:
			attributes = append(attributes, attribute.Int(key, int(val)))
		case int64:
			attributes = append(attributes, attribute.Int64(key, val))
		case float32:
			attributes = append(attributes, attribute.Float64(key, float64(val)))
		case float64:
			attributes = append(attributes, attribute.Float64(key, val))
		case []int:
			attributes = append(attributes, attribute.IntSlice(key, val))
		case []int64:
			attributes = append(attributes, attribute.Int64Slice(key, val))
		case []float64:
			attributes = append(attributes, attribute.Float64Slice(key, val))
		case struct{}, map[string]interface{}, map[string]string:
			attributes = append(attributes, attribute.String(key, fmt.Sprintf("%+v", val)))
		default:
			attributes = append(attributes, attribute.String(key, fmt.Sprintf("%v", val)))
		}
	}

	s.span.SetAttributes(attributes...)
}

// RecordError records an error to the current span.
func (s *Span) RecordError(err error) {
	if err == nil {
		return
	}
	s.span.RecordError(err)
	s.span.SetStatus(codes.Error, err.Error())
}

// AddEvent adds an event to the current span at current timeframe.
func (s *Span) AddEvent(events ...string) {
	for _, event := range events {
		s.span.AddEvent(event,
			trace.WithTimestamp(time.Now()),
			trace.WithStackTrace(true))
	}
}

// Context returns SpanContext.
func (s *Span) Context() tengcoruxTracer.SpanContext {
	return s.spanContext
}

// SpanContext is a wrapper around a Go context for which it stores the underlying
// context of a span. It provides convenient methods for interacting with tracing
// information. The SpanContext is designed to be used wherever tracing context
// needs to be passed or extracted within an application.
type SpanContext struct {
	ctx context.Context
}

// TraceID returns the SpanContext's TraceID. If it does not exist
// then it returns an empty string.
func (sc *SpanContext) TraceID() string {
	span := trace.SpanFromContext(sc.ctx)
	if span == nil {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

// SpanID returns the SpanContext's SpanID. If it does not exist
// then it returns an empty string.
func (sc *SpanContext) SpanID() string {
	span := trace.SpanFromContext(sc.ctx)
	if span == nil {
		return ""
	}
	return span.SpanContext().SpanID().String()
}

// Context returns the SpanContext's underlying context.
func (sc *SpanContext) Context() context.Context {
	return sc.ctx
}

// mapSpanKind maps a given span type and layer to open telemetry span kind.
func mapSpanKind(
	spanType tengcoruxTracer.SpanType,
	spanLayer tengcoruxTracer.SpanLayer,
) trace.SpanKind {
	switch spanType {
	case tengcoruxTracer.SpanTypeLocal:
		return trace.SpanKindInternal
	case tengcoruxTracer.SpanTypeEntry:
		switch spanLayer {
		case tengcoruxTracer.SpanLayerMQ:
			return trace.SpanKindConsumer
		default: // HTTP and others
			return trace.SpanKindServer
		}
	case tengcoruxTracer.SpanTypeExit:
		switch spanLayer {
		case tengcoruxTracer.SpanLayerMQ:
			return trace.SpanKindProducer
		default:
			return trace.SpanKindClient
		}
	default:
		return trace.SpanKindUnspecified
	}
}
