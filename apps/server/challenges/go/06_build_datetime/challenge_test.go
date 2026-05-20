package builddatetime

import (
	"errors"
	"testing"
	"time"
)

func TestBuildRFC3339DateTime(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.May, 20, 9, 45, 30, 0, time.UTC)

	tests := []struct {
		name      string
		dateValue string
		timeValue string
		want      string
		wantErr   error
	}{
		{name: "full values", dateValue: "2026-05-21", timeValue: "14:05", want: "2026-05-21T14:05:00Z"},
		{name: "default date", dateValue: "", timeValue: "08:15", want: "2026-05-20T08:15:00Z"},
		{name: "default time", dateValue: "2026-05-21", timeValue: "", want: "2026-05-21T09:45:00Z"},
		{name: "default both", dateValue: "", timeValue: "", want: "2026-05-20T09:45:00Z"},
		{name: "invalid date", dateValue: "2026-13-99", timeValue: "09:45", wantErr: ErrInvalidDateTime},
		{name: "invalid time", dateValue: "2026-05-21", timeValue: "25:61", wantErr: ErrInvalidDateTime},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := BuildRFC3339DateTime(tt.dateValue, tt.timeValue, now)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("BuildRFC3339DateTime() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("BuildRFC3339DateTime() = %q, want %q", got, tt.want)
			}
		})
	}
}
