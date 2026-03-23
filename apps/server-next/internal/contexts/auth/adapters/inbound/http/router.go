package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *AuthHandler) RegisterRoutes(r chi.Router, requireSession func(http.Handler) http.Handler) {
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)
	r.With(requireSession).Get("/session", h.Session)
}
