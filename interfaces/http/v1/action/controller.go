package action

import (
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/action"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config    config.Config
	Logger    zerolog.Logger
	FiberOapi *fiberoapi.OApiGroup
	ActionApp *action.ActionApp
}
