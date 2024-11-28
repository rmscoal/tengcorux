package tracetest

import (
	"sync"
	"testing"
	"time"
)

func TestSpanRecorder_OnStart(t *testing.T) {
	recorder := NewSpanRecorder()

	span := &Span{
		StartTime: time.Now(),
		Name:      "testing",
		TraceID:   1,
		SpanID:    2,
	}

	recorder.OnStart(span)
	started := recorder.StartedSpans()
	if len(started) != 1 {
		t.Fatalf("got %d spans, want 1", len(started))
	}
}

func TestSpanRecorder_OnEnd(t *testing.T) {
	recorder := NewSpanRecorder()

	span := &Span{
		StartTime: time.Now().Add(-1 * time.Second),
		EndTime:   time.Now(),
		Name:      "testing",
		TraceID:   1,
		SpanID:    2,
	}

	recorder.OnEnd(span)
	ended := recorder.EndedSpans()
	if len(ended) != 1 {
		t.Fatalf("got %d spans, want 1", len(ended))
	}
}

func TestSpanRecorder_OnStart_Concurrent(t *testing.T) {
	recorder := NewSpanRecorder()

	firstSpan := &Span{
		StartTime: time.Now(),
		Name:      "testing_first_span",
		TraceID:   1,
		SpanID:    1,
	}
	secondSpan := &Span{
		StartTime: time.Now(),
		Name:      "testing_second_span",
		TraceID:   2,
		SpanID:    2,
	}

	// Launch 2 goroutines to add the span
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		recorder.OnStart(firstSpan)
	}(wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		recorder.OnStart(secondSpan)
	}(wg)

	wg.Wait()

	started := recorder.StartedSpans()
	if len(started) != 2 {
		t.Fatalf("got %d spans, want 2", len(started))
	}

	// Checks for each span whether information is correct
	for _, span := range started {
		switch span.SpanID {
		case 1:
			if span.Name != "testing_first_span" {
				t.Fatalf("span of id %d has invalid name, got %s but want %s", span.SpanID, span.Name, "testing_first_span")
			}
		case 2:
			if span.Name != "testing_second_span" {
				t.Fatalf("span of id %d has invalid name, got %s but want %s", span.SpanID, span.Name, "testing_second_span")
			}
		}
	}
}

func TestSpanRecorder_OnEnd_Concurrent(t *testing.T) {
	recorder := NewSpanRecorder()

	firstSpan := &Span{
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Name:      "testing_first_span",
		TraceID:   1,
		SpanID:    1,
	}
	secondSpan := &Span{
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Name:      "testing_second_span",
		TraceID:   2,
		SpanID:    2,
	}

	// Launch 2 goroutines to add the span
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		recorder.OnEnd(firstSpan)
	}(wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		recorder.OnEnd(secondSpan)
	}(wg)

	wg.Wait()

	ended := recorder.EndedSpans()
	if len(ended) != 2 {
		t.Fatalf("got %d spans, want 2", len(ended))
	}

	// Checks for each span whether information is correct
	for _, span := range ended {
		switch span.SpanID {
		case 1:
			if span.Name != "testing_first_span" {
				t.Fatalf("span of id %d has invalid name, got %s but want %s", span.SpanID, span.Name, "testing_first_span")
			}
		case 2:
			if span.Name != "testing_second_span" {
				t.Fatalf("span of id %d has invalid name, got %s but want %s", span.SpanID, span.Name, "testing_second_span")
			}
		}
	}
}
