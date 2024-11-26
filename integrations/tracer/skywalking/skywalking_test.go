package skywalking

import (
	"context"
	"os"
	"testing"

	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"

	"github.com/SkyAPM/go2sky"
)

var (
	exportAddress         = "127.0.0.1:11800"
	serviceName           = "testing_service"
	skywalkTestingOptions = []go2sky.TracerOption{go2sky.WithInstance("testing")}

	startTestingTracer = func() (*Tracer, error) {
		envAddress := os.Getenv("TRACER_EXPORTER_ADDRESS")
		if envAddress != "" {
			exportAddress = envAddress
		}

		return NewTracer(exportAddress, serviceName, skywalkTestingOptions...)
	}
)

func TestSkyWalking_New(t *testing.T) {
	tracer, err := startTestingTracer()
	if err != nil {
		t.Error("should not throw error")
	} else if tracer == nil {
		t.Error("tracer should not be nil")
	}

	tengcoruxTracer.SetGlobalTracer(tracer)
}

func TestSkyWalking_EndToEnd(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, err := startTestingTracer()
	if err != nil {
		t.Fatal("unable to start testing tracer")
	} else if tracer == nil {
		t.Fatal("tracer should not be nil")
	}

	ctx, span := tracer.StartSpan(context.Background(), "GET /hello",
		tengcoruxTracer.WithSpanType(tengcoruxTracer.SpanTypeEntry),
		tengcoruxTracer.WithSpanLayer(tengcoruxTracer.SpanLayerHttp),
	)
	defer span.End()

	ctx1, span1 := tracer.StartSpan(ctx, "DO_SOMETHING_ONE",
		tengcoruxTracer.WithSpanType(tengcoruxTracer.SpanTypeLocal),
		tengcoruxTracer.WithSpanLayer(tengcoruxTracer.SpanLayerUnknown),
	)
	defer span1.End()

	_, span2 := tracer.StartSpan(ctx1, "INSERT",
		tengcoruxTracer.WithSpanType(tengcoruxTracer.SpanTypeLocal),
		tengcoruxTracer.WithSpanLayer(tengcoruxTracer.SpanLayerDatabase),
	)
	defer span2.End()
}
