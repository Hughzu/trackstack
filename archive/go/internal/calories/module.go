package calories

import (
	"net/http"

	"github.com/23St/trackstack/internal/common/db"
)

// RouteRegistrar is the interface for registering routes
type RouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux)
}

type Module struct {
	// External
	Service Service

	// Internal
	repo Repository
}

func NewModule(database *db.DB) *Module {
	repo := NewRepository(database)
	service := NewService(repo)

	return &Module{
		repo:    repo,
		Service: service,
	}
}

// RegisterWithHandler registers routes using the provided handler
func (m *Module) RegisterWithHandler(mux *http.ServeMux, handler RouteRegistrar) {
	handler.RegisterRoutes(mux)
}
