package domain

type Target struct {
	ID                 string `json:"id"`
	UserID             string `json:"userId"`
	TargetCalories     int    `json:"targetCalories"`
	TargetProteinGrams int    `json:"targetProteinGrams"`
	TargetCarbGrams    *int   `json:"targetCarbGrams,omitempty"`
	TargetFatGrams     *int   `json:"targetFatGrams,omitempty"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}
