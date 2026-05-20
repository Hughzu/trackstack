package normalizeemail

import "errors"

var ErrEmailRequired = errors.New("email is required")

func NormalizeEmail(raw string) (string, error) {
	panic("TODO")
}
