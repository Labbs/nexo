package v1

import (
	"github.com/labbs/nexo/infrastructure"
	"github.com/labbs/nexo/interfaces/http/app"
	"github.com/labbs/nexo/interfaces/http/v1/action"
	"github.com/labbs/nexo/interfaces/http/v1/admin"
	"github.com/labbs/nexo/interfaces/http/v1/apikey"
	"github.com/labbs/nexo/interfaces/http/v1/auth"
	"github.com/labbs/nexo/interfaces/http/v1/database"
	"github.com/labbs/nexo/interfaces/http/v1/document"
	"github.com/labbs/nexo/interfaces/http/v1/drawing"
	"github.com/labbs/nexo/interfaces/http/v1/space"
	"github.com/labbs/nexo/interfaces/http/v1/user"
	"github.com/labbs/nexo/interfaces/http/v1/webhook"
)

func SetupRouterV1(deps infrastructure.Deps) {
	deps.Logger.Info().Str("component", "http.router.v1").Msg("Setting up API v1 routes")
	grp := deps.Http.FiberOapi.Group("/api/v1")

	authCtrl := auth.Controller{
		Config:    deps.Config,
		Logger:    deps.Logger,
		FiberOapi: grp.Group("/auth"),
		AuthApp:   deps.AuthApp,
	}
	auth.SetupAuthRouter(authCtrl)

	userCtrl := user.Controller{
		Config:    deps.Config,
		Logger:    deps.Logger,
		FiberOapi: grp.Group("/user"),
		UserApp:   deps.UserApp,
		SpaceApp:  deps.SpaceApp,
	}
	user.SetupUserRouter(userCtrl)

	spaceCtrl := space.Controller{
		Config:    deps.Config,
		Logger:    deps.Logger,
		FiberOapi: grp.Group("/space"),
		SpaceApp:  deps.SpaceApp,
	}
	space.SetupSpaceRouter(spaceCtrl)

	documentCtrl := document.Controller{
		Config:      deps.Config,
		Logger:      deps.Logger,
		FiberOapi:   grp.Group("/document"),
		SpaceApp:    deps.SpaceApp,
		DocumentApp: deps.DocumentApp,
	}
	document.SetupDocumentRouter(documentCtrl)

	apiKeyCtrl := apikey.Controller{
		Config:    deps.Config,
		Logger:    deps.Logger,
		FiberOapi: grp.Group("/apikeys"),
		ApiKeyApp: deps.ApiKeyApp,
	}
	apikey.SetupApiKeyRouter(apiKeyCtrl)

	webhookCtrl := webhook.Controller{
		Config:     deps.Config,
		Logger:     deps.Logger,
		FiberOapi:  grp.Group("/webhooks"),
		WebhookApp: deps.WebhookApp,
	}
	webhook.SetupWebhookRouter(webhookCtrl)

	databaseCtrl := database.Controller{
		Config:      deps.Config,
		Logger:      deps.Logger,
		FiberOapi:   grp.Group("/databases"),
		DatabaseApp: deps.DatabaseApp,
	}
	database.SetupDatabaseRouter(databaseCtrl)

	drawingCtrl := drawing.Controller{
		Config:     deps.Config,
		Logger:     deps.Logger,
		FiberOapi:  grp.Group("/drawings"),
		DrawingApp: deps.DrawingApp,
	}
	drawing.SetupDrawingRouter(drawingCtrl)

	actionCtrl := action.Controller{
		Config:    deps.Config,
		Logger:    deps.Logger,
		FiberOapi: grp.Group("/actions"),
		ActionApp: deps.ActionApp,
	}
	action.SetupActionRouter(actionCtrl)

	adminCtrl := admin.Controller{
		Config:    deps.Config,
		Logger:    deps.Logger,
		FiberOapi: grp.Group("/admin"),
		UserApp:   deps.UserApp,
		SpaceApp:  deps.SpaceApp,
		ApiKeyApp: deps.ApiKeyApp,
		GroupApp:  deps.GroupApp,
	}
	admin.SetupAdminRouter(adminCtrl)

	app.SetupRouterApp(deps)
}
