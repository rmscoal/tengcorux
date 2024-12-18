package tracing

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"github.com/rmscoal/tengcorux/tracer/tracetest"
)

func TestNewHook_Options(t *testing.T) {
	t.Run("WithAttributes", func(t *testing.T) {
		hook := NewHook(WithAttributes(attribute.KeyValuePair("ho", "ok")))
		if len(hook.spanAttributes) != 2 {
			t.Errorf("expecting 2 span attributes field, got %d",
				len(hook.spanAttributes))
		}
	})

	t.Run("WithConnectionString", func(t *testing.T) {
		hook := NewHook(WithConnectionString("redis://localhost:6379"))
		if len(hook.spanAttributes) > 1 {
			t.Errorf("expecting 1 span attribute, got %d, it should not create new because not enabled",
				len(hook.spanAttributes))
		}

		hook = NewHook(IncludeAddress(true),
			WithConnectionString("redis://localhost:6379"))
		if len(hook.spanAttributes) != 2 {
			t.Errorf("expecting 2 span attributes field, got %d",
				len(hook.spanAttributes))
		}
	})

	t.Run("WithClientType", func(t *testing.T) {
		hook := NewHook(WithClientType("random"))
		if len(hook.spanAttributes) > 1 {
			t.Errorf("expecting 1 span attribute, got %d, it should not create new because mistaken type",
				len(hook.spanAttributes))
		}

		hook = NewHook(WithClientType("client"))
		if len(hook.spanAttributes) != 2 {
			t.Errorf("expecting 2 span attributes field, got %d, client is valid type",
				len(hook.spanAttributes))
		}

		hook = NewHook(WithClientType("cluster"))
		if len(hook.spanAttributes) != 2 {
			t.Errorf("expecting 2 span attributes field, got %d, cluster is valid type",
				len(hook.spanAttributes))
		}

		hook = NewHook(WithClientType("ring"))
		if len(hook.spanAttributes) != 2 {
			t.Errorf("expecting 2 span attributes field, got %d, ring is valid type",
				len(hook.spanAttributes))
		}
	})

	t.Run("WithServerAddress", func(t *testing.T) {
		t.Run("Disabled", func(t *testing.T) {
			hook := NewHook(WithServerAddress("localhost:6379"))
			if len(hook.spanAttributes) > 1 {
				t.Errorf("expecting 1 span attribute, got %d, it should not create new because not enabled",
					len(hook.spanAttributes))
			}
		})

		t.Run("Enabled", func(t *testing.T) {
			hook := NewHook(IncludeAddress(true),
				WithServerAddress("localhost:6379"))
			if len(hook.spanAttributes) != 3 {
				t.Errorf("expecting 3 span attributes field, got %d",
					len(hook.spanAttributes))
			}
		})

		t.Run("IncorrectAddress", func(t *testing.T) {
			hook := NewHook(IncludeAddress(true),
				WithServerAddress("redis://localhost:6379"))
			if len(hook.spanAttributes) != 1 {
				t.Errorf("expecting 1 span attributes field, got %d",
					len(hook.spanAttributes))
			}
		})
	})

	t.Run("IncludeAddress", func(t *testing.T) {
		t.Run("Disabled", func(t *testing.T) {
			hook := NewHook(WithServerAddress("redis://localhost:6379"),
				IncludeAddress(false))
			if len(hook.spanAttributes) > 1 {
				t.Errorf("expecting 1 span attribute, got %d, it should not create new because not enabled",
					len(hook.spanAttributes))
			}
		})

		t.Run("PreviouslySet", func(t *testing.T) {
			hook := NewHook(WithServerAddress("localhost:6379"),
				WithConnectionString("some_connection_string"),
				IncludeAddress(false))
			if len(hook.spanAttributes) != 1 {
				t.Errorf("expecting 1 span attributes field, got %d",
					len(hook.spanAttributes))
			}
		})
	})
}

