package functionurl

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestWrapBufferedReturnsJSONBodyAndCookies(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "abc123", Path: "/", HttpOnly: true})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	wrapped := WrapBuffered(handler)
	response, err := wrapped(context.Background(), events.LambdaFunctionURLRequest{
		RawPath: "/api/auth/session",
		RequestContext: events.LambdaFunctionURLRequestContext{
			DomainName: "example.lambda-url.eu-west-1.on.aws",
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method:   http.MethodGet,
				SourceIP: "127.0.0.1",
			},
		},
	})
	if err != nil {
		t.Fatalf("wrapped handler returned error: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}
	if response.Body != `{"ok":true}` {
		t.Fatalf("expected JSON body, got %q", response.Body)
	}
	if response.Headers["Content-Type"] != "application/json" {
		t.Fatalf("expected content type header, got %q", response.Headers["Content-Type"])
	}
	if len(response.Cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(response.Cookies))
	}
}

func TestWrapBufferedMapsCookiesIntoRequestHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("Cookie")))
	})

	wrapped := WrapBuffered(handler)
	response, err := wrapped(context.Background(), events.LambdaFunctionURLRequest{
		RawPath: "/api/auth/session",
		Cookies: []string{"trackstack_session=abc123", "theme=dark"},
		RequestContext: events.LambdaFunctionURLRequestContext{
			DomainName: "example.lambda-url.eu-west-1.on.aws",
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method:   http.MethodGet,
				SourceIP: "127.0.0.1",
			},
		},
	})
	if err != nil {
		t.Fatalf("wrapped handler returned error: %v", err)
	}

	if response.Body != "trackstack_session=abc123; theme=dark" {
		t.Fatalf("expected cookie header to be reconstructed, got %q", response.Body)
	}
}

func TestWrapBufferedEncodesBinaryResponses(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte{0x00, 0xff, 0x10})
	})

	wrapped := WrapBuffered(handler)
	response, err := wrapped(context.Background(), events.LambdaFunctionURLRequest{
		RawPath: "/binary",
		RequestContext: events.LambdaFunctionURLRequestContext{
			DomainName: "example.lambda-url.eu-west-1.on.aws",
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method:   http.MethodGet,
				SourceIP: "127.0.0.1",
			},
		},
	})
	if err != nil {
		t.Fatalf("wrapped handler returned error: %v", err)
	}

	if !response.IsBase64Encoded {
		t.Fatal("expected binary response to be base64 encoded")
	}
	if response.Body == "" {
		t.Fatal("expected binary response body")
	}
	if bytes.Equal([]byte(response.Body), []byte{0x00, 0xff, 0x10}) {
		t.Fatal("expected encoded body, got raw bytes")
	}
}
