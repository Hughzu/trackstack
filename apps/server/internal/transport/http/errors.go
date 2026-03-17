package httptransport

import (
	"errors"
	"net/http"

	heatdomain "github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, heatdomain.ErrInvalidInput) {
		status = http.StatusBadRequest
	}

	writeJSON(w, status, errorResponse{Error: err.Error()})
}
