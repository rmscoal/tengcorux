package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-resty/resty/v2"
	"github.com/rmscoal/tengcorux/reqid"
	"github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
)

// Rest is a wrapper for *resty.Client with extra fields and features.
type Rest struct {
	*resty.Client

	tracerEnabled bool
}

func New(opts ...Option) *Rest {
	rest := &Rest{
		Client:        resty.New(),
		tracerEnabled: false,
	}

	for _, opt := range opts {
		opt(rest)
	}

	if rest.tracerEnabled {
		rest.registerTracerMiddleware()
	}

	return rest
}

// registerTracerMiddleware registers OnBeforeRequest and OnAfterResponse middleware
// for tracer to start the span and captures the attributes.
func (r *Rest) registerTracerMiddleware() {
	r.Client = r.Client.
		OnBeforeRequest(
			func(client *resty.Client, request *resty.Request) error {
				// Here, we are going to start a span. We should also fill the span
				// with attributes. Therefore, before resty makes a request, we need
				// to capture the following for tracer:
				// 1. Method (like GET, POST, and others),
				// 2. URL Target, this should include path and query params,
				// 3. Body of the payload if it is a POST/PUT/PATCH request, and
				// 4. Headers of the request.
				method := request.Method
				ctx, span := tracer.StartSpan(
					request.Context(), "HTTP "+method+" Request",
					tracer.WithSpanType(tracer.SpanTypeExit),
					tracer.WithSpanLayer(tracer.SpanLayerHttp),
				)

				// Injects X-Request-Id to the request by reading from the context
				requestID := reqid.RetrieveFromContext(ctx)
				request.SetHeader(reqid.HeaderKey, requestID)

				// Add spans to the attribute
				span.SetAttributes(
					attribute.HTTPRequestID(requestID),
					attribute.HTTPRequestMethod(method),
					attribute.KeyValuePair(
						"http.request.headers",
						generateHeaderAttribute(request.Header),
					),
				)

				if method == resty.MethodPost ||
					method == resty.MethodPut ||
					method == resty.MethodPatch {
					span.SetAttributes(
						attribute.HTTPRequestBody(generateBodyAttribute(request.Body)),
					)
				}

				request.SetContext(ctx)
				return nil
			},
		).
		OnError(
			func(request *resty.Request, err error) {
				// OnError should be triggered when resty failed to make a request.
				// Thus, the tracer should mark the current span as error.
				span := tracer.SpanFromContext(request.Context())
				if span == nil {
					// Ignore if there are no span.
					return
				}
				defer span.End()

				span.SetAttributes(attribute.HTTPUrl(request.URL))
				span.RecordError(err)
			},
		).
		OnPanic(
			func(request *resty.Request, err error) {
				// OnPanic should be triggered after resty makes a request. Marks
				// the span as failed and record error.
				span := tracer.SpanFromContext(request.Context())
				if span == nil {
					// Ignore if there are no span.
					return
				}
				defer span.End()

				span.SetAttributes(attribute.HTTPUrl(request.URL))
				span.RecordError(err)
			},
		).
		OnAfterResponse(
			func(client *resty.Client, response *resty.Response) error {
				// Here we are going to get the span from the request's context.
				// After request was made, we need to capture the following attributes:
				// 1. The request body if there are any,
				// 2. The response status code, and
				// 3. The response header.
				span := tracer.SpanFromContext(response.Request.Context())
				if span == nil {
					// Ignore if there are no span.
					return nil
				}
				defer span.End()

				span.SetAttributes(
					attribute.HTTPResponseStatus(response.StatusCode()),
					attribute.HTTPUrl(response.Request.URL),
					attribute.KeyValuePair(
						"http.response.headers",
						generateHeaderAttribute(response.Header()),
					),
				)

				if response.Body() != nil {
					span.SetAttributes(
						attribute.HTTPResponseBody(
							generateBodyAttribute(response.Body()),
						),
					)
				}

				return nil
			},
		)
}

func generateBodyAttribute(body any) any {
	if body == nil {
		return nil
	}

	typeof := reflect.TypeOf(body)
	switch typeof.Kind() {
	case reflect.Map, reflect.Struct:
		b, err := json.MarshalIndent(body, "", "  ")
		if err != nil {
			return body
		}
		return string(b)
	case reflect.Slice:
		if b, ok := body.([]byte); ok {
			return string(b)
		} else {
			return fmt.Sprintf("%v", body)
		}
	default:
		return body
	}
}

func generateHeaderAttribute(header http.Header) any {
	if header == nil {
		return nil
	}

	b, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		return header
	}

	return string(b)
}
