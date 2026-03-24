package domain

type NutritionSummary struct {
	Calories     int `json:"calories"`
	ProteinGrams int `json:"proteinGrams"`
	CarbGrams    int `json:"carbGrams"`
	FatGrams     int `json:"fatGrams"`
}

type DashboardNutritionSummary struct {
	ConsumedCalories int    `json:"consumedCalories"`
	ProteinGrams     int    `json:"proteinGrams"`
	CarbGrams        int    `json:"carbGrams"`
	FatGrams         int    `json:"fatGrams"`
	Target           Target `json:"target"`
}

type Dashboard struct {
	Summary     DashboardNutritionSummary `json:"summary"`
	Logs        []Log                     `json:"logs"`
	RecentMeals []Log                     `json:"recentMeals"`
}
