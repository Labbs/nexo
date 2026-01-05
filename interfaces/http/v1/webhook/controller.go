package webhook

import (
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/webhook"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config     config.Config
	Logger     zerolog.Logger
	FiberOapi  *fiberoapi.OApiGroup
	WebhookApp *webhook.WebhookApp
}
