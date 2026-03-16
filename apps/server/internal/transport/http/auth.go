package httptransport

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/users"
)

type AuthHandler struct {
	authService       *auth.Service
	usersService      *users.Service
	cookieName        string
	cookieSecure      bool
	cookieSameSiteRaw string
}

type authSessionResponse struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer logAuthHTTPTime(r, "login.total", start, nil)

	payload, err := readAuthPayload(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(payload.Email))
	password := payload.Password
	if email == "" || password == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing credentials"})
		return
	}

	findUserStart := time.Now()
	user, err := h.usersService.FindByEmail(r.Context(), email)
	logAuthHTTPTime(r, "login.find_user", findUserStart, err)
	if err != nil {
		if !errors.Is(err, users.ErrNotFound) {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
			return
		}
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

	verifyStart := time.Now()
	if !auth.VerifyPassword(password, user.PasswordHash) {
		logAuthHTTPTime(r, "login.verify_password", verifyStart, auth.ErrUnauthorized)
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}
	logAuthHTTPTime(r, "login.verify_password", verifyStart, nil)

	createSessionStart := time.Now()
	rawToken, session, err := h.authService.CreateSession(r.Context(), auth.CreateSessionRequest{
		UserID:  user.ID,
		Context: getClientContext(r),
	})
	logAuthHTTPTime(r, "login.create_session", createSessionStart, err)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	setAuthCookie(w, h.cookieName, rawToken, h.cookieSecure, h.cookieSameSiteRaw, expiresAt)

	now := time.Now().UTC()

	h.updateLastLoginAsync(r, user.ID, now.Format(time.RFC3339))

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieName)
	if err == nil && cookie != nil && cookie.Value != "" {
		_ = h.authService.RevokeSessionByRawToken(r.Context(), auth.RevokeSessionRequest{RawToken: cookie.Value})
	}

	clearAuthCookie(w, h.cookieName, h.cookieSecure, h.cookieSameSiteRaw)

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieName)
	if err != nil || cookie == nil || strings.TrimSpace(cookie.Value) == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

	result, err := h.authService.Authenticate(r.Context(), auth.AuthenticateRequest{
		RawToken: cookie.Value,
		Context:  getClientContext(r),
	})
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	if result.CookieExpires != nil {
		tokenValue := cookie.Value
		if result.ReplacementRaw != nil {
			tokenValue = *result.ReplacementRaw
		}

		setAuthCookie(w, h.cookieName, tokenValue, h.cookieSecure, h.cookieSameSiteRaw, *result.CookieExpires)
	}

	writeJSON(w, http.StatusOK, authSessionResponse{
		UserID:    result.UserID,
		SessionID: result.SessionID,
	})
}

type authPayload struct {
	Email    string
	Password string
}

func readAuthPayload(r *http.Request) (authPayload, error) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		return authPayload{}, err
	}
	return authPayload{Email: payload.Email, Password: payload.Password}, nil
}

func parseSameSite(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

func getClientContext(r *http.Request) auth.ClientContext {
	userAgent := strings.TrimSpace(r.UserAgent())
	var userAgentPtr *string
	if userAgent != "" {
		userAgentPtr = &userAgent
	}

	ip := strings.TrimSpace(firstHeaderValue(r, "X-Forwarded-For"))
	if ip == "" {
		ip = strings.TrimSpace(firstHeaderValue(r, "X-Real-IP"))
	}
	ipPrefix := hashIPPrefix(ip)

	return auth.ClientContext{
		UserAgent: userAgentPtr,
		IPPrefix:  ipPrefix,
	}
}

func firstHeaderValue(r *http.Request, key string) string {
	value := r.Header.Get(key)
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(parts[0])
}

func hashIPPrefix(ip string) *string {
	if ip == "" {
		return nil
	}

	if strings.Contains(ip, ".") {
		parts := strings.Split(ip, ".")
		if len(parts) >= 3 {
			value := parts[0] + "." + parts[1] + "." + parts[2] + ".0"
			return &value
		}
		value := ip
		return &value
	}

	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		if len(parts) > 4 {
			parts = parts[:4]
		}
		value := strings.Join(parts, ":") + "::"
		return &value
	}

	value := ip
	return &value
}

func logAuthHTTPTime(r *http.Request, step string, start time.Time, err error) {
	attrs := []any{"step", step, "duration", time.Since(start), "path", r.URL.Path}
	if err != nil {
		attrs = append(attrs, "error", err)
	}
	slog.DebugContext(r.Context(), "auth http timing", attrs...)
}

func (h *AuthHandler) updateLastLoginAsync(r *http.Request, userID string, lastLoginAt string) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 2*time.Second)

	go func() {
		defer cancel()

		start := time.Now()
		err := h.usersService.UpdateLastLogin(ctx, userID, lastLoginAt)
		logAuthHTTPTime(r, "login.update_last_login", start, err)
	}()
}
