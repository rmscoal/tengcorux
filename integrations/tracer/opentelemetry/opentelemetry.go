package opentelemetry

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewTracer(serviceName string, opts ...Option) *Tracer {
	tracer := &Tracer{
		serviceName: serviceName,
	}

	for _, opt := range opts {
		opt(tracer)
	}

	// Propagation
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	// Start the tracer
	tracer.tracer = otel.Tracer(tracer.serviceName)

	return tracer
}

// Version returns the current tracer's version
func (t *Tracer) Version() string {
	return "v0.1.0"
}

type Option func(*Tracer)

// WithExporter we'll use the given exporter are the tracer provider. This enables
// easy extensibility for library that adheres to OpenTelemetry.
func WithExporter(exporter sdktrace.SpanExporter) Option {
	return func(tracer *Tracer) {
		provider := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(time.Second)),
		)
		tracer.shutdowns = append(tracer.shutdowns, provider.Shutdown)
		otel.SetTracerProvider(provider)
	}
}

// WithEnvironment exports the environment settings later in the trace
// that will be embedded into the service name with a ":" separator.
func WithEnvironment(env string) Option {
	return func(tracer *Tracer) {
		tracer.environment = env
		tracer.serviceName = fmt.Sprintf("%s:%s", tracer.serviceName, env)
	}
}
