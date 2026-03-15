package webhook

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type WebhookApplication struct {
	Config              config.Config
	Logger              zerolog.Logger
	WebhookPers         domain.WebhookPers
	WebhookDeliveryPers domain.WebhookDeliveryPers
}

func NewWebhookApplication(config config.Config, logger zerolog.Logger, webhookPers domain.WebhookPers, webhookDeliveryPers domain.WebhookDeliveryPers) *WebhookApplication {
	return &WebhookApplication{
		Config:              config,
		Logger:              logger,
		WebhookPers:         webhookPers,
		WebhookDeliveryPers: webhookDeliveryPers,
	}
}

// WebhookApp is a type alias for backward compatibility
type WebhookApp = WebhookApplication
