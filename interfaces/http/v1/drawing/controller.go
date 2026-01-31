package drawing

import (
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/drawing"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config     config.Config
	Logger     zerolog.Logger
	FiberOapi  *fiberoapi.OApiGroup
	DrawingApp *drawing.DrawingApp
}
