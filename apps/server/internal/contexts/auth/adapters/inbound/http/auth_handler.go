package http

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/domain"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/authcontext"
)

type CookieConfig struct {
	Name     string
	Domain   string
	Path     string
	Secure   bool
	SameSite http.SameSite
}

type AuthHandler struct {
	useCase      ports.AuthUseCase
	cookieConfig CookieConfig
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type tokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresAt   string `json:"expiresAt"`
	UserID      string `json:"userId"`
}

func NewAuthHandler(useCase ports.AuthUseCase, cookieConfig CookieConfig) *AuthHandler {
	return &AuthHandler{useCase: useCase, cookieConfig: cookieConfig}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	result, err := h.useCase.Login(r.Context(), ports.LoginCommand{
		Email:     payload.Email,
		Password:  payload.Password,
		UserAgent: r.UserAgent(),
		IP:        clientIP(r),
	})
	if err != nil {
		h.writeAuthError(w, err)
		return
	}

	h.setRefreshCookie(w, result.RefreshToken, result.RefreshExpiresAt)
	h.writeTokenResponse(w, http.StatusOK, result)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, ok := h.readRefreshToken(r)
	if !ok {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

	result, err := h.useCase.Refresh(r.Context(), ports.RefreshCommand{
		RefreshToken: refreshToken,
		UserAgent:    r.UserAgent(),
		IP:           clientIP(r),
	})
	if err != nil {
		h.clearRefreshCookie(w)
		h.writeAuthError(w, err)
		return
	}

	h.setRefreshCookie(w, result.RefreshToken, result.RefreshExpiresAt)
	h.writeTokenResponse(w, http.StatusOK, result)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, _ := h.readRefreshToken(r)
	if err := h.useCase.Logout(r.Context(), ports.LogoutCommand{RefreshToken: refreshToken}); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	h.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	userID, ok := authcontext.GetUserID(r.Context())
	if !ok || userID == "" {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

	sessionID, _ := authcontext.GetSessionID(r.Context())
	h.writeJSON(w, http.StatusOK, map[string]string{
		"userId":    userID,
		"sessionId": sessionID,
	})
}

func (h *AuthHandler) writeTokenResponse(w http.ResponseWriter, status int, result ports.AuthTokens) {
	h.writeJSON(w, status, tokenResponse{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresAt:   result.ExpiresAt.UTC().Format(time.RFC3339),
		UserID:      result.UserID,
	})
}

func (h *AuthHandler) writeAuthError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrUnauthorized) || errors.Is(err, domain.ErrInvalidInput) || errors.Is(err, domain.ErrSessionNotFound) {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

	h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
}

func (h *AuthHandler) readRefreshToken(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(h.cookieConfig.Name)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return "", false
	}

	return cookie.Value, true
}

func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, value string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieConfig.Name,
		Value:    value,
		Path:     h.cookieConfig.Path,
		Domain:   h.cookieConfig.Domain,
		HttpOnly: true,
		Secure:   h.cookieConfig.Secure,
		SameSite: h.cookieConfig.SameSite,
		Expires:  expiresAt.UTC(),
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	})
}

func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieConfig.Name,
		Value:    "",
		Path:     h.cookieConfig.Path,
		Domain:   h.cookieConfig.Domain,
		HttpOnly: true,
		Secure:   h.cookieConfig.Secure,
		SameSite: h.cookieConfig.SameSite,
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
	})
}

func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
