package app

import (
	"github.com/labbs/nexo/infrastructure"
	"github.com/labbs/nexo/infrastructure/static"
)

func SetupRouterApp(deps infrastructure.Deps) {
	deps.Logger.Info().Str("component", "http.router.app").Msg("Setting up application routes")

	static.NewStatic(deps.Http.Fiber)
}
