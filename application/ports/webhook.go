package ports

import (
	"github.com/labbs/nexo/application/webhook/dto"
)

type WebhookPort interface {
	CreateWebhook(input dto.CreateWebhookInput) (*dto.CreateWebhookOutput, error)
	ListWebhooks(input dto.ListWebhooksInput) (*dto.ListWebhooksOutput, error)
	GetWebhook(input dto.GetWebhookInput) (*dto.GetWebhookOutput, error)
	UpdateWebhook(input dto.UpdateWebhookInput) error
	DeleteWebhook(input dto.DeleteWebhookInput) error
	GetDeliveries(input dto.GetDeliveriesInput) (*dto.GetDeliveriesOutput, error)
	TriggerWebhooks(input dto.TriggerWebhookInput)
}
