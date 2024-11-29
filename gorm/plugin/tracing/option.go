package tracing

import "github.com/rmscoal/tengcorux/tracer"

type Option func(*tracing)

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
