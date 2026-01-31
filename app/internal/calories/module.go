package calories

import (
	"github.com/23St/trackstack/internal/common/db"
)

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

type NewHandlerFunc func(service Service) interface{}
