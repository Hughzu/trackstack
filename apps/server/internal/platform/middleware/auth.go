package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/domain"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/authcontext"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func ResolveSession(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenValue, ok := bearerToken(r)
			if !ok {
				writeUnauthorized(w)
				return
			}

			claims := &domain.SessionClaims{}
			token, err := jwtv5.ParseWithClaims(tokenValue, claims, func(token *jwtv5.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			}, jwtv5.WithValidMethods([]string{jwtv5.SigningMethodHS256.Alg()}))

			if err != nil || !token.Valid || claims.UserID == "" {
				writeUnauthorized(w)
				return
			}

			ctx := authcontext.WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(r *http.Request) (string, bool) {
	fields := strings.Fields(strings.TrimSpace(r.Header.Get("Authorization")))
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") || fields[1] == "" {
		fields = strings.Fields(strings.TrimSpace(r.Header.Get("X-Trackstack-Authorization")))
	}

	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") || fields[1] == "" {
		return "", false
	}

	return fields[1], true
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
}
