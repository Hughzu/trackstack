package jsonhandler

import (
	"context"
	"net/http"
)

type WidgetRequest struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type WidgetResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type WidgetStore interface {
	Save(ctx context.Context, req WidgetRequest) (WidgetResponse, error)
}

func NewWidgetHandler(store WidgetStore) http.HandlerFunc {
	panic("TODO")
}
