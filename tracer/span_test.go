package tracer

import "testing"

func TestDefaultStartSpanConfig(t *testing.T) {
	defaultStartSpanConfig := DefaultStartSpanConfig()

	if defaultStartSpanConfig.SpanLayer != SpanLayerUnknown {
		t.Error("default start span layer should be SpanLayerUnknown")
	}

	if defaultStartSpanConfig.SpanType != SpanTypeLocal {
		t.Error("default start span type should be SpanTypeLocal")
	}

	if defaultStartSpanConfig.TraceID != "" {
		t.Error("default start span TraceID should be empty")
	}

	if defaultStartSpanConfig.ParentSpanID != "" {
		t.Error("default start span ParentSpanID should be empty")
	}
}

func TestStartSpanOption(t *testing.T) {
	t.Run("WithTraceID", func(t *testing.T) {
		spanConfig := DefaultStartSpanConfig()
		opt := WithTraceID("some_trace_id")
		opt(spanConfig)
		if spanConfig.TraceID != "some_trace_id" {
			t.Error("trace id should be `some_trace_id`")
		}
	})

	t.Run("WithParentSpanID", func(t *testing.T) {
		spanConfig := DefaultStartSpanConfig()
		opt := WithParentSpanID("some_span_id")
		opt(spanConfig)
		if spanConfig.ParentSpanID != "some_span_id" {
			t.Error("span id should be `some_span_id`")
		}
	})

	t.Run("WithSpanLayer", func(t *testing.T) {
		t.Run("WithSpanLayer_SpanLayerUnknown", func(t *testing.T) {
			spanConfig := DefaultStartSpanConfig()
			opt := WithSpanLayer(SpanLayerUnknown)
			opt(spanConfig)
			if spanConfig.SpanLayer != SpanLayerUnknown {
				t.Error("span layer should be `SpanLayerUnknown`")
			}
		})
		t.Run("WithSpanLayer_SpanLayerDatabase", func(t *testing.T) {
			spanConfig := DefaultStartSpanConfig()
			opt := WithSpanLayer(SpanLayerDatabase)
			opt(spanConfig)
			if spanConfig.SpanLayer != SpanLayerDatabase {
				t.Error("span layer should be `SpanLayerDatabase`")
			}
		})
		t.Run("WithSpanLayer_SpanLayerHttp", func(t *testing.T) {
			spanConfig := DefaultStartSpanConfig()
			opt := WithSpanLayer(SpanLayerHttp)
			opt(spanConfig)
			if spanConfig.SpanLayer != SpanLayerHttp {
				t.Error("span layer should be `SpanLayerHttp`")
			}
		})
		t.Run("WithSpanLayer_SpanLayerMQ", func(t *testing.T) {
			spanConfig := DefaultStartSpanConfig()
			opt := WithSpanLayer(SpanLayerMQ)
			opt(spanConfig)
			if spanConfig.SpanLayer != SpanLayerMQ {
				t.Error("span layer should be `SpanLayerMQ`")
			}
		})
	})

	t.Run("WithSpanType", func(t *testing.T) {
		t.Run("WithSpanType_SpanTypeLocal", func(t *testing.T) {
			spanConfig := DefaultStartSpanConfig()
			opt := WithSpanType(SpanTypeLocal)
			opt(spanConfig)
			if spanConfig.SpanType != SpanTypeLocal {
				t.Error("span type should be `SpanTypeLocal`")
			}
		})
		t.Run("WithSpanType_SpanTypeEntry", func(t *testing.T) {
			spanConfig := DefaultStartSpanConfig()
			opt := WithSpanType(SpanTypeEntry)
			opt(spanConfig)
			if spanConfig.SpanType != SpanTypeEntry {
				t.Error("span type should be `SpanTypeEntry`")
			}
		})
		t.Run("WithSpanType_SpanTypeExit", func(t *testing.T) {
			spanConfig := DefaultStartSpanConfig()
			opt := WithSpanType(SpanTypeExit)
			opt(spanConfig)
			if spanConfig.SpanType != SpanTypeExit {
				t.Error("span type should be `SpanTypeExit`")
			}
		})
	})
}
