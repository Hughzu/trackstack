package httptransport

import "net/http"

type HealthResponse struct {
	OK bool `json:"ok"`
}

func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{OK: true})
}
