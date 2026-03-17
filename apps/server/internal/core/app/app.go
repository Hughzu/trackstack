package app

import bootstrap "github.com/Hughzu/trackstack/apps/server/internal/app/bootstrap"

type Runtime = bootstrap.Runtime

func NewRuntime() (*Runtime, error) {
	return bootstrap.NewRuntime()
}
