package domain

type DashboardBalance struct {
	Remaining float64 `json:"remaining"`
	Income    float64 `json:"income"`
}

type DashboardSpent struct {
	Fund   float64 `json:"fund"`
	Fun    float64 `json:"fun"`
	Future float64 `json:"future"`
}

type DashboardBudget struct {
	Fund   int `json:"fund"`
	Fun    int `json:"fun"`
	Future int `json:"future"`
}

type DashboardRatio struct {
	Percent    int     `json:"percent"`
	CategoryID string  `json:"categoryId"`
	Label      string  `json:"label"`
	Value      float64 `json:"value"`
	Budget     int     `json:"budget"`
	Target     int     `json:"target"`
	Over       bool    `json:"over"`
}

type Dashboard struct {
	PeriodKey          string           `json:"periodKey"`
	Balance            DashboardBalance `json:"balance"`
	Spent              DashboardSpent   `json:"spent"`
	Budget             DashboardBudget  `json:"budget"`
	Ratios             []DashboardRatio `json:"ratios"`
	PendingObligations []ChecklistItem  `json:"pendingObligations"`
	History            []Entry          `json:"history"`
}
