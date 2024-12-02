package rest

type Option func(r *Rest)

// WithTracerEnabled marks tracer flag as true and registers request and response
// middleware to start and end a span.
func WithTracerEnabled() Option {
	return func(r *Rest) {
		r.tracerEnabled = true
	}
}
