package tracer

import (
	"context"
	"github.com/rmscoal/tengcorux/tracer/attribute"
)

type Span interface {
	// End ends the current span journey.
	End()

	// SetAttributes sets attributes to the current span.
	SetAttributes(attributes ...attribute.KeyValue)

	// RecordError sets a new error to the span.
	RecordError(err error)

	// AddEvent adds certain event into the span.
	AddEvent(descriptions ...string)

	// Context returns the SpanContext of the current span.
	Context() SpanContext
}

type SpanContext interface {
	// TraceID return the tracer ID of the SpanContext as string.
	TraceID() string

	// SpanID return the span ID of the SpanContext as string.
	SpanID() string

	// Context returns the Go's context.
	Context() context.Context
}

// SpanType determines what type is the current span.
type SpanType int32

const (
	// SpanTypeLocal determines that the current span is a local process.
	// For example invoking local function.
	SpanTypeLocal SpanType = 0
	// SpanTypeEntry determines that the current span is an entry span.
	// For example http server receiving http request.
	SpanTypeEntry SpanType = 1
	// SpanTypeExit determines that the current span is an exit span.
	// For example http client making http calls.
	SpanTypeExit SpanType = 2
)

// SpanLayer determines on what medium is the process of the span running in.
type SpanLayer int32

const (
	// SpanLayerUnknown is an unknown layer and this could be anything. It should the default value.
	SpanLayerUnknown SpanLayer = 0
	// SpanLayerDatabase is a database layer and determines running a database operations from the client.
	SpanLayerDatabase SpanLayer = 1
	// SpanLayerHttp is a http layer and determines running a http process.
	SpanLayerHttp SpanLayer = 2
	// SpanLayerMQ is a MQ layer and determines running a MQ.
	SpanLayerMQ SpanLayer = 3
)

type StartSpanConfig struct {
	// TraceID to manually change the implementation's trace id of the current span.
	TraceID string

	// ParentSpanID to manually change the implementation's span id of the current span.
	ParentSpanID string

	// SpanType should default to SpanTypeLocal.
	SpanType SpanType

	// SpanLayer should default to SpanLayerUnknown.
	SpanLayer SpanLayer
}

// StartSpanOption provides options to inject to the span when
// creating a new span.
type StartSpanOption func(*StartSpanConfig)

// WithTraceID configures the trace id of the span.
func WithTraceID(traceId string) StartSpanOption {
	return func(cfg *StartSpanConfig) {
		cfg.TraceID = traceId
	}
}

// WithParentSpanID configures the id for the span.
func WithParentSpanID(spanId string) StartSpanOption {
	return func(cfg *StartSpanConfig) {
		cfg.ParentSpanID = spanId
	}
}

// WithSpanLayer configures the SpanLayer.
func WithSpanLayer(layer SpanLayer) StartSpanOption {
	return func(cfg *StartSpanConfig) {
		cfg.SpanLayer = layer
	}
}

// WithSpanType configures the type SpanType.
func WithSpanType(t SpanType) StartSpanOption {
	return func(cfg *StartSpanConfig) {
		cfg.SpanType = t
	}
}

// DefaultStartSpanConfig returns the default StartSpanConfig.
func DefaultStartSpanConfig() *StartSpanConfig {
	return &StartSpanConfig{
		SpanType:  SpanTypeLocal,
		SpanLayer: SpanLayerUnknown,
	}
}
