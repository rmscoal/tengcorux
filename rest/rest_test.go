package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/rmscoal/tengcorux/reqid"
	"github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"github.com/rmscoal/tengcorux/tracer/tracetest"
	"github.com/stretchr/testify/assert"
)

func TestRest_New(t *testing.T) {
	server := testServer()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	go func() {
		_ = server.ListenAndServe()
	}()

	t.Run("NoOption", func(t *testing.T) {
		rest := New().SetBaseURL("http://localhost:8123")
		assert.NotNil(t, rest, "rest should not be nil")

		t.Run("Hit Get /success", func(t *testing.T) {
			resp, err := rest.R().SetContext(context.Background()).Get("/success")
			assert.NoError(t, err, "error should be nil")
			assert.Equal(t, "Success", string(resp.Body()),
				"body should be Success")
			assert.Equal(t, http.StatusOK, resp.StatusCode(),
				"status code should be 200")
		})

		t.Run("Hit Post /success", func(t *testing.T) {
			resp, err := rest.R().SetContext(context.Background()).Post("/success")
			assert.NoError(t, err, "error should be nil")
			assert.Equal(t, "Success", string(resp.Body()),
				"body should be Success")
			assert.Equal(t, http.StatusOK, resp.StatusCode(),
				"status code should be 200")
		})
	})

	t.Run("WithTracerEnabled", func(t *testing.T) {
		t.Run("Hit Get /success", func(t *testing.T) {
			tr := tracetest.NewTracer()
			tracer.SetGlobalTracer(tr)
			rest := New(WithTracerEnabled()).SetBaseURL("http://localhost:8123")
			assert.NotNil(t, rest, "rest should not be nil")

			resp, err := rest.R().SetContext(reqid.Inject(context.Background())).Get("/success")
			assert.NoError(t, err, "error should be nil")
			assert.Equal(t, "Success", string(resp.Body()),
				"body should be Success")
			assert.Equal(t, http.StatusOK, resp.StatusCode(),
				"status code should be 200")

			ended := tr.Recorder().EndedSpans()
			assert.Len(t, ended, 1, "ended spans should be 1")
			span := ended[0]
			assert.Equal(t, tracer.SpanTypeExit, span.Type,
				"span.Type should be correct")
			assert.Equal(t, tracer.SpanLayerHttp, span.Layer,
				"span.Layer should be correct")

			for _, attr := range span.Attributes {
				switch attr.Key {
				case attribute.HTTPRequestIDKey:
					assert.NotEmpty(t, attr.Value,
						"the request id span attribute should not be empty")
				case attribute.HTTPRequestMethodKey:
					assert.Equal(t, "GET", attr.Value,
						"http request method should be GET")
				case "http.request.headers":
					assert.NotNil(t, attr.Value,
						"http request headers should not be nil")
				case attribute.HTTPUrlKey:
					assert.Equal(t, "http://localhost:8123/success", attr.Value,
						"http url should be correct")
				case attribute.HTTPResponseStatusKey:
					assert.Equal(t, 200, attr.Value,
						"http response status should be 200")
				case attribute.HTTPResponseBodyKey:
					assert.Equal(t, "Success", attr.Value,
						"http response body should be correct")
				}
			}
		})

		t.Run("Hit Get /error", func(t *testing.T) {
			tr := tracetest.NewTracer()
			tracer.SetGlobalTracer(tr)
			rest := New(WithTracerEnabled()).SetBaseURL("http://localhost:8123")
			assert.NotNil(t, rest, "rest should not be nil")

			resp, err := rest.R().SetContext(reqid.Inject(context.Background())).Get("/error")
			assert.NoError(t, err, "error should be nil")
			assert.Equal(t, "Bad Request",
				strings.TrimSpace(string(resp.Body())),
				"body should be Bad Request")
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode(),
				"status code should be 400")

			ended := tr.Recorder().EndedSpans()
			assert.Len(t, ended, 1, "ended spans should be 1")
			span := ended[0]
			assert.Equal(t, tracer.SpanTypeExit, span.Type,
				"span.Type should be correct")
			assert.Equal(t, tracer.SpanLayerHttp, span.Layer,
				"span.Layer should be correct")

			for _, attr := range span.Attributes {
				switch attr.Key {
				case attribute.HTTPRequestIDKey:
					assert.NotEmpty(t, attr.Value,
						"the request id span attribute should not be empty")
				case attribute.HTTPRequestMethodKey:
					assert.Equal(t, "GET", attr.Value,
						"http request method should be GET")
				case "http.request.headers":
					assert.NotNil(t, attr.Value,
						"http request headers should not be nil")
				case attribute.HTTPUrlKey:
					assert.Equal(t, "http://localhost:8123/error", attr.Value,
						"http url should be correct")
				case attribute.HTTPResponseStatusKey:
					assert.Equal(t, 400, attr.Value,
						"http response status should be 400")
				case attribute.HTTPResponseBodyKey:
					assert.Equal(t, "Bad Request", attr.Value,
						"http response body should be correct")
				}
			}
		})

		t.Run("Hit Post /error", func(t *testing.T) {
			tr := tracetest.NewTracer()
			tracer.SetGlobalTracer(tr)
			rest := New(WithTracerEnabled()).SetBaseURL("http://localhost:8123")
			assert.NotNil(t, rest, "rest should not be nil")

			body := map[string]interface{}{
				"hello": "world",
			}
			resp, err := rest.R().SetContext(reqid.Inject(context.Background())).SetBody(body).Post("/error")
			assert.NoError(t, err, "error should be nil")
			assert.Equal(t, "Bad Request",
				strings.TrimSpace(string(resp.Body())),
				"body should be Bad Request")
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode(),
				"status code should be 400")

			ended := tr.Recorder().EndedSpans()
			assert.Len(t, ended, 1, "ended spans should be 1")
			span := ended[0]
			assert.Equal(t, tracer.SpanTypeExit, span.Type,
				"span.Type should be correct")
			assert.Equal(t, tracer.SpanLayerHttp, span.Layer,
				"span.Layer should be correct")

			for _, attr := range span.Attributes {
				switch attr.Key {
				case attribute.HTTPRequestIDKey:
					assert.NotEmpty(t, attr.Value,
						"the request id span attribute should not be empty")
				case attribute.HTTPRequestMethodKey:
					assert.Equal(t, "POST", attr.Value,
						"http request method should be POST")
				case "http.request.headers":
					assert.NotNil(t, attr.Value,
						"http request headers should not be nil")
				case attribute.HTTPUrlKey:
					assert.Equal(t, "http://localhost:8123/error", attr.Value,
						"http url should be correct")
				case attribute.HTTPRequestBodyKey:
					assert.NotNil(t, attr.Value,
						"http request body should not be nil")
				case attribute.HTTPResponseStatusKey:
					assert.Equal(t, 400, attr.Value,
						"http response status should be 400")
				case attribute.HTTPResponseBodyKey:
					assert.Equal(t, "Bad Request", attr.Value,
						"http response body should be correct")
				}
			}
		})

		t.Run("Error", func(t *testing.T) {
			tr := tracetest.NewTracer()
			tracer.SetGlobalTracer(tr)
			rest := New(WithTracerEnabled()).SetBaseURL("http://localhost:8123")
			rest.SetTransport(&http.Transport{
				DialContext: func(
					ctx context.Context, network, addr string,
				) (net.Conn, error) {
					return nil, errors.New("panic")
				},
			})
			_, err := rest.R().SetContext(context.TODO()).Get("/error")
			assert.Error(t, err)
		})
	})
}

