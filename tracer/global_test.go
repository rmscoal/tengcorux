package tracer

import (
	"context"
	"errors"
	"testing"
)

func TestGlobalTracer(t *testing.T) {
	SetGlobalTracer(&NoopTracer{})

	_, ok := GetGlobalTracer().(*NoopTracer)
	if !ok {
		t.Error("tracer is not of type *NoopTracer")
	}
}

type errorTracer struct {
	*NoopTracer
}

func (*errorTracer) Shutdown(_ context.Context) error {
	return errors.New("always error")
}

func TestSetGlobalTracer(t *testing.T) {
	t.Run("Previously Shutdown Successful", func(t *testing.T) {
		SetGlobalTracer(new(NoopTracer))
		_, ok := globalTracer.(*NoopTracer)
		if !ok {
			t.Error("global tracer is not of type *NoopTracer")
		}
	})
	t.Run("Previous Shutdown Return Error", func(t *testing.T) {
		defer func() {
			globalTracer = new(NoopTracer)
		}()

		globalTracer = &errorTracer{}
		SetGlobalTracer(new(NoopTracer))
		_, ok := globalTracer.(*errorTracer)
		if !ok {
			t.Error("it should be still of type *errorTracer")
		}
	})
}
