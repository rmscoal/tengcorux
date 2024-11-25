package skywalking

import (
	"context"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	v3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
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

	go2skySpan, ctx, _ := t.tracer.CreateLocalSpan(ctx,
		go2sky.WithOperationName(name),
		go2sky.WithSpanType(t.mapSpanType(startSpanConfig.SpanType)),
		go2sky.WithContext(
			&propagation.SpanContext{
				TraceID:      startSpanConfig.TraceID,
				ParentSpanID: stringToSpanID(startSpanConfig.ParentSpanID),
			},
		),
	)
	go2skySpan.SetSpanLayer(t.mapSpanLayer(startSpanConfig.SpanLayer))

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

func (t *Tracer) mapSpanType(option tengcoruxTracer.SpanType) go2sky.SpanType {
	switch option {
	case tengcoruxTracer.SpanTypeLocal:
		return go2sky.SpanTypeLocal
	case tengcoruxTracer.SpanTypeEntry:
		return go2sky.SpanTypeEntry
	case tengcoruxTracer.SpanTypeExit:
		return go2sky.SpanTypeExit
	default:
		return go2sky.SpanTypeLocal
	}
}

func (t *Tracer) mapSpanLayer(option tengcoruxTracer.SpanLayer) v3.SpanLayer {
	switch option {
	case tengcoruxTracer.SpanLayerUnknown:
		return v3.SpanLayer_Unknown
	case tengcoruxTracer.SpanLayerDatabase:
		return v3.SpanLayer_Database
	case tengcoruxTracer.SpanLayerHttp:
		return v3.SpanLayer_Http
	case tengcoruxTracer.SpanLayerMQ:
		return v3.SpanLayer_MQ
	default:
		return v3.SpanLayer_Unknown
	}
}
