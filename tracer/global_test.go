package tracer

import "testing"

func TestGlobalTracer(t *testing.T) {
	SetGlobalTracer(&NoopTracer{})

	_, ok := GetGlobalTracer().(*NoopTracer)
	if !ok {
		t.Error("tracer is not of type *NoopTracer")
	}
}

func TestSetGlobalTracer(t *testing.T) {
	SetGlobalTracer(new(NoopTracer))
	_, ok := globalTracer.(*NoopTracer)
	if !ok {
		t.Error("global tracer is not of type *NoopTracer")
	}
}
