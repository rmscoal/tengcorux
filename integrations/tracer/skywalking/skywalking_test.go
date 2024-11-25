package skywalking

import (
	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	"testing"

	"github.com/SkyAPM/go2sky"
)

var (
	exportAddress         = "127.0.0.1:11800"
	serviceName           = "testing_service"
	skywalkTestingOptions = []go2sky.TracerOption{go2sky.WithInstance("testing")}

	startTestingTracer = func() (*Tracer, error) {
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
