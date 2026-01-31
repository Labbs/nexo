package database

import (
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/database"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config      config.Config
	Logger      zerolog.Logger
	FiberOapi   *fiberoapi.OApiGroup
	DatabaseApp *database.DatabaseApp
}
