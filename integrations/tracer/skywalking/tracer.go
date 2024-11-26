package skywalking

import (
	"context"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
)

type Tracer struct {
	tracer   *go2sky.Tracer
	reporter go2sky.Reporter
}

func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...tengcoruxTracer.StartSpanOption) (context.Context, tengcoruxTracer.Span) {
	startSpanConfig := tengcoruxTracer.DefaultStartSpanConfig()
	for _, opt := range opts {
		opt(startSpanConfig)
	}

	go2skySpan, ctx, _ := t.tracer.CreateLocalSpan(ctx, t.generateSkywalkSpanOptions(name, startSpanConfig)...)
	go2skySpan.SetSpanLayer(mapSpanLayer(startSpanConfig.SpanLayer))
	go2skySpan.SetComponent(mapComponentLibrary(startSpanConfig.SpanLayer).AsInt32())

	return ctx, &Span{
		tracer: t,
		span:   go2skySpan,
		context: &SpanContext{
			ctx: ctx,
		},
		name: name,
	}
}

func (t *Tracer) Shutdown(_ context.Context) error {
	t.reporter.Close()
	return nil
}

func (t *Tracer) SpanFromContext(ctx context.Context) tengcoruxTracer.Span {
	go2skySpan := go2sky.ActiveSpan(ctx)
	return &Span{
		tracer: t,
		span:   go2skySpan,
		context: &SpanContext{
			ctx: ctx,
		},
		name: go2skySpan.GetOperationName(),
	}
}

////////////// Tracer's PRIVATE METHODS //////////////////

// generateSkywalkSpanOptions generates a slice of go2sky SpanOptions from a given operation name and start span config.
func (t *Tracer) generateSkywalkSpanOptions(operationName string, startSpanConfig *tengcoruxTracer.StartSpanConfig) []go2sky.SpanOption {
	options := []go2sky.SpanOption{
		go2sky.WithOperationName(operationName),
		go2sky.WithSpanType(mapSpanType(startSpanConfig.SpanType)),
	}

	if startSpanConfig.TraceID != "" {
		options = append(options, go2sky.WithContext(&propagation.SpanContext{
			TraceID:      startSpanConfig.TraceID,
			ParentSpanID: stringToSpanID(startSpanConfig.ParentSpanID),
		}))
	}

	return options
}
