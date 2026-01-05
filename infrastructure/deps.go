package infrastructure

import (
	"github.com/labbs/nexo/application/action"
	"github.com/labbs/nexo/application/apikey"
	"github.com/labbs/nexo/application/auth"
	databaseApp "github.com/labbs/nexo/application/database"
	"github.com/labbs/nexo/application/document"
	"github.com/labbs/nexo/application/group"
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

	UserApp     *user.UserApp
	SessionApp  *session.SessionApp
	AuthApp     *auth.AuthApp
	SpaceApp    *space.SpaceApp
	DocumentApp *document.DocumentApp
	ApiKeyApp   *apikey.ApiKeyApp
	WebhookApp  *webhook.WebhookApp
	DatabaseApp *databaseApp.DatabaseApp
	ActionApp   *action.ActionApp
	GroupApp    *group.GroupApp
}
