package webhook

import (
	"fmt"

	"github.com/labbs/nexo/application/webhook/dto"
)

func (app *WebhookApplication) ListWebhooks(input dto.ListWebhooksInput) (*dto.ListWebhooksOutput, error) {
	webhooks, err := app.WebhookPers.GetByUserId(input.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}

	output := &dto.ListWebhooksOutput{
		Webhooks: make([]dto.WebhookItem, len(webhooks)),
	}

	for i, w := range webhooks {
		var events []string
		if w.Events != nil {
			if e, ok := w.Events["events"].([]interface{}); ok {
				for _, ev := range e {
					if str, ok := ev.(string); ok {
						events = append(events, str)
					}
				}
			}
		}

		item := dto.WebhookItem{
			Id:           w.Id,
			Name:         w.Name,
			Url:          w.Url,
			SpaceId:      w.SpaceId,
			Events:       events,
			Active:       w.Active,
			LastError:    w.LastError,
			LastErrorAt:  w.LastErrorAt,
			SuccessCount: w.SuccessCount,
			FailureCount: w.FailureCount,
			CreatedAt:    w.CreatedAt,
		}

		if w.Space != nil {
			item.SpaceName = &w.Space.Name
		}

		output.Webhooks[i] = item
	}

	return output, nil
}
