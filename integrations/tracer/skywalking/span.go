package skywalking

import (
	"context"
	"fmt"
	"github.com/SkyAPM/go2sky"
	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"go/types"
	"strconv"
	"time"
)

var _ tengcoruxTracer.Span = (*Span)(nil)

// Span represents a span in a trace.
type Span struct {
	span    go2sky.Span
	tracer  *Tracer
	context *SpanContext
	name    string
}

// End ends the current span.
func (s *Span) End() {
	s.span.End()
}

// SetAttributes sets attributes to the current span.
func (s *Span) SetAttributes(attributes ...attribute.KeyValue) {
	for _, attr := range attributes {
		var value string

		switch attr.Value.(type) {
		case int, int32, int64, uint, uint32, uint64, float32, float64:
			value = fmt.Sprintf("%d", attr.Value)
		case string:
			value = attr.Value.(string)
		case struct{}, types.Slice:
			value = fmt.Sprintf("%v", attr.Value)
		}

		s.span.Tag(go2sky.Tag(attr.Key), value)
	}
}

// RecordError records an error to the current span at current timeframe.
func (s *Span) RecordError(err error) {
	s.span.Error(time.Now(), err.Error())
}

// AddEvent adds an event to the current span at current timeframe.
func (s *Span) AddEvent(descriptions ...string) {
	s.span.Log(time.Now(), descriptions...)
}

// Context returns SpanContext.
func (s *Span) Context() tengcoruxTracer.SpanContext {
	return s.context
}

// SpanContext stores the underlying context of the current span.
type SpanContext struct {
	ctx context.Context
}

// TraceID returns the SpanContext's TraceID. If it does not exist
// then it returns an empty string.
func (sc *SpanContext) TraceID() string {
	traceId := go2sky.TraceID(sc.ctx)
	if traceId == "N/A" {
		traceId = ""
	}

	return traceId
}

// SpanID returns the SpanContext's SpanID. If it does not exist
// then it returns an empty string.
func (sc *SpanContext) SpanID() string {
	span, ok := go2sky.ActiveSpan(sc.ctx).(go2sky.ReportedSpan)
	if !ok {
		return ""
	}
	return strconv.Itoa(int(span.Context().SpanID))
}

// Context returns the SpanContext's underlying context.
func (sc *SpanContext) Context() context.Context {
	return sc.ctx
}
