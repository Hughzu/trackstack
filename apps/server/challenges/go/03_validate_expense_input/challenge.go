package expenseinput

import "errors"

var (
	ErrUserIDRequired = errors.New("user id is required")
	ErrLabelRequired  = errors.New("label is required")
	ErrAmountInvalid  = errors.New("amount must be greater than zero")
)

type CreateExpenseInput struct {
	UserID string
	Label  string
	Amount int
}

func (in CreateExpenseInput) Validate() error {
	panic("TODO")
}
