package parsedate

import (
	"errors"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr error
	}{
		{
			name:  "rfc3339",
			input: "2026-05-20T14:30:00+02:00",
			want:  time.Date(2026, time.May, 20, 12, 30, 0, 0, time.UTC),
		},
		{
			name:  "plain date",
			input: "2026-05-20",
			want:  time.Date(2026, time.May, 20, 0, 0, 0, 0, time.UTC),
		},
		{name: "blank", input: "   ", wantErr: ErrDateRequired},
		{name: "invalid", input: "wat", wantErr: ErrInvalidDate},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseDate(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ParseDate() error = %v, want %v", err, tt.wantErr)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("ParseDate() = %s, want %s", got.Format(time.RFC3339), tt.want.Format(time.RFC3339))
			}
		})
	}
}
