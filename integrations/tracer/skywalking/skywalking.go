package skywalking

import (
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

func NewTracer(exportAddr, serviceName string, opts ...go2sky.TracerOption) (*Tracer, error) {
	r, err := reporter.NewGRPCReporter(exportAddr)
	if err != nil {
		return nil, err
	}

	opts = append(opts, go2sky.WithReporter(r))
	tracer, err := go2sky.NewTracer(serviceName, opts...)
	if err != nil {
		return nil, err
	}

	return &Tracer{tracer: tracer, reporter: r}, nil
}
