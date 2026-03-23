package ports

import "context"

type AuthUseCase interface {
	Login(ctx context.Context, email string, password string) (string, error)
}
