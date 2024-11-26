package skywalking

import (
	"strconv"

	"github.com/SkyAPM/go2sky"
	tengcoruxTracer "github.com/rmscoal/tengcorux/tracer"
	v3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

// stringToSpanID converts a string to int32.
// If the string given is not a valid integer, it returns 0.
func stringToSpanID(s string) int32 {
	id, _ := strconv.ParseInt(s, 10, 32)
	return int32(id)
}

// mapSpanType maps a given tengcorux's SpanType to go2sky's SpanType.
func mapSpanType(option tengcoruxTracer.SpanType) go2sky.SpanType {
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

// mapSpanLayer maps a given tengcorux's SpanLayer to go2sky's SpanLayer.
func mapSpanLayer(option tengcoruxTracer.SpanLayer) v3.SpanLayer {
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

// mapComponentLibrary maps a go2sky span layer to a component library.
// TODO: Make it extensible....
func mapComponentLibrary(option tengcoruxTracer.SpanLayer) ComponentLibrary {
	switch option {
	case tengcoruxTracer.SpanLayerUnknown:
		return Unknown
	case tengcoruxTracer.SpanLayerMQ:
		return GoKafka
	case tengcoruxTracer.SpanLayerHttp:
		return GoHttpServer
	case tengcoruxTracer.SpanLayerDatabase:
		return GoMysql
	default:
		return Unknown
	}
}
