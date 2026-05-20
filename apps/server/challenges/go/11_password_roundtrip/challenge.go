package passwordroundtrip

import "errors"

const HashPrefix = "$scrypt$"

var ErrPasswordRequired = errors.New("password is required")

func HashPassword(password string) (string, error) {
	panic("TODO")
}

func VerifyPassword(password string, storedHash string) bool {
	panic("TODO")
}
