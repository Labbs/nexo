package infrastructure

import (
	"github.com/labbs/nexo/application/action"
	"github.com/labbs/nexo/application/apikey"
	"github.com/labbs/nexo/application/auth"
	databaseApp "github.com/labbs/nexo/application/database"
	"github.com/labbs/nexo/application/document"
	"github.com/labbs/nexo/application/drawing"
	"github.com/labbs/nexo/application/favorite"
	"github.com/labbs/nexo/application/group"
	"github.com/labbs/nexo/application/permission"
	"github.com/labbs/nexo/application/session"
	"github.com/labbs/nexo/application/space"
	"github.com/labbs/nexo/application/user"
	"github.com/labbs/nexo/application/webhook"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/labbs/nexo/infrastructure/cronscheduler"
	"github.com/labbs/nexo/infrastructure/database"
	"github.com/labbs/nexo/infrastructure/http"
	"github.com/rs/zerolog"
)

type Deps struct {
	Config        config.Config
	Logger        zerolog.Logger
	Http          http.Config
	CronScheduler cronscheduler.Config
	Database      database.Config

	UserApp       *user.UserApplication
	SessionApp    *session.SessionApplication
	AuthApp       *auth.AuthApplication
	SpaceApp      *space.SpaceApplication
	DocumentApp   *document.DocumentApplication
	ApiKeyApp     *apikey.ApiKeyApplication
	WebhookApp    *webhook.WebhookApplication
	DatabaseApp   *databaseApp.DatabaseApplication
	DrawingApp    *drawing.DrawingApplication
	ActionApp     *action.ActionApplication
	GroupApp      *group.GroupApp
	FavoriteApp   *favorite.FavoriteApp
	PermissionApp *permission.PermissionApplication
}
