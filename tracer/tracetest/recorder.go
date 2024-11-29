package tracetest

import "sync"

// SpanRecorder is responsible for recorder the spans during the transaction
// of the tracer. Such that it can later be tested and checked what are the
// produced spans.
type SpanRecorder struct {
	startedMU sync.RWMutex
	starts    []*ReadWriteSpan

	endedMU sync.RWMutex
	ends    []*ReadOnlySpan
}

func NewSpanRecorder() *SpanRecorder { return new(SpanRecorder) }

// OnStart make the span as a ReadWriteSpan and insert to a slice of started spans.
func (sr *SpanRecorder) OnStart(s *Span) {
	sr.startedMU.Lock()
	defer sr.startedMU.Unlock()

	if s == nil {
		return
	}

	rwSpan := ReadWriteSpan(*s)
	sr.starts = append(sr.starts, &rwSpan)
}

// OnEnd make the span as a ReadOnlySpan and insert to a slice of ended spans.
func (sr *SpanRecorder) OnEnd(s *Span) {
	sr.endedMU.Lock()
	defer sr.endedMU.Unlock()

	if s == nil {
		return
	}

	roSpan := ReadOnlySpan(*s)
	sr.ends = append(sr.ends, &roSpan)
}

// StartedSpans returns a copy of the started slice spans.
func (sr *SpanRecorder) StartedSpans() []*ReadWriteSpan {
	sr.startedMU.RLock()
	defer sr.startedMU.RUnlock()
	dst := make([]*ReadWriteSpan, len(sr.starts))
	copy(dst, sr.starts)
	return dst
}

// EndedSpans returns a copy of the ended slice spans.
func (sr *SpanRecorder) EndedSpans() []*ReadOnlySpan {
	sr.endedMU.RLock()
	defer sr.endedMU.RUnlock()
	dst := make([]*ReadOnlySpan, len(sr.ends))
	copy(dst, sr.ends)
	return dst
}
