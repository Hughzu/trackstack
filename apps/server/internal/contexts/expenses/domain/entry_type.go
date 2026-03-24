package domain

type EntryType string

const (
	EntryTypeManual    EntryType = "manual"
	EntryTypeRecurring EntryType = "recurring"
	EntryTypeChecklist EntryType = "checklist"
)
