package skywalking

import (
	"testing"

	"github.com/SkyAPM/go2sky"
	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	v3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

func TestStringToSpanID(t *testing.T) {
	tests := []struct {
		input    string
		expected int32
	}{
		{
			input:    "0",
			expected: 0,
		},
		{
			input:    "1",
			expected: 1,
		},
		{
			input:    "0.4",
			expected: 0,
		},
		{
			input:    "1s3hb",
			expected: 0,
		},
		{
			input:    "9321",
			expected: 9321,
		},
	}

	for _, test := range tests {
		id := stringToSpanID(test.input)
		if id != test.expected {
			t.Errorf("expected %d, got %d", test.expected, id)
		}
	}
}

func TestMapSpanType(t *testing.T) {
	if go2skySpanType := mapSpanType(tengcoruxTracer.SpanTypeLocal); go2skySpanType != go2sky.SpanTypeLocal {
		t.Errorf("expects %v but got %v", go2sky.SpanTypeLocal, go2skySpanType)
	}

	if go2skySpanType := mapSpanType(tengcoruxTracer.SpanTypeEntry); go2skySpanType != go2sky.SpanTypeEntry {
		t.Errorf("expects %v but got %v", go2sky.SpanTypeEntry, go2skySpanType)
	}

	if go2skySpanType := mapSpanType(tengcoruxTracer.SpanTypeExit); go2skySpanType != go2sky.SpanTypeExit {
		t.Errorf("expects %v but got %v", go2sky.SpanTypeExit, go2skySpanType)
	}
}

func TestMapSpanLayer(t *testing.T) {
	if go2skySpanLayer := mapSpanLayer(tengcoruxTracer.SpanLayerUnknown); go2skySpanLayer != v3.SpanLayer_Unknown {
		t.Errorf("expects %v but got %v", v3.SpanLayer_Unknown, go2skySpanLayer)
	}

	if go2skySpanLayer := mapSpanLayer(tengcoruxTracer.SpanLayerDatabase); go2skySpanLayer != v3.SpanLayer_Database {
		t.Errorf("expects %v but got %v", v3.SpanLayer_Database, go2skySpanLayer)
	}

	if go2skySpanLayer := mapSpanLayer(tengcoruxTracer.SpanLayerHttp); go2skySpanLayer != v3.SpanLayer_Http {
		t.Errorf("expects %v but got %v", v3.SpanLayer_Http, go2skySpanLayer)
	}

	if go2skySpanLayer := mapSpanLayer(tengcoruxTracer.SpanLayerMQ); go2skySpanLayer != v3.SpanLayer_MQ {
		t.Errorf("expects %v but got %v", v3.SpanLayer_MQ, go2skySpanLayer)
	}
}
