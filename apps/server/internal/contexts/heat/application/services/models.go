package services

import "github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"

type Refill = domain.Refill

type ListRefillsRequest struct {
	UserID string
	From   string
	To     string
}

type CreateRefillRequest struct {
	UserID      string
	Date        string
	WeightKg    float64
	Bags        int
	Temperature *float64
}

type CreateRefillInput struct {
	Date        string
	WeightKg    float64
	Bags        int
	Temperature *float64
	Season      *string
}

type DeleteRefillRequest struct {
	UserID string
	ID     string
}

type SeasonSnapshot struct {
	SeasonLabel      string `json:"seasonLabel"`
	SeasonToDate     int    `json:"seasonToDate"`
	LastSeasonToDate int    `json:"lastSeasonToDate"`
	Delta            int    `json:"delta"`
	DeltaPct         *int   `json:"deltaPct"`
}

type DashboardViewModel struct {
	DaysSinceRefill int            `json:"daysSinceRefill"`
	SeasonSnapshot  SeasonSnapshot `json:"seasonSnapshot"`
	History         []Refill       `json:"history"`
}

type GetDashboardRequest struct {
	UserID string
	Page   int
	Limit  int
}
