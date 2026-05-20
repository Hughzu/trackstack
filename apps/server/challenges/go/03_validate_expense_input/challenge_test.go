package expenseinput

import (
	"errors"
	"testing"
)

func TestCreateExpenseInputValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   CreateExpenseInput
		wantErr error
	}{
		{name: "missing user id", input: CreateExpenseInput{Label: "Lunch", Amount: 12}, wantErr: ErrUserIDRequired},
		{name: "missing label", input: CreateExpenseInput{UserID: "user-1", Amount: 12}, wantErr: ErrLabelRequired},
		{name: "invalid amount", input: CreateExpenseInput{UserID: "user-1", Label: "Lunch", Amount: 0}, wantErr: ErrAmountInvalid},
		{name: "valid input", input: CreateExpenseInput{UserID: "user-1", Label: "Lunch", Amount: 12}, wantErr: nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.input.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
