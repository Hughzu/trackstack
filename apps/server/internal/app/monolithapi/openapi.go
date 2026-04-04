package monolithapi

import (
	"embed"
	"net/http"
)

//go:embed openapi.yaml
var openAPISpec embed.FS

func openapi(w http.ResponseWriter, _ *http.Request) {
	data, err := openAPISpec.ReadFile("openapi.yaml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
