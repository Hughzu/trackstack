package httptransport

import (
	"net/http"
	"strings"
)

func isJSONRequest(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

func redirectWithError(w http.ResponseWriter, r *http.Request, fallback string) {
	referrer := strings.TrimSpace(r.Referer())
	if referrer != "" {
		fallback = referrer
	}

	if strings.Contains(fallback, "?") {
		fallback = fallback + "&error=1"
	} else {
		fallback = fallback + "?error=1"
	}

	redirect(w, fallback)
}

func redirect(w http.ResponseWriter, location string) {
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusSeeOther)
}
