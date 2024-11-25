package skywalking

import (
	"context"
	"github.com/SkyAPM/go2sky"
	v3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	"testing"

	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
)

func TestSkyWalkingTracer_StartSpan(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, err := startTestingTracer()
	if err != nil {
		t.Error("should not throw error")
		t.FailNow()
	} else if tracer == nil {
		t.Error("tracer should not be nil")
		t.FailNow()
	}

	t.Run("WithoutOptions", func(t *testing.T) {
		ctx, span := tracer.StartSpan(context.Background(), "testing")
		if ctx == nil {
			t.Error("ctx should not be nil")
		} else if span == nil {
			t.Error("span should not be nil")
		} else if _, ok := span.(*Span); !ok {
			t.Error("span should be of type Span")
		}
	})

	t.Run("WithOptions", func(t *testing.T) {
		ctx, span := tracer.StartSpan(context.Background(), "testing",
			tengcoruxTracer.WithSpanLayer(tengcoruxTracer.SpanLayerDatabase),
			tengcoruxTracer.WithSpanType(tengcoruxTracer.SpanTypeExit),
		)
		if ctx == nil {
			t.Error("ctx should not be nil")
		} else if span == nil {
			t.Error("span should not be nil")
		} else if _, ok := span.(*Span); !ok {
			t.Error("span should be of type Span")
		}

		go2skySpan, ok := span.(*Span)
		if !ok {
			t.Error("span is not of type *Span")
			t.FailNow()
		}

		if !go2skySpan.span.IsValid() {
			t.Error("the go2sky span is not valid")
		}

		if !go2skySpan.span.IsExit() {
			t.Error("the go2sky span is not an exit span")
		}

		if go2skySpan.span.GetOperationName() != "testing" {
			t.Error("the go2sky span operation name is not testing")
		}

		reportedSpan, ok := go2skySpan.span.(go2sky.ReportedSpan)
		if !ok {
			t.Error("the go2sky span is not a reported span")
		}

		if spanType := reportedSpan.SpanType(); spanType != v3.SpanType_Exit {
			t.Errorf("the go2sky span type expects %v, but got %v", v3.SpanType_Exit, spanType)
		}

		if layer := reportedSpan.SpanLayer(); layer != v3.SpanLayer_Database {
			t.Errorf("the reported span layer expects %v, but got %v", v3.SpanLayer_Database, layer)
		}
	})
}

func TestSkyWalkingTracer_Shutdown(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, err := startTestingTracer()
	if err != nil {
		t.Error("should not throw error")
		t.FailNow()
	} else if tracer == nil {
		t.Error("tracer should not be nil")
		t.FailNow()
	}

	err = tracer.Shutdown(context.Background())
	if err != nil {
		t.Error("should not throw error")
	}
}

func TestSkyWalkingTracer_SpanFromContext(t *testing.T) {
	defer t.Cleanup(recoverPanic(t))

	tracer, _ := startTestingTracer()
	ctx, firstSpan := tracer.StartSpan(context.Background(), "testing",
		tengcoruxTracer.WithSpanLayer(tengcoruxTracer.SpanLayerDatabase),
		tengcoruxTracer.WithSpanType(tengcoruxTracer.SpanTypeExit),
	)
	secondSpan := tracer.SpanFromContext(ctx)

	firstReportedSpan, ok := firstSpan.(*Span).span.(go2sky.ReportedSpan)
	if !ok {
		t.Error("firstReportedSpan should be of type ReportedSpan")
	}

	secondReportedSpan, ok := secondSpan.(*Span).span.(go2sky.ReportedSpan)
	if !ok {
		t.Error("secondReportedSpan should be of type ReportedSpan")
	}

	if firstReportedSpan.OperationName() != secondReportedSpan.OperationName() {
		t.Error("firstReportedSpan should have the same operation name with the secondReportedSpan")
	}

	if firstReportedSpan.SpanLayer() != secondReportedSpan.SpanLayer() {
		t.Error("firstReportedSpan should have the same span layer with the secondReportedSpan")
	}

	if firstReportedSpan.SpanType() != secondReportedSpan.SpanType() {
		t.Error("firstReportedSpan should have the same span type with the secondReportedSpan")
	}
}

////////////// Testing Tracer's PRIVATE METHODS //////////////////

func TestMapSpanType(t *testing.T) {
	tracer := &Tracer{}

	if go2skySpanType := tracer.mapSpanType(tengcoruxTracer.SpanTypeLocal); go2skySpanType != go2sky.SpanTypeLocal {
		t.Errorf("expects %v but got %v", go2sky.SpanTypeLocal, go2skySpanType)
	}

	if go2skySpanType := tracer.mapSpanType(tengcoruxTracer.SpanTypeEntry); go2skySpanType != go2sky.SpanTypeEntry {
		t.Errorf("expects %v but got %v", go2sky.SpanTypeEntry, go2skySpanType)
	}

	if go2skySpanType := tracer.mapSpanType(tengcoruxTracer.SpanTypeExit); go2skySpanType != go2sky.SpanTypeExit {
		t.Errorf("expects %v but got %v", go2sky.SpanTypeExit, go2skySpanType)
	}
}

func TestMapSpanLayer(t *testing.T) {
	tracer := &Tracer{}

	if go2skySpanLayer := tracer.mapSpanLayer(tengcoruxTracer.SpanLayerUnknown); go2skySpanLayer != v3.SpanLayer_Unknown {
		t.Errorf("expects %v but got %v", v3.SpanLayer_Unknown, go2skySpanLayer)
	}

	if go2skySpanLayer := tracer.mapSpanLayer(tengcoruxTracer.SpanLayerDatabase); go2skySpanLayer != v3.SpanLayer_Database {
		t.Errorf("expects %v but got %v", v3.SpanLayer_Database, go2skySpanLayer)
	}

	if go2skySpanLayer := tracer.mapSpanLayer(tengcoruxTracer.SpanLayerHttp); go2skySpanLayer != v3.SpanLayer_Http {
		t.Errorf("expects %v but got %v", v3.SpanLayer_Http, go2skySpanLayer)
	}

	if go2skySpanLayer := tracer.mapSpanLayer(tengcoruxTracer.SpanLayerMQ); go2skySpanLayer != v3.SpanLayer_MQ {
		t.Errorf("expects %v but got %v", v3.SpanLayer_MQ, go2skySpanLayer)
	}
}
