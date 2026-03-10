package httptransport

import (
	"errors"
	"net/http"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
)

func authMiddleware(handler *AuthHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handler == nil || handler.authService == nil {
				next.ServeHTTP(w, r)
				return
			}

			path := r.URL.Path
			if isPublicPath(path) {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie(handler.cookieName)
			if err != nil || cookie == nil || cookie.Value == "" {
				writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
				return
			}

			result, err := handler.authService.Authenticate(r.Context(), auth.AuthenticateRequest{
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

			if result.ReplacementRaw != nil && result.CookieExpires != nil {
				now := time.Now().UTC()
				maxAge := int(result.CookieExpires.Sub(now).Seconds())
				if maxAge < 0 {
					maxAge = 0
				}

				http.SetCookie(w, &http.Cookie{
					Name:     handler.cookieName,
					Value:    *result.ReplacementRaw,
					Path:     "/",
					HttpOnly: true,
					Secure:   handler.cookieSecure,
					SameSite: parseSameSite(handler.cookieSameSiteRaw),
					MaxAge:   maxAge,
				})
			}

			next.ServeHTTP(w, r.WithContext(withAuthContext(r.Context(), result.UserID, result.SessionID)))
		})
	}
}

func isPublicPath(path string) bool {
	switch path {
	case "/health", "/api/health", "/openapi.yaml", "/api/auth/login", "/api/auth/logout", "/api/auth/session":
		return true
	default:
		return false
	}
}
