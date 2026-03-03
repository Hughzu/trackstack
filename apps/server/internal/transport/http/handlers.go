package httptransport

import (
	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/users"
)

type Handlers struct {
	Heat     *HeatHandler
	Expenses *ExpensesHandler
	Calories *CaloriesHandler
	Auth     *AuthHandler
}

type Deps struct {
	HeatService        *heat.Service
	ExpensesService    *expenses.Service
	CaloriesService    *calories.Service
	UsersService       *users.Service
	AuthService        *auth.Service
	HardcodedUserID    string
	AuthCookieName     string
	AuthCookieSecure   bool
	AuthCookieSameSite string
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
		Auth: &AuthHandler{
			authService:       deps.AuthService,
			usersService:      deps.UsersService,
			cookieName:        deps.AuthCookieName,
			cookieSecure:      deps.AuthCookieSecure,
			cookieSameSiteRaw: deps.AuthCookieSameSite,
		},
	}
}
