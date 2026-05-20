package optionalint

import "testing"

func TestParseOptionalInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  *int
	}{
		{name: "empty string", input: "", want: nil},
		{name: "spaces only", input: "   ", want: nil},
		{name: "positive int", input: "42", want: intPtr(42)},
		{name: "negative int", input: " -7 ", want: intPtr(-7)},
		{name: "invalid string", input: "abc", want: nil},
		{name: "decimal string", input: "3.14", want: nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ParseOptionalInt(tt.input)
			if !sameIntPtr(got, tt.want) {
				t.Fatalf("ParseOptionalInt(%q) = %v, want %v", tt.input, ptrValue(got), ptrValue(tt.want))
			}
		})
	}
}

func intPtr(value int) *int {
	return &value
}

func sameIntPtr(left *int, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func ptrValue(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}
