package httptransport

import (
	"net/http"
	"time"
)

func setAuthCookie(w http.ResponseWriter, name string, value string, secure bool, sameSiteRaw string, expiresAt time.Time) {
	now := time.Now().UTC()
	maxAge := int(expiresAt.Sub(now).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: parseSameSite(sameSiteRaw),
		MaxAge:   maxAge,
	})
}

func clearAuthCookie(w http.ResponseWriter, name string, secure bool, sameSiteRaw string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: parseSameSite(sameSiteRaw),
		MaxAge:   -1,
	})
}
