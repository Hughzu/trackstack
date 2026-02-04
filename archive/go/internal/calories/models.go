package calories

import "time"

// Meal represents a single meal entry
type Meal struct {
	ID       string
	UserID   string
	Name     string
	Calories int
	Protein  int
	Carbs    int
	Fat      int
	LoggedAt time.Time
}

// UserTargets represents a user's daily nutritional goals
type UserTargets struct {
	UserID        string
	CalorieTarget int
	ProteinTarget int
	UpdatedAt     time.Time
}

// DailySummary represents the dashboard view model
type DailySummary struct {
	Calories    MetricStatus
	Protein     MetricStatus
	RecentMeals []Meal
}

// MetricStatus contains calculated status for a single metric
type MetricStatus struct {
	Consumed   int    // Total consumed today
	Target     int    // Daily target
	Remaining  int    // Target - Consumed
	Percentage int    // Percentage consumed (0-100+)
	Status     string // "safe" | "warn" | "danger"
}
