package monthsnapshot

import "time"

type Entry struct {
	OccurredAt time.Time
	Amount     int
}

type Snapshot struct {
	CurrentTotal  int
	PreviousTotal int
	Delta         int
	DeltaPct      *int
}

func BuildMonthSnapshot(now time.Time, entries []Entry) Snapshot {
	panic("TODO")
}
