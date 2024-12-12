package opentelemetry

import (
	"testing"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

func TestOpentelemetry_NewTracer(t *testing.T) {
	t.Run("WithEnvironment", func(t *testing.T) {
		tracer := NewTracer("some_service_name", WithEnvironment("STAGING"))
		if tracer.serviceName != "some_service_name:STAGING" {
			t.Errorf("expected service name to be 'some_service_name:STAGING', got %s", tracer.serviceName)
		}
	})
	t.Run("WithExporter", func(t *testing.T) {
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatal(err)
		}
		tracer := NewTracer("some_service_name", WithExporter(exporter))
		if len(tracer.shutdowns) == 0 {
			t.Error("expected at least one shutdown")
		}
	})
}
