package calories

import "errors"

var ErrInvalidInput = errors.New("invalid input")

type Service struct {
	store CaloriesStore
}

func NewService(store CaloriesStore) *Service {
	return &Service{store: store}
}
