package webhook

import (
	"fmt"

	"github.com/labbs/nexo/application/webhook/dto"
)

func (app *WebhookApplication) GetDeliveries(input dto.GetDeliveriesInput) (*dto.GetDeliveriesOutput, error) {
	// Verify ownership
	webhook, err := app.WebhookPers.GetById(input.WebhookId)
	if err != nil {
		return nil, fmt.Errorf("webhook not found: %w", err)
	}

	if webhook.UserId != input.UserId {
		return nil, fmt.Errorf("access denied")
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	deliveries, err := app.WebhookDeliveryPers.GetByWebhookId(input.WebhookId, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliveries: %w", err)
	}

	output := &dto.GetDeliveriesOutput{
		Deliveries: make([]dto.DeliveryItem, len(deliveries)),
	}

	for i, d := range deliveries {
		output.Deliveries[i] = dto.DeliveryItem{
			Id:         d.Id,
			Event:      d.Event,
			StatusCode: d.StatusCode,
			Success:    d.Success,
			Duration:   d.Duration,
			CreatedAt:  d.CreatedAt,
		}
	}

	return output, nil
}
