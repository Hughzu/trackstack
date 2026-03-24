package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/domain"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/authcontext"
)

type AuthHandler struct {
	useCase ports.AuthUseCase
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type loginResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresAt   string `json:"expiresAt"`
	UserID      string `json:"userId"`
}

func NewAuthHandler(useCase ports.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	result, err := h.useCase.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		if err == domain.ErrUnauthorized || err == domain.ErrInvalidInput {
			h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	h.writeJSON(w, http.StatusOK, loginResponse{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresAt:   result.ExpiresAt.UTC().Format(time.RFC3339),
		UserID:      result.UserID,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	userID, ok := authcontext.GetUserID(r.Context())
	if !ok || userID == "" {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

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
