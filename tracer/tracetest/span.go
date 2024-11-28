package tracetest

import (
	"time"

	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
)

type Span struct {
	StartTime    time.Time
	EndTime      time.Time
	Events       []string
	Attributes   []attribute.KeyValue
	Name         string
	TraceID      uint64
	SpanID       uint64
	ParentSpanID uint64
	Layer        tengcoruxTracer.SpanLayer
	Type         tengcoruxTracer.SpanType
	Error        error

	tracer *Tracer
}

// End ends the span by marking the EndTime as now as well as
// appending the current span to the ended list of span by the
// SpanRecorder.
func (s *Span) End() {
	s.EndTime = time.Now()
	s.tracer.recorder.OnEnd(s)
}

// SetAttributes appends the given attributes into the Attributes slice.
func (s *Span) SetAttributes(kv ...attribute.KeyValue) {
	s.Attributes = append(s.Attributes, kv...)
}

// RecordError marks that the current test span has an error.
func (s *Span) RecordError(err error) {
	s.Error = err
}

// AddEvent appends the string of events into the Events slice.
func (s *Span) AddEvent(events ...string) {
	s.Events = append(s.Events, events...)
}

// Context returns nil.
func (s *Span) Context() tengcoruxTracer.SpanContext {
	return nil
}

// ReadWriteSpan allows the span to be read and written.
type ReadWriteSpan Span

// ReadOnlySpan only allows the span to be read.
type ReadOnlySpan Span

type prevSpanContextKey struct{}

// prevSpanKey is the key that holds the *Span value inside a context
// carried along during the transaction.
var prevSpanKey prevSpanContextKey
