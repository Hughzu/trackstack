package httptransport

import (
	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
)

type Handlers struct {
	Heat     *HeatHandler
	Expenses *ExpensesHandler
	Calories *CaloriesHandler
}

type Deps struct {
	HeatService     *heat.Service
	ExpensesService *expenses.Service
	CaloriesService *calories.Service
	HardcodedUserID string
}

func NewHandlers(deps Deps) Handlers {
	return Handlers{
		Heat: &HeatHandler{
			svc:    deps.HeatService,
			userID: deps.HardcodedUserID,
		},
		Expenses: &ExpensesHandler{
			svc:    deps.ExpensesService,
			userID: deps.HardcodedUserID,
		},
		Calories: &CaloriesHandler{
			svc:    deps.CaloriesService,
			userID: deps.HardcodedUserID,
		},
	}
}
