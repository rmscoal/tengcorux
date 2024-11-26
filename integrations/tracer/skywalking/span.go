package skywalking

import (
	"context"
	"fmt"
	"go/types"
	"strconv"
	"time"

	"github.com/SkyAPM/go2sky"
	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
)

var _ tengcoruxTracer.Span = (*Span)(nil)

// Span represents a unit of work in a distributed trace.
//
// A span captures timing, metadata, and contextual information
// about a specific operation within a trace. This struct provides
// abstraction over the underlying SkyWalking `go2sky.Span`,
// enabling more intuitive interaction and additional functionality.
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

// ComponentLibrary represents a specific library or component where
// the current span is running. This type helps identify the origin of
// the span within SkyWalking's tracing UI, making it easier to
// debug and understand distributed traces.
//
// Each library or component is associated with a unique identifier
// that is displayed in the SkyWalking UI. These identifiers allow
// developers to distinguish between various technologies in use
// within their traced application.
//
// A complete list of available component libraries and their IDs can
// be found in the SkyWalking repository:
// https://github.com/apache/skywalking/blob/master/oap-server/server-starter/src/main/resources/component-libraries.yml
type ComponentLibrary int32

// Predefined constants representing commonly used libraries/components.
// These values correspond to the IDs specified in SkyWalking's component library file.

const (
	Unknown      ComponentLibrary = 0
	GoRedis      ComponentLibrary = 7
	PostgreSQL   ComponentLibrary = 22
	GoKafka      ComponentLibrary = 27
	RabbitMQ     ComponentLibrary = 51
	GoHttpServer ComponentLibrary = 5004
	GoMysql      ComponentLibrary = 5012
)

func (c ComponentLibrary) AsInt32() int32 {
	return int32(c)
}
