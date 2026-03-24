package domain

type Settings struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userId"`
	Income      float64 `json:"income"`
	RatioFund   int     `json:"ratioFund"`
	RatioFun    int     `json:"ratioFun"`
	RatioFuture int     `json:"ratioFuture"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type SettingsView struct {
	Settings  Settings   `json:"settings"`
	Checklist []Template `json:"checklist"`
	Recurring []Template `json:"recurring"`
}
