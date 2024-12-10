package tracing

import "github.com/rmscoal/tengcorux/tracer/attribute"

type Option func(*tracingHook)

// WithConnectionString allows to show the conn string in the span attribute.
func WithConnectionString(on bool) Option {
	return func(hook *tracingHook) {
		hook.showConnString = on
	}
}

// WithAttributes adds given attributes to the span later on.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(hook *tracingHook) {
		hook.spanAttrs = append(hook.spanAttrs, attrs...)
	}
}
