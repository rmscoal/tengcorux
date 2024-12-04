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

// WithEnvironment exports the environment settings later in the trace.
func WithEnvironment(env string) Option {
	return func(tracer *Tracer) {
		tracer.environment = env
		tracer.serviceName = fmt.Sprintf("%s:%s", tracer.serviceName, env)
	}
}
