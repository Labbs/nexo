package webhook

import (
	"fmt"

	"github.com/labbs/nexo/application/webhook/dto"
)

func (app *WebhookApplication) GetWebhook(input dto.GetWebhookInput) (*dto.GetWebhookOutput, error) {
	webhook, err := app.WebhookPers.GetById(input.WebhookId)
	if err != nil {
		return nil, fmt.Errorf("webhook not found: %w", err)
	}

	if webhook.UserId != input.UserId {
		return nil, fmt.Errorf("access denied")
	}

	var events []string
	if webhook.Events != nil {
		if e, ok := webhook.Events["events"].([]interface{}); ok {
			for _, ev := range e {
				if str, ok := ev.(string); ok {
					events = append(events, str)
				}
			}
		}
	}

	output := &dto.GetWebhookOutput{
		Id:           webhook.Id,
		Name:         webhook.Name,
		Url:          webhook.Url,
		Secret:       webhook.Secret,
		SpaceId:      webhook.SpaceId,
		Events:       events,
		Active:       webhook.Active,
		LastError:    webhook.LastError,
		LastErrorAt:  webhook.LastErrorAt,
		SuccessCount: webhook.SuccessCount,
		FailureCount: webhook.FailureCount,
		CreatedAt:    webhook.CreatedAt,
		UpdatedAt:    webhook.UpdatedAt,
	}

	if webhook.Space != nil {
		output.SpaceName = &webhook.Space.Name
	}

	return output, nil
}
