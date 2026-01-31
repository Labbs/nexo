package domain

import (
	"time"

	"gorm.io/gorm"
)

type Webhook struct {
	Id string

	UserId string
	User   User `gorm:"foreignKey:UserId;references:Id"`

	// Optional: scope webhook to a specific space
	SpaceId *string
	Space   *Space `gorm:"foreignKey:SpaceId;references:Id"`

	Name   string
	Url    string
	Secret string // Used for signature verification

	// Events to trigger on
	Events JSONB // ["document.created", "document.updated", etc.]

	// Status
	Active       bool
	LastError    string
	LastErrorAt  *time.Time
	SuccessCount int
	FailureCount int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (w *Webhook) TableName() string {
	return "webhook"
}

// Webhook events
type WebhookEvent string

const (
	WebhookEventDocumentCreated WebhookEvent = "document.created"
	WebhookEventDocumentUpdated WebhookEvent = "document.updated"
	WebhookEventDocumentDeleted WebhookEvent = "document.deleted"
	WebhookEventCommentCreated  WebhookEvent = "comment.created"
	WebhookEventCommentResolved WebhookEvent = "comment.resolved"
	WebhookEventSpaceCreated    WebhookEvent = "space.created"
	WebhookEventSpaceUpdated    WebhookEvent = "space.updated"
)

func (w *Webhook) HasEvent(event WebhookEvent) bool {
	if w.Events == nil {
		return false
	}
	events, ok := w.Events["events"].([]interface{})
	if !ok {
		return false
	}
	for _, e := range events {
		if str, ok := e.(string); ok && str == string(event) {
			return true
		}
	}
	return false
}

type WebhookPers interface {
	Create(webhook *Webhook) error
	GetById(id string) (*Webhook, error)
	GetByUserId(userId string) ([]Webhook, error)
	GetActiveByEvent(event WebhookEvent, spaceId *string) ([]Webhook, error)
	Update(webhook *Webhook) error
	Delete(id string) error
	IncrementSuccess(id string) error
	RecordFailure(id string, errorMsg string) error
}

// WebhookDelivery records individual delivery attempts
type WebhookDelivery struct {
	Id string

	WebhookId string
	Webhook   Webhook `gorm:"foreignKey:WebhookId;references:Id"`

	Event      string
	Payload    JSONB
	StatusCode int
	Response   string
	Duration   int // milliseconds
	Success    bool

	CreatedAt time.Time
}

func (d *WebhookDelivery) TableName() string {
	return "webhook_delivery"
}

type WebhookDeliveryPers interface {
	Create(delivery *WebhookDelivery) error
	GetByWebhookId(webhookId string, limit int) ([]WebhookDelivery, error)
}
