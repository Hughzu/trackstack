package httptransport

import "github.com/Hughzu/trackstack/apps/server/internal/modules/heat"

type Handlers struct {
	Heat *HeatHandler
}

type Deps struct {
	HeatService     *heat.Service
	HardcodedUserID string
}

func NewHandlers(deps Deps) Handlers {
	return Handlers{
		Heat: &HeatHandler{
			svc:    deps.HeatService,
			userID: deps.HardcodedUserID,
		},
	}
}
