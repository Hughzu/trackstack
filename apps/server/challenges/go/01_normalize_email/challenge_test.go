package normalizeemail

import (
	"errors"
	"testing"
)

func TestNormalizeEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{name: "trim and lowercase", input: "  Alice.Example@Example.COM  ", want: "alice.example@example.com"},
		{name: "already normalized", input: "bob@example.com", want: "bob@example.com"},
		{name: "blank input", input: "   ", wantErr: ErrEmailRequired},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NormalizeEmail(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NormalizeEmail() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NormalizeEmail() = %q, want %q", got, tt.want)
			}
		})
	}
}
