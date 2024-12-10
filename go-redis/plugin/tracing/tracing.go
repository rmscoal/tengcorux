package tracing

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/redis/go-redis/extra/rediscmd/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
)

func InstrumentTracing(rd redis.UniversalClient, opts ...Option) error {
	switch rd := rd.(type) {
	case *redis.Client:
		opt := rd.Options()
		connString := formatDBConnString(opt.Network, opt.Addr)
		rd.AddHook(newHook(connString, opts...))
		return nil
	case *redis.ClusterClient:
		rd.AddHook(newHook("", opts...))
		rd.OnNewNode(func(rdb *redis.Client) {
			opt := rdb.Options()
			connString := formatDBConnString(opt.Network, opt.Addr)
			rdb.AddHook(newHook(connString, opts...))
		})
		return nil
	case *redis.Ring:
		rd.AddHook(newHook("", opts...))
		rd.OnNewNode(func(rdb *redis.Client) {
			opt := rdb.Options()
			connString := formatDBConnString(opt.Network, opt.Addr)
			rdb.AddHook(newHook(connString, opts...))
		})
		return nil
	default:
		return fmt.Errorf("tracing: unsupported redis client type: %T", rd)
	}
}

var _defaultAttributes = []attribute.KeyValue{attribute.DBSystem("redis")}

type tracingHook struct {
	spanAttrs        []attribute.KeyValue
	spanStartOptions []tracer.StartSpanOption

	showConnString bool
	withVariable   bool
}

func newHook(connString string, opts ...Option) *tracingHook {
	hook := &tracingHook{
		spanAttrs: _defaultAttributes,
		spanStartOptions: []tracer.StartSpanOption{
			tracer.WithSpanLayer(tracer.SpanLayerDatabase),
			tracer.WithSpanType(tracer.SpanTypeExit),
		},
	}

	for _, opt := range opts {
		opt(hook)
	}

	if connString != "" && hook.showConnString {
		hook.spanAttrs = append(hook.spanAttrs,
			attribute.KeyValuePair("db.connection_string", connString))
	}

	return hook
}

func (th *tracingHook) DialHook(hook redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		ctx, span := tracer.StartSpan(ctx, "redis.dial", th.spanStartOptions...)
		defer span.End()

		span.SetAttributes(th.spanAttrs...)

		conn, err := hook(ctx, network, addr)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		return conn, nil
	}
}

func (th *tracingHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		fn, file, line := funcFileLine("github.com/redis/go-redis")

		attrs := make([]attribute.KeyValue, 0, 8)
		attrs = append(attrs,
			attribute.KeyValuePair("code.function", fn),
			attribute.KeyValuePair("code.filepath", file),
			attribute.KeyValuePair("code.lineno", line),
		)
		attrs = append(attrs, th.spanAttrs...)
		attrs = append(attrs,
			attribute.DBOperation(cmd.Name()),
			attribute.DBStatement(rediscmd.CmdString(cmd)))

		ctx, span := tracer.StartSpan(ctx, "redis."+cmd.FullName(),
			th.spanStartOptions...)
		defer span.End()

		span.SetAttributes(attrs...)

		if err := hook(ctx, cmd); err != nil {
			span.RecordError(err)
			return err
		}

		return nil
	}
}

func (th *tracingHook) ProcessPipelineHook(hook redis.ProcessPipelineHook,
) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		fn, file, line := funcFileLine("github.com/redis/go-redis")

		attrs := make([]attribute.KeyValue, 0, 8)
		attrs = append(attrs,
			attribute.KeyValuePair("code.function", fn),
			attribute.KeyValuePair("code.filepath", file),
			attribute.KeyValuePair("code.lineno", line),
		)
		attrs = append(attrs, th.spanAttrs...)

		summary, cmdsString := rediscmd.CmdsString(cmds)
		attrs = append(attrs, attribute.DBStatement(cmdsString))

		ctx, span := tracer.StartSpan(ctx, "redis.pipeline->"+
			strings.ReplaceAll(summary, " ", "->"),
			th.spanStartOptions...)
		defer span.End()

		span.SetAttributes(attrs...)

		if err := hook(ctx, cmds); err != nil {
			span.RecordError(err)
			return err
		}

		return nil
	}
}
