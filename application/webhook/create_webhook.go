package webhook

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/webhook/dto"
	"github.com/labbs/nexo/domain"
)

func (app *WebhookApplication) CreateWebhook(input dto.CreateWebhookInput) (*dto.CreateWebhookOutput, error) {
	secret, err := generateSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	webhook := &domain.Webhook{
		Id:      uuid.New().String(),
		UserId:  input.UserId,
		SpaceId: input.SpaceId,
		Name:    input.Name,
		Url:     input.Url,
		Secret:  secret,
		Events: domain.JSONB{
			"events": input.Events,
		},
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := app.WebhookPers.Create(webhook); err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return &dto.CreateWebhookOutput{
		Id:     webhook.Id,
		Name:   webhook.Name,
		Url:    webhook.Url,
		Secret: secret,
		Events: input.Events,
		Active: webhook.Active,
	}, nil
}
