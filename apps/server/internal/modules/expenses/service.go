package expenses

import (
	"errors"
)

var ErrInvalidInput = errors.New("invalid input")
var ErrNotFound = errors.New("not found")

type Service struct {
	store ExpensesStore
}

func NewService(store ExpensesStore) *Service {
	return &Service{store: store}
}
