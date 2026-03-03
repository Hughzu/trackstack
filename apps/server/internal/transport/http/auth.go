package httptransport

import (
	"errors"
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	isJSON := strings.Contains(r.Header.Get("Content-Type"), "application/json")

	payload, ok := readAuthPayload(r, isJSON)
	if !ok {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(payload.Email))
	password := payload.Password
	if email == "" || password == "" {
		if !isJSON {
			redirect(w, "/login?error=1")
			return
		}
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing credentials"})
		return
	}

	user, err := h.usersService.FindByEmail(r.Context(), email)
	if err != nil {
		if !errors.Is(err, users.ErrNotFound) {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
			return
		}
		if !isJSON {
			redirect(w, "/login?error=1")
			return
		}
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

	if !auth.VerifyPassword(password, user.PasswordHash) {
		if !isJSON {
			redirect(w, "/login?error=1")
			return
		}
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

	rawToken, session, err := h.authService.CreateSession(r.Context(), auth.CreateSessionRequest{
		UserID:  user.ID,
		Context: getClientContext(r),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	now := time.Now().UTC()
	maxAge := int(expiresAt.Sub(now).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    rawToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: parseSameSite(h.cookieSameSiteRaw),
		MaxAge:   maxAge,
	})

	_ = h.usersService.UpdateLastLogin(r.Context(), user.ID, now.Format(time.RFC3339))

	if !isJSON {
		redirect(w, "/")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieName)
	if err == nil && cookie != nil && cookie.Value != "" {
		_ = h.authService.RevokeSessionByRawToken(r.Context(), auth.RevokeSessionRequest{RawToken: cookie.Value})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: parseSameSite(h.cookieSameSiteRaw),
		MaxAge:   -1,
	})

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		redirect(w, "/login")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type authPayload struct {
	Email    string
	Password string
}

func readAuthPayload(r *http.Request, isJSON bool) (authPayload, bool) {
	if isJSON {
		var payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := decodeJSON(r, &payload); err != nil {
			return authPayload{}, false
		}
		return authPayload{Email: payload.Email, Password: payload.Password}, true
	}

	if err := r.ParseForm(); err != nil {
		return authPayload{}, true
	}

	return authPayload{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}, true
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
