package tracing

import (
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"net"
	"strconv"
)

type Option func(*Tracing)

// WithAttributes adds given attributes to the span later on.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(tr *Tracing) {
		tr.spanAttributes = append(tr.spanAttributes, attrs...)
	}
}

// WithConnectionString add connection string attribute.
func WithConnectionString(str string) Option {
	return func(tr *Tracing) {
		if str != "" && tr.includeAddress {
			tr.spanAttributes = append(tr.spanAttributes,
				attribute.KeyValuePair("db.connection_string", str))
		}
	}
}

// WithClientType determines the redis client type. Choose from
// "client", "cluster", or "ring".
func WithClientType(t string) Option {
	return func(tr *Tracing) {
		switch t {
		case "client", "cluster", "ring":
			tr.spanAttributes = append(tr.spanAttributes,
				attribute.KeyValuePair("redis.client_type", t))
		}
	}
}

// WithServerAddress adds redis server address to attributes.
func WithServerAddress(addr string) Option {
	return func(tr *Tracing) {
		if addr == "" {
			return
		} else if !tr.includeAddress {
			return
		}

		host, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			return
		}
		tr.spanAttributes = append(tr.spanAttributes,
			attribute.KeyValuePair("server.host", host))

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return
		}
		tr.spanAttributes = append(tr.spanAttributes,
			attribute.KeyValuePair("server.port", port))
	}
}

// IncludeAddress determines to include redis address in span attributes.
func IncludeAddress(on bool) Option {
	return func(tr *Tracing) {
		tr.includeAddress = on
		if !on {
			// Find spans that has server.host or server.port as key and
			// remove them
			for i := 0; i < len(tr.spanAttributes); i++ {
				switch tr.spanAttributes[i].Key {
				case "server.port", "server.host", "db.connection_string":
					tr.spanAttributes = append(tr.spanAttributes[:i], tr.spanAttributes[i+1:]...)
					i--
				}
			}
		}
	}
}
