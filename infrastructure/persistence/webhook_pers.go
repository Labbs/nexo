package persistence

import (
	"time"

	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type webhookPers struct {
	db *gorm.DB
}

func NewWebhookPers(db *gorm.DB) *webhookPers {
	return &webhookPers{db: db}
}

func (p *webhookPers) Create(webhook *domain.Webhook) error {
	return p.db.Create(webhook).Error
}

func (p *webhookPers) GetById(id string) (*domain.Webhook, error) {
	var webhook domain.Webhook
	err := p.db.
		Preload("User").
		Preload("Space").
		Where("id = ?", id).
		First(&webhook).Error
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (p *webhookPers) GetByUserId(userId string) ([]domain.Webhook, error) {
	var webhooks []domain.Webhook
	err := p.db.
		Preload("Space").
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Find(&webhooks).Error
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (p *webhookPers) GetActiveByEvent(event domain.WebhookEvent, spaceId *string) ([]domain.Webhook, error) {
	var webhooks []domain.Webhook
	query := p.db.Where("active = ?", true)

	// Filter by space if provided, or get global webhooks (no space)
	if spaceId != nil {
		query = query.Where("space_id = ? OR space_id IS NULL", *spaceId)
	} else {
		query = query.Where("space_id IS NULL")
	}

	err := query.Find(&webhooks).Error
	if err != nil {
		return nil, err
	}

	// Filter by event in application code since Events is JSONB
	var filtered []domain.Webhook
	for _, w := range webhooks {
		if w.HasEvent(event) {
			filtered = append(filtered, w)
		}
	}

	return filtered, nil
}

func (p *webhookPers) Update(webhook *domain.Webhook) error {
	return p.db.Save(webhook).Error
}

func (p *webhookPers) Delete(id string) error {
	return p.db.Where("id = ?", id).Delete(&domain.Webhook{}).Error
}

func (p *webhookPers) IncrementSuccess(id string) error {
	return p.db.Model(&domain.Webhook{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"success_count": gorm.Expr("success_count + 1"),
			"last_error":    "",
			"last_error_at": nil,
		}).Error
}

func (p *webhookPers) RecordFailure(id string, errorMsg string) error {
	now := time.Now()
	return p.db.Model(&domain.Webhook{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"failure_count": gorm.Expr("failure_count + 1"),
			"last_error":    errorMsg,
			"last_error_at": &now,
		}).Error
}

// WebhookDelivery persistence
type webhookDeliveryPers struct {
	db *gorm.DB
}

func NewWebhookDeliveryPers(db *gorm.DB) *webhookDeliveryPers {
	return &webhookDeliveryPers{db: db}
}

func (p *webhookDeliveryPers) Create(delivery *domain.WebhookDelivery) error {
	return p.db.Create(delivery).Error
}

func (p *webhookDeliveryPers) GetByWebhookId(webhookId string, limit int) ([]domain.WebhookDelivery, error) {
	var deliveries []domain.WebhookDelivery
	err := p.db.
		Where("webhook_id = ?", webhookId).
		Order("created_at DESC").
		Limit(limit).
		Find(&deliveries).Error
	if err != nil {
		return nil, err
	}
	return deliveries, nil
}
