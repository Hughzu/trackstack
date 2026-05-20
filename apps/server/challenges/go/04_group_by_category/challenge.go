package groupbycategory

const DefaultCategory = "uncategorized"

type Entry struct {
	Category string
	Amount   int
}

func GroupByCategory(entries []Entry) map[string]int {
	panic("TODO")
}
