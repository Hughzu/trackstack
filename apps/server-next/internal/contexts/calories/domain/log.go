package domain

type Log struct {
	ID           string  `json:"id"`
	UserID       string  `json:"userId"`
	DateTime     string  `json:"dateTime"`
	Calories     int     `json:"calories"`
	ProteinGrams int     `json:"proteinGrams"`
	CarbGrams    *int    `json:"carbGrams,omitempty"`
	FatGrams     *int    `json:"fatGrams,omitempty"`
	Title        *string `json:"title,omitempty"`
}
