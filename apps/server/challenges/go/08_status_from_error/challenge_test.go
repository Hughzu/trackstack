package errstatus

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestStatusFromError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "invalid input", err: fmt.Errorf("wrap: %w", ErrInvalidInput), want: http.StatusBadRequest},
		{name: "unauthorized", err: fmt.Errorf("wrap: %w", ErrUnauthorized), want: http.StatusUnauthorized},
		{name: "not found", err: fmt.Errorf("wrap: %w", ErrNotFound), want: http.StatusNotFound},
		{name: "other error", err: errors.New("boom"), want: http.StatusInternalServerError},
		{name: "nil error", err: nil, want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := StatusFromError(tt.err)
			if got != tt.want {
				t.Fatalf("StatusFromError() = %d, want %d", got, tt.want)
			}
		})
	}
}
