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
		Config:          deps.Config,
		Logger:          deps.Logger,
		FiberOapi:       grp.Group("/auth"),
		AuthApplication: deps.AuthApplication,
	}
	auth.SetupAuthRouter(authCtrl)

	userCtrl := user.Controller{
		Config:              deps.Config,
		Logger:              deps.Logger,
		FiberOapi:           grp.Group("/user"),
		UserApplication:     deps.UserApplication,
		SpaceApplication:    deps.SpaceApplication,
		FavoriteApplication: deps.FavoriteApplication,
		OAuthProviderPers:   deps.OAuthProviderPers,
	}
	user.SetupUserRouter(userCtrl)

	spaceCtrl := space.Controller{
		Config:                deps.Config,
		Logger:                deps.Logger,
		FiberOapi:             grp.Group("/space"),
		SpaceApplication:      deps.SpaceApplication,
		PermissionApplication: deps.PermissionApplication,
	}
	space.SetupSpaceRouter(spaceCtrl)

	documentCtrl := document.Controller{
		Config:                deps.Config,
		Logger:                deps.Logger,
		FiberOapi:             grp.Group("/document"),
		SpaceApplication:      deps.SpaceApplication,
		DocumentApplication:   deps.DocumentApplication,
		PermissionApplication: deps.PermissionApplication,
	}
	document.SetupDocumentRouter(documentCtrl)

	apiKeyCtrl := apikey.Controller{
		Config:            deps.Config,
		Logger:            deps.Logger,
		FiberOapi:         grp.Group("/apikeys"),
		ApiKeyApplication: deps.ApiKeyApplication,
	}
	apikey.SetupApiKeyRouter(apiKeyCtrl)

	webhookCtrl := webhook.Controller{
		Config:             deps.Config,
		Logger:             deps.Logger,
		FiberOapi:          grp.Group("/webhooks"),
		WebhookApplication: deps.WebhookApplication,
	}
	webhook.SetupWebhookRouter(webhookCtrl)

	databaseCtrl := database.Controller{
		Config:                deps.Config,
		Logger:                deps.Logger,
		FiberOapi:             grp.Group("/databases"),
		DatabaseApplication:   deps.DatabaseApplication,
		PermissionApplication: deps.PermissionApplication,
	}
	database.SetupDatabaseRouter(databaseCtrl)

	drawingCtrl := drawing.Controller{
		Config:                deps.Config,
		Logger:                deps.Logger,
		FiberOapi:             grp.Group("/drawings"),
		DrawingApplication:    deps.DrawingApplication,
		PermissionApplication: deps.PermissionApplication,
	}
	drawing.SetupDrawingRouter(drawingCtrl)

	actionCtrl := action.Controller{
		Config:            deps.Config,
		Logger:            deps.Logger,
		FiberOapi:         grp.Group("/actions"),
		ActionApplication: deps.ActionApplication,
	}
	action.SetupActionRouter(actionCtrl)

	adminCtrl := admin.Controller{
		Config:            deps.Config,
		Logger:            deps.Logger,
		FiberOapi:         grp.Group("/admin"),
		UserApplication:   deps.UserApplication,
		SpaceApplication:  deps.SpaceApplication,
		ApiKeyApplication: deps.ApiKeyApplication,
		GroupApplication:  deps.GroupApplication,
		PermissionPers:    deps.PermissionPers,
	}
	admin.SetupAdminRouter(adminCtrl)

	app.SetupRouterApp(deps)
}