func TestTracingHook_ProcessHook(t *testing.T) {
	tt := tracetest.NewTracer()
	tracer.SetGlobalTracer(tt)

	hook := NewHook()
	ctx := context.TODO()
	cmd := redis.NewCmd(ctx, "ping")

	processHook := hook.ProcessHook(func(
		ctx context.Context, cmd redis.Cmder,
	) error {
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

	hook := NewHook()
	ctx := context.TODO()
	commands := []redis.Cmder{
		redis.NewCmd(ctx, "ping"),
		redis.NewStringCmd(ctx, "get", "key"),
	}

	processPipelineHook := hook.ProcessPipelineHook(func(
		ctx context.Context, cmds []redis.Cmder,
	) error {
		span := tracer.SpanFromContext(ctx)
		ttSpan, ok := span.(*tracetest.Span)
		if !ok {
			t.Fatal("span was not recorded")
		}

		if ttSpan.Name != "redis.pipeline->ping->get" {
			t.Errorf("expected name to be redis.pipeline->ping->get but got %s",
				ttSpan.Name)
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
				if val != "ping\nget key" {
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
			}
		}

		return nil
	})

	err := processPipelineHook(ctx, commands)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOption_WithSpanNameGenerator_Process(t *testing.T) {
	tt := tracetest.NewTracer()
	tracer.SetGlobalTracer(tt)
	hook := NewHook(WithSpanNameGenerator(SpanNameGenerator{
		Process: func(cmd redis.Cmder) string {
			return fmt.Sprintf("hello_world_%s", cmd.FullName())
		},
	}))
	ctx := context.TODO()
	cmd := redis.NewCmd(ctx, "ping")
	processHook := hook.ProcessHook(func(
		ctx context.Context, cmd redis.Cmder,
	) error {
		span := tracer.SpanFromContext(ctx)
		ttSpan, ok := span.(*tracetest.Span)
		if !ok {
			t.Fatal("span was not recorded")
		}

		if ttSpan.Name != "hello_world_ping" {
			t.Errorf("expected name to be hello_world_ping but got %s",
				ttSpan.Name)
		}

		return nil
	})

	err := processHook(ctx, cmd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOption_WithSpanNameGenerator_PipelineProcess(t *testing.T) {
	tt := tracetest.NewTracer()
	tracer.SetGlobalTracer(tt)
	hook := NewHook(WithSpanNameGenerator(SpanNameGenerator{
		PipelineProcess: func(cmd []redis.Cmder) string {
			return "hello_world"
		},
	}))
	ctx := context.TODO()
	commands := []redis.Cmder{
		redis.NewCmd(ctx, "ping"),
		redis.NewStringCmd(ctx, "get", "key"),
	}
	pipelineHook := hook.ProcessPipelineHook(func(
		ctx context.Context, commands []redis.Cmder,
	) error {
		span := tracer.SpanFromContext(ctx)
		ttSpan, ok := span.(*tracetest.Span)
		if !ok {
			t.Fatal("span was not recorded")
		}

		if ttSpan.Name != "hello_world" {
			t.Errorf("expected name to be hello_world but got %s",
				ttSpan.Name)
		}

		return nil
	})

	err := pipelineHook(ctx, commands)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOption_WithSpanNameGenerator_Dial(t *testing.T) {
	tt := tracetest.NewTracer()
	tracer.SetGlobalTracer(tt)
	hook := NewHook(WithSpanNameGenerator(SpanNameGenerator{
		Dial: func(network, addr string) string {
			return "dialing_911"
		},
	}))
	ctx := context.TODO()
	dialHook := hook.DialHook(func(
		ctx context.Context, network, addr string,
	) (net.Conn, error) {
		span := tracer.SpanFromContext(ctx)
		ttSpan, ok := span.(*tracetest.Span)
		if !ok {
			t.Fatal("span was not recorded")
		}

		if ttSpan.Name != "dialing_911" {
			t.Errorf("expected name to be dialing_911 but got %s",
				ttSpan.Name)
		}

		return nil, nil
	})

	_, err := dialHook(ctx, "", "")
	if err != nil {
		t.Fatal(err)
	}
}
