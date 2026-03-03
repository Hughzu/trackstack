package calories

type Log struct {
	ID       string  `json:"id"`
	UserID   string  `json:"userId"`
	DateTime string  `json:"dateTime"`
	Calories int     `json:"calories"`
	ProteinG int     `json:"proteinG"`
	CarbsG   *int    `json:"carbsG,omitempty"`
	FatG     *int    `json:"fatG,omitempty"`
	Title    *string `json:"title,omitempty"`
}

type Target struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	TargetKcal     int    `json:"targetKcal"`
	TargetProteinG int    `json:"targetProteinG"`
	TargetCarbsG   *int   `json:"targetCarbsG,omitempty"`
	TargetFatG     *int   `json:"targetFatG,omitempty"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type GetTargetRequest struct {
	UserID string
}

type UpdateTargetRequest struct {
	UserID         string
	TargetKcal     *int
	TargetProteinG *int
	TargetCarbsG   *int
	TargetFatG     *int
}

type AddLogRequest struct {
	UserID   string
	Calories *int
	ProteinG *int
	CarbsG   *int
	FatG     *int
	Title    *string
	Date     *string
	Time     *string
}

type DeleteLogRequest struct {
	UserID string
	ID     string
}
