package jsonhandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewWidgetHandler(t *testing.T) {
	t.Parallel()

	t.Run("method not allowed", func(t *testing.T) {
		t.Parallel()

		handler := NewWidgetHandler(fakeWidgetStore{})
		req := httptest.NewRequest(http.MethodGet, "/widgets", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assertJSONStatus(t, res, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		handler := NewWidgetHandler(fakeWidgetStore{})
		req := httptest.NewRequest(http.MethodPost, "/widgets", strings.NewReader(`{"name":`))
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assertJSONStatus(t, res, http.StatusBadRequest, map[string]any{"error": "invalid json body"})
	})

	t.Run("unknown field", func(t *testing.T) {
		t.Parallel()

		handler := NewWidgetHandler(fakeWidgetStore{})
		req := httptest.NewRequest(http.MethodPost, "/widgets", strings.NewReader(`{"name":"box","count":2,"wat":true}`))
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assertJSONStatus(t, res, http.StatusBadRequest, map[string]any{"error": "invalid json body"})
	})

	t.Run("missing required fields", func(t *testing.T) {
		t.Parallel()

		handler := NewWidgetHandler(fakeWidgetStore{})
		req := httptest.NewRequest(http.MethodPost, "/widgets", strings.NewReader(`{"name":"","count":0}`))
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assertJSONStatus(t, res, http.StatusBadRequest, map[string]any{"error": "missing required fields"})
	})

	t.Run("store error", func(t *testing.T) {
		t.Parallel()

		handler := NewWidgetHandler(fakeWidgetStore{err: errors.New("boom")})
		req := httptest.NewRequest(http.MethodPost, "/widgets", strings.NewReader(`{"name":"box","count":2}`))
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assertJSONStatus(t, res, http.StatusInternalServerError, map[string]any{"error": "internal server error"})
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		store := fakeWidgetStore{response: WidgetResponse{ID: "widget-1", Name: "box", Count: 2}}
		handler := NewWidgetHandler(store)
		req := httptest.NewRequest(http.MethodPost, "/widgets", strings.NewReader(`{"name":"box","count":2}`))
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assertJSONStatus(t, res, http.StatusOK, map[string]any{"id": "widget-1", "name": "box", "count": float64(2)})
	})
}

type fakeWidgetStore struct {
	response WidgetResponse
	err      error
}

func (f fakeWidgetStore) Save(_ context.Context, _ WidgetRequest) (WidgetResponse, error) {
	if f.err != nil {
		return WidgetResponse{}, f.err
	}
	return f.response, nil
}

func assertJSONStatus(t *testing.T, res *httptest.ResponseRecorder, wantStatus int, wantBody map[string]any) {
	t.Helper()

	if res.Code != wantStatus {
		t.Fatalf("status = %d, want %d", res.Code, wantStatus)
	}

	if got := res.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type = %q, want application/json", got)
	}

	var gotBody map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &gotBody); err != nil {
		t.Fatalf("decode response json: %v", err)
	}

	if len(gotBody) != len(wantBody) {
		t.Fatalf("response body = %#v, want %#v", gotBody, wantBody)
	}

	for key, wantValue := range wantBody {
		if gotBody[key] != wantValue {
			t.Fatalf("response body[%q] = %#v, want %#v", key, gotBody[key], wantValue)
		}
	}
}
