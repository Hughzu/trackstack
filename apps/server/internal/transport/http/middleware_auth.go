package httptransport

import (
	"errors"
	"net/http"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
)

func authMiddleware(handler *AuthHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handler == nil || handler.authService == nil {
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

			if result.CookieExpires != nil {
				tokenValue := cookie.Value
				if result.ReplacementRaw != nil {
					tokenValue = *result.ReplacementRaw
				}

				setAuthCookie(w, handler.cookieName, tokenValue, handler.cookieSecure, handler.cookieSameSiteRaw, *result.CookieExpires)
			}

			next.ServeHTTP(w, r.WithContext(withAuthContext(r.Context(), result.UserID, result.SessionID)))
		})
	}
}
