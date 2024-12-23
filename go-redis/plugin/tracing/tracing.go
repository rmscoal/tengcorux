package tracing

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/redis/go-redis/extra/rediscmd/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
)

const goRedisPkgName = "github.com/redis/go-redis"

type Tracing struct {
	spanAttributes    []attribute.KeyValue
	spanStartOptions  []tracer.StartSpanOption
	spanNameGenerator SpanNameGenerator

	includeAddress bool
}

func (tr *Tracing) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		var attrs []attribute.KeyValue
		attrs = append(attrs, tr.spanAttributes...)
		if tr.includeAddress {
			host, portStr, err := net.SplitHostPort(addr)
			if err != nil {
				goto Dial
			}
			attrs = append(attrs,
				attribute.KeyValuePair("server.address", host))

			port, err := strconv.Atoi(portStr)
			if err != nil {
				goto Dial
			}
			attrs = append(attrs, attribute.KeyValuePair("server.port", port))
		}

	Dial:
		ctx, span := tracer.StartSpan(ctx,
			tr.spanNameGenerator.Dial(network, addr),
			tr.spanStartOptions...)
		defer span.End()

		span.SetAttributes(attrs...)

		conn, err := next(ctx, network, addr)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		return conn, nil
	}
}

func (tr *Tracing) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		fn, file, line := funcFileLine(goRedisPkgName)

		var attrs []attribute.KeyValue
		attrs = append(attrs,
			attribute.KeyValuePair("code.function", fn),
			attribute.KeyValuePair("code.filepath", file),
			attribute.KeyValuePair("code.lineno", line))
		attrs = append(attrs, tr.spanAttributes...)
		attrs = append(attrs,
			attribute.DBOperation(cmd.Name()),
			attribute.DBStatement(rediscmd.CmdString(cmd)))

		ctx, span := tracer.StartSpan(ctx,
			tr.spanNameGenerator.Process(cmd),
			tr.spanStartOptions...)
		defer span.End()

		span.SetAttributes(attrs...)

		if err := next(ctx, cmd); err != nil {
			span.RecordError(err)
			return err
		}

		return nil
	}
}

func (tr *Tracing) ProcessPipelineHook(
	next redis.ProcessPipelineHook,
) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		fn, file, line := funcFileLine("github.com/redis/go-redis")

		var attrs []attribute.KeyValue
		attrs = append(attrs, tr.spanAttributes...)
		attrs = append(attrs,
			attribute.KeyValuePair("code.function", fn),
			attribute.KeyValuePair("code.filepath", file),
			attribute.KeyValuePair("code.lineno", line))

		_, commands := rediscmd.CmdsString(cmds)
		attrs = append(attrs, attribute.DBStatement(commands))

		ctx, span := tracer.StartSpan(ctx,
			tr.spanNameGenerator.PipelineProcess(cmds),
			tr.spanStartOptions...)
		defer span.End()

		span.SetAttributes(attrs...)

		if err := next(ctx, cmds); err != nil {
			span.RecordError(err)
			return err
		}

		return nil
	}
}

type SpanNameDialGenerator func(network, addr string) string
type SpanNameProcessGenerator func(cmd redis.Cmder) string
type SpanNamePipelineGenerator func(cmd []redis.Cmder) string

type SpanNameGenerator struct {
	Dial            SpanNameDialGenerator
	Process         SpanNameProcessGenerator
	PipelineProcess SpanNamePipelineGenerator
}

var (
	_defaultSpanNameDialGenerator = func(network, addr string) string {
		return "redis.dial"
	}
	_defaultSpanNameProcessGenerator = func(cmd redis.Cmder) string {
		return "redis." + cmd.FullName()
	}
	_defaultSpanNamePipelineGenerator = func(commands []redis.Cmder) string {
		summary, _ := rediscmd.CmdsString(commands)
		return "redis.pipeline->" +
			strings.ReplaceAll(summary, " ", "->")
	}
)

// NewHook creates a new Tracing instance.
func NewHook(opts ...Option) *Tracing {
	tr := &Tracing{
		spanAttributes: []attribute.KeyValue{attribute.DBSystem("redis")},
		spanStartOptions: []tracer.StartSpanOption{
			tracer.WithSpanLayer(tracer.SpanLayerDatabase),
			tracer.WithSpanType(tracer.SpanTypeExit),
		},
		spanNameGenerator: SpanNameGenerator{
			Dial:            _defaultSpanNameDialGenerator,
			Process:         _defaultSpanNameProcessGenerator,
			PipelineProcess: _defaultSpanNamePipelineGenerator,
		},
	}

	for _, opt := range opts {
		opt(tr)
	}

	return tr
}

// InstrumentTracing automatically registers hook given the redis instance. It
// also detects the client type such that on multi-nodes environment it is able
// to register tracing with the OnNewNode hook.
func InstrumentTracing(rd redis.UniversalClient, opts ...Option) error {
	var options []Option
	options = append(options, opts...)

	switch rd := rd.(type) {
	case *redis.Client:
		redisOption := rd.Options()
		connString := formatDBConnString(redisOption.Network, redisOption.Addr)
		options = append(options, WithClientType("client"),
			WithConnectionString(connString))

		rd.AddHook(NewHook(options...))
		return nil
	case *redis.ClusterClient:
		rd.AddHook(NewHook(options...))
		rd.OnNewNode(func(rdb *redis.Client) {
			redisOption := rdb.Options()
			connString := formatDBConnString(redisOption.Network,
				redisOption.Addr)
			nodeOptions := append([]Option(nil), options...)
			nodeOptions = append(nodeOptions, WithClientType("cluster"),
				WithConnectionString(connString))

			rdb.AddHook(NewHook(nodeOptions...))
		})
		return nil
	case *redis.Ring:
		rd.AddHook(NewHook(options...))
		rd.OnNewNode(func(rdb *redis.Client) {
			redisOption := rdb.Options()
			connString := formatDBConnString(redisOption.Network,
				redisOption.Addr)
			nodeOptions := append([]Option(nil), options...)
			nodeOptions = append(nodeOptions, WithClientType("ring"),
				WithConnectionString(connString))

			rdb.AddHook(NewHook(nodeOptions...))
		})
		return nil
	default:
		return fmt.Errorf("tracing: unsupported redis client type: %T", rd)
	}
}
