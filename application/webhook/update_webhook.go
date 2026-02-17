package webhook

import (
	"fmt"
	"time"

	"github.com/labbs/nexo/application/webhook/dto"
	"github.com/labbs/nexo/domain"
)

func (app *WebhookApplication) UpdateWebhook(input dto.UpdateWebhookInput) error {
	webhook, err := app.WebhookPers.GetById(input.WebhookId)
	if err != nil {
		return fmt.Errorf("webhook not found: %w", err)
	}

	if webhook.UserId != input.UserId {
		return fmt.Errorf("access denied")
	}

	if input.Name != nil {
		webhook.Name = *input.Name
	}

	if input.Url != nil {
		webhook.Url = *input.Url
	}

	if input.Events != nil {
		webhook.Events = domain.JSONB{
			"events": *input.Events,
		}
	}

	if input.Active != nil {
		webhook.Active = *input.Active
	}

	webhook.UpdatedAt = time.Now()

	if err := app.WebhookPers.Update(webhook); err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	return nil
}
