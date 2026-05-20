package groupbycategory

import (
	"reflect"
	"testing"
)

func TestGroupByCategory(t *testing.T) {
	t.Parallel()

	entries := []Entry{
		{Category: "food", Amount: 12},
		{Category: "food", Amount: 8},
		{Category: " transport ", Amount: 15},
		{Category: "", Amount: 7},
		{Category: "   ", Amount: 3},
	}

	want := map[string]int{
		"food":          20,
		"transport":     15,
		DefaultCategory: 10,
	}

	got := GroupByCategory(entries)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GroupByCategory() = %#v, want %#v", got, want)
	}
}
