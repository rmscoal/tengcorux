package tracing

import "github.com/rmscoal/tengcorux/tracer"

type Option func(*tracing)

func WithSpanNameGenerator(f SpanNameGenerator) Option {
	return func(t *tracing) {
		if f != nil {
			t.spanNameGenerator = f
		}
	}
}

// WithTracer registers the given tracer to the tracing instance.
func WithTracer(tracer tracer.Tracer) Option {
	return func(t *tracing) {
		t.provider = tracer
	}
}

func WithSQLVariables() Option {
	return func(t *tracing) {
		t.showSQLVariable = true
	}
}
