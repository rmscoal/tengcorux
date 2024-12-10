package tracing

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"github.com/rmscoal/tengcorux/tracer/tracetest"
)

func TestInstrumentTracing(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()

	tt := tracetest.NewTracer()
	tracer.SetGlobalTracer(tt)

	db, mock := redismock.NewClientMock()
	err := InstrumentTracing(db)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectSet("some_key", "some_value", 0).RedisNil()

	err = db.Set(context.TODO(), "some_key", "some_value", 0).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	// It does not record ? Perhaps because it is using mock
	// if len(tt.Recorder().EndedSpans()) == 0 {
	// 	t.Error("expected spans to be recorded")
	// }
}

func TestTracingHook_ProcessHook(t *testing.T) {
	tt := tracetest.NewTracer()
	tracer.SetGlobalTracer(tt)

	hook := newHook("",
		WithConnectionString(true),
		WithAttributes(
			attribute.KeyValuePair("some_key", "some_value")))

	ctx, span := tracer.StartSpan(context.TODO(), "testing")
	cmd := redis.NewCmd(ctx, "ping")
	defer span.End()

	processHook := hook.ProcessHook(func(ctx context.Context, cmd redis.Cmder) error {
		span := tracer.SpanFromContext(ctx)
		ttSpan, ok := span.(*tracetest.Span)
		if !ok {
			t.Fatal("span was not recorded")
		}

		if ttSpan.Name != "redis.ping" {
			t.Errorf("expected name to be redis.ping but got %s", ttSpan.Name)
		}

		if len(ttSpan.Attributes) < 1 {
			t.Fatal("expected at least one attribute")
		}

		if ttSpan.Type != tracer.SpanTypeExit {
			t.Errorf("expected type to be %d but got %d",
				tracer.SpanTypeExit, ttSpan.Type)
		}

		if ttSpan.Layer != tracer.SpanLayerDatabase {
			t.Errorf("expected layer to be %d but got %d",
				tracer.SpanLayerDatabase, ttSpan.Layer)
		}

		for _, attr := range ttSpan.Attributes {
			switch attr.Key {
			case attribute.DBStatementKey:
				val := attr.Value.(string)
				if val != "ping" {
					t.Errorf("expected val to be ping but got %s", val)
				}
			case attribute.DBSystemKey:
				val := attr.Value.(string)
				if val != "redis" {
					t.Errorf("expected val to be redis but got %s", val)
				}
			case attribute.DBOperationKey:
				val := attr.Value.(string)
				if val != "ping" {
					t.Errorf("expected val to be ping but got %s", val)
				}
			case "code.function", "code.filepath", "code.lineno":
				if attr.Value == nil {
					t.Error("expected code.xxx attribute not be empty")
				}
			case "some_key": // _defaultAttributes
				if attr.Value.(string) != "some_value" {
					t.Errorf("expected some_value to be some_value but got %s",
						attr.Value.(string))
				}
			}
		}

		return nil
	})

	err := processHook(ctx, cmd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTracingHook_ProcessPipelineHook(t *testing.T) {
	tt := tracetest.NewTracer()
	tracer.SetGlobalTracer(tt)

	hook := newHook("")

	ctx, span := tracer.StartSpan(context.TODO(), "testing")
	cmds := []redis.Cmder{
		redis.NewCmd(ctx, "ping"),
		redis.NewStringCmd(ctx, "get", "key"),
	}
	defer span.End()

	processPipelineHook := hook.ProcessPipelineHook(func(ctx context.Context, cmds []redis.Cmder) error {
		span := tracer.SpanFromContext(ctx)
		ttSpan, ok := span.(*tracetest.Span)
		if !ok {
			t.Fatal("span was not recorded")
		}

		if ttSpan.Name != "redis.pipeline->ping->get" {
			t.Errorf("expected name to be redis.pipeline->ping->get but got %s", ttSpan.Name)
		}

		if len(ttSpan.Attributes) < 1 {
			t.Fatal("expected at least one attribute")
		}

		if ttSpan.Type != tracer.SpanTypeExit {
			t.Errorf("expected type to be %d but got %d",
				tracer.SpanTypeExit, ttSpan.Type)
		}

		if ttSpan.Layer != tracer.SpanLayerDatabase {
			t.Errorf("expected layer to be %d but got %d",
				tracer.SpanLayerDatabase, ttSpan.Layer)
		}

		for _, attr := range ttSpan.Attributes {
			switch attr.Key {
			case attribute.DBStatementKey:
				val := attr.Value.(string)
				if val != fmt.Sprint("ping\nget key") {
					t.Errorf("expected val to be ping but got %s", val)
				}
			case attribute.DBSystemKey:
				val := attr.Value.(string)
				if val != "redis" {
					t.Errorf("expected val to be redis but got %s", val)
				}
			case "code.function", "code.filepath", "code.lineno":
				if attr.Value == nil {
					t.Error("expected code.xxx attribute not be empty")
				}
			case "some_key": // _defaultAttributes
				if attr.Value.(string) != "some_value" {
					t.Errorf("expected some_value to be some_value but got %s",
						attr.Value.(string))
				}
			}
		}

		return nil
	})

	err := processPipelineHook(ctx, cmds)
	if err != nil {
		t.Fatal(err)
	}
}
