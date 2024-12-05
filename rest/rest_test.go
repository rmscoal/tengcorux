package rest

import (
	"context"
	"github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"github.com/rmscoal/tengcorux/tracer/tracetest"
	"net/http"
	"strings"
	"testing"
	"time"
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
		if rest == nil {
			t.Fatal("rest should not be nil")
		}

		t.Run("Hit Get /success", func(t *testing.T) {
			resp, err := rest.R().
				SetContext(context.Background()).
				Get("/success")
			if err != nil {
				t.Fatalf("error should be nil, got %v", err)
			} else if string(resp.Body()) != "Success" {
				t.Fatalf("body should be Success, got %v", string(resp.Body()))
			} else if resp.StatusCode() != http.StatusOK {
				t.Fatal("status code should be 200")
			}
		})

		t.Run("Hit Post /success", func(t *testing.T) {
			resp, err := rest.R().SetContext(context.Background()).Post("/success")
			if err != nil {
				t.Fatalf("error should be nil, got %v", err)
			} else if string(resp.Body()) != "Success" {
				t.Fatalf("body should be Success, got %v", string(resp.Body()))
			} else if resp.StatusCode() != http.StatusOK {
				t.Fatal("status code should be 200")
			}
		})

	})

	t.Run("WithTracerEnabled", func(t *testing.T) {
		t.Run("Hit Get /success", func(t *testing.T) {
			tr := tracetest.NewTracer()
			tracer.SetGlobalTracer(tr)
			rest := New(WithTracerEnabled()).SetBaseURL("http://localhost:8123")
			if rest == nil {
				t.Fatal("rest should not be nil")
			}

			resp, err := rest.R().SetContext(context.Background()).Get("/success")
			if err != nil {
				t.Fatalf("error should be nil, got %v", err)
			} else if string(resp.Body()) != "Success" {
				t.Fatalf("body should be Success, got %v", string(resp.Body()))
			} else if resp.StatusCode() != http.StatusOK {
				t.Fatal("status code should be 200")
			}

			// Check the tracer
			ended := tr.Recorder().EndedSpans()
			if len(ended) != 1 {
				t.Fatalf("ended spans should be 1, got %v", len(ended))
			}
			span := ended[0]
			if span.Type != tracer.SpanTypeExit {
				t.Fatalf("span.Type should be %d, got %d", tracer.SpanTypeExit, span.Type)
			} else if span.Layer != tracer.SpanLayerHttp {
				t.Fatalf("layer should be %d, got %d", tracer.SpanLayerHttp, span.Layer)
			}

			for _, attr := range span.Attributes {
				switch attr.Key {
				case attribute.HTTPRequestIDKey:
					if attr.Value == "" {
						t.Fatalf("the request id span attribute should not be empty")
					}
				case attribute.HTTPRequestMethodKey:
					if attr.Value != "GET" {
						t.Fatalf("http request method should be GET, got %v", attr.Value)
					}
				case attribute.Key("http.request.headers"):
					if attr.Value == nil {
						t.Fatalf("http request headers should not be nil")
					}
				case attribute.HTTPUrlKey:
					if attr.Value != "http://localhost:8123/success" {
						t.Fatalf("http url should be http://localhost:8123/success, got %v", attr.Value)
					}
				case attribute.HTTPResponseStatusKey:
					if attr.Value != 200 {
						t.Fatalf("http response status should be 200, got %v", attr.Value)
					}
				case attribute.HTTPResponseBodyKey:
					if attr.Value != "Success" {
						t.Fatalf("http response body should be Success, got %v", attr.Value)
					}
				}
			}
		})

		t.Run("Hit Get /error", func(t *testing.T) {
			tr := tracetest.NewTracer()
			tracer.SetGlobalTracer(tr)
			rest := New(WithTracerEnabled()).SetBaseURL("http://localhost:8123")
			if rest == nil {
				t.Fatal("rest should not be nil")
			}

			resp, err := rest.R().SetContext(context.Background()).Get("/error")
			if err != nil {
				t.Fatalf("error should be nil, got %v", err)
			} else if strings.Trim(string(resp.Body()), " ") != "Bad Request" {
				t.Fatalf("body should be Bad Request, got %v", string(resp.Body()))
			} else if resp.StatusCode() != http.StatusBadRequest {
				t.Fatal("status code should be 400")
			}

			// Check the tracer
			ended := tr.Recorder().EndedSpans()
			if len(ended) != 1 {
				t.Fatalf("ended spans should be 1, got %v", len(ended))
			}
			span := ended[0]
			if span.Type != tracer.SpanTypeExit {
				t.Fatalf("span.Type should be %d, got %d", tracer.SpanTypeExit, span.Type)
			} else if span.Layer != tracer.SpanLayerHttp {
				t.Fatalf("layer should be %d, got %d", tracer.SpanLayerHttp, span.Layer)
			}

			for _, attr := range span.Attributes {
				switch attr.Key {
				case attribute.HTTPRequestIDKey:
					if attr.Value == "" {
						t.Fatalf("the request id span attribute should not be empty")
					}
				case attribute.HTTPRequestMethodKey:
					if attr.Value != "GET" {
						t.Fatalf("http request method should be GET, got %v", attr.Value)
					}
				case attribute.Key("http.request.headers"):
					if attr.Value == nil {
						t.Fatalf("http request headers should not be nil")
					}
				case attribute.HTTPUrlKey:
					if attr.Value != "http://localhost:8123/error" {
						t.Fatalf("http url should be http://localhost:8123/error, got %v", attr.Value)
					}
				case attribute.HTTPResponseStatusKey:
					if attr.Value != 400 {
						t.Fatalf("http response status should be 400, got %v", attr.Value)
					}
				case attribute.HTTPResponseBodyKey:
					if attr.Value != "Bad Request" {
						t.Fatalf("http response body should be Bad Request, got %v", attr.Value)
					}
				}
			}
		})

		t.Run("Hit Post /error", func(t *testing.T) {
			tr := tracetest.NewTracer()
			tracer.SetGlobalTracer(tr)
			rest := New(WithTracerEnabled()).SetBaseURL("http://localhost:8123")
			if rest == nil {
				t.Fatal("rest should not be nil")
			}

			body := map[string]interface{}{
				"hello": "world",
			}
			resp, err := rest.R().
				SetContext(context.Background()).
				SetBody(body).
				Post("/error")
			if err != nil {
				t.Fatalf("error should be nil, got %v", err)
			} else if strings.Trim(string(resp.Body()), " ") != "Bad Request" {
				t.Fatalf("body should be Bad Request, got %v", string(resp.Body()))
			} else if resp.StatusCode() != http.StatusBadRequest {
				t.Fatal("status code should be 400")
			}

			// Check the tracer
			ended := tr.Recorder().EndedSpans()
			if len(ended) != 1 {
				t.Fatalf("ended spans should be 1, got %v", len(ended))
			}
			span := ended[0]
			if span.Type != tracer.SpanTypeExit {
				t.Fatalf("span.Type should be %d, got %d", tracer.SpanTypeExit, span.Type)
			} else if span.Layer != tracer.SpanLayerHttp {
				t.Fatalf("layer should be %d, got %d", tracer.SpanLayerHttp, span.Layer)
			}

			for _, attr := range span.Attributes {
				switch attr.Key {
				case attribute.HTTPRequestIDKey:
					if attr.Value == "" {
						t.Fatalf("the request id span attribute should not be empty")
					}
				case attribute.HTTPRequestMethodKey:
					if attr.Value != "POST" {
						t.Fatalf("http request method should be POST, got %v", attr.Value)
					}
				case attribute.Key("http.request.headers"):
					if attr.Value == nil {
						t.Fatalf("http request headers should not be nil")
					}
				case attribute.HTTPUrlKey:
					if attr.Value != "http://localhost:8123/error" {
						t.Fatalf("http url should be http://localhost:8123/error, got %v", attr.Value)
					}
				case attribute.HTTPRequestBodyKey:
					asMap := attr.Value.(map[string]interface{})
					if asMap["hello"] != "world" {
						t.Errorf("http request body should be %v, got %v", body, attr.Value)
					}
				case attribute.HTTPResponseStatusKey:
					if attr.Value != 400 {
						t.Fatalf("http response status should be 400, got %v", attr.Value)
					}
				case attribute.HTTPResponseBodyKey:
					if attr.Value != "Bad Request" {
						t.Fatalf("http response body should be Bad Request, got %v", attr.Value)
					}
				}
			}
		})
	})
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