func TestGenerateBodyAttribute(t *testing.T) {
	t.Run("Nil Body", func(t *testing.T) {
		assert.Nil(t, generateBodyAttribute(nil))
	})

	t.Run("Byte Body", func(t *testing.T) {
		res := generateBodyAttribute([]byte(`{"hello": "world"}`))
		expected := `{"hello": "world"}`
		assert.Equal(t, expected, res)
		t.Log("GenerateBodyAttribute Result:", res)
	})

	t.Run("Struct Body", func(t *testing.T) {
		type req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		body := req{
			Username: "test",
			Password: "test",
		}

		marshalled, _ := json.MarshalIndent(body, "", "  ")
		expected := string(marshalled)

		res := generateBodyAttribute(body)
		assert.Equal(t, expected, res)
		t.Log("GenerateBodyAttribute Result:", res)
	})

	t.Run("Map Any Body", func(t *testing.T) {
		body := map[string]interface{}{
			"username": "test",
			"password": 1,
		}

		marshalled, _ := json.MarshalIndent(body, "", "  ")
		expected := string(marshalled)

		res := generateBodyAttribute(body)
		assert.Equal(t, expected, res)
		t.Log("GenerateBodyAttribute Result:", res)
	})

	t.Run("Map String Body", func(t *testing.T) {
		body := map[string]string{
			"username": "test",
			"password": "test",
		}

		marshalled, _ := json.MarshalIndent(body, "", "  ")
		expected := string(marshalled)

		res := generateBodyAttribute(body)
		assert.Equal(t, expected, res)
		t.Log("GenerateBodyAttribute Result:", res)
	})

	t.Run("Int Body", func(t *testing.T) {
		assert.Equal(t, 1, generateBodyAttribute(1))
	})

	t.Run("Slice of Int Body", func(t *testing.T) {
		assert.Equal(t, "[1 2]", generateBodyAttribute([]int{1, 2}))
	})
}

func TestGenerateHeaderAttribute(t *testing.T) {
	head := http.Header{}
	head.Set("Content-Type", "application/x-www-form-urlencoded")
	head.Set("X-Request-ID", "some_request_id")
	head.Set("Topic", "RestPackage")

	res := generateHeaderAttribute(head)
	expected, _ := json.MarshalIndent(head, "", "  ")
	assert.Equal(t, string(expected), res)
	t.Log("GenerateHeaderAttribute Result:", res)
}

func testServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Success"))
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request"))
	})

	server := &http.Server{Handler: mux, Addr: ":8123"}
	return server
}
