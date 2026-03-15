package webhook

import (
	"fmt"

	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/application/webhook/dto"
)

func (app *WebhookApplication) DeleteWebhook(input dto.DeleteWebhookInput) error {
	webhook, err := app.WebhookPers.GetById(input.WebhookId)
	if err != nil {
		return fmt.Errorf("webhook not found: %w", err)
	}

	if webhook.UserId != input.UserId {
		return apperrors.ErrAccessDenied
	}

	if err := app.WebhookPers.Delete(input.WebhookId); err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}
