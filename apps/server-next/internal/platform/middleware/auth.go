package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/auth/domain"
	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/authcontext"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func ResolveSession(jwtSecret string, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookieName)
			if err != nil || cookie == nil || cookie.Value == "" {
				writeUnauthorized(w)
				return
			}

			claims := &domain.SessionClaims{}
			token, err := jwtv5.ParseWithClaims(cookie.Value, claims, func(token *jwtv5.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid || claims.UserID == "" {
				writeUnauthorized(w)
				return
			}

			ctx := authcontext.WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
}
