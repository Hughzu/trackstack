package http

import (
	"encoding/json"
	"net/http"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/auth/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/auth/domain"
	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/authcontext"
)

type AuthHandler struct {
	useCase      ports.AuthUseCase
	cookieName   string
	cookieSecure bool
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewAuthHandler(useCase ports.AuthUseCase, cookieName string, cookieSecure bool) *AuthHandler {
	return &AuthHandler{
		useCase:      useCase,
		cookieName:   cookieName,
		cookieSecure: cookieSecure,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	token, err := h.useCase.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		if err == domain.ErrUnauthorized || err == domain.ErrInvalidInput {
			h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	// 30 days max-age
	maxAge := 30 * 24 * 60 * 60
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	userID, ok := authcontext.GetUserID(r.Context())
	if !ok || userID == "" {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

	// For compatibility with old `AuthHandler.Session` response shape
	h.writeJSON(w, http.StatusOK, map[string]string{
		"userId":    userID,
		"sessionId": "stateless-jwt",
	})
}

func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
