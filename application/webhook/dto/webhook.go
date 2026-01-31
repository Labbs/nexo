package dto

import "time"

// Create webhook
type CreateWebhookInput struct {
	UserId  string
	SpaceId *string
	Name    string
	Url     string
	Events  []string
}

type CreateWebhookOutput struct {
	Id     string
	Name   string
	Url    string
	Secret string
	Events []string
	Active bool
}

// List webhooks
type ListWebhooksInput struct {
	UserId string
}

type WebhookItem struct {
	Id           string
	Name         string
	Url          string
	SpaceId      *string
	SpaceName    *string
	Events       []string
	Active       bool
	LastError    string
	LastErrorAt  *time.Time
	SuccessCount int
	FailureCount int
	CreatedAt    time.Time
}

type ListWebhooksOutput struct {
	Webhooks []WebhookItem
}

// Get webhook
type GetWebhookInput struct {
	UserId    string
	WebhookId string
}

type GetWebhookOutput struct {
	Id           string
	Name         string
	Url          string
	Secret       string
	SpaceId      *string
	SpaceName    *string
	Events       []string
	Active       bool
	LastError    string
	LastErrorAt  *time.Time
	SuccessCount int
	FailureCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Update webhook
type UpdateWebhookInput struct {
	UserId    string
	WebhookId string
	Name      *string
	Url       *string
	Events    *[]string
	Active    *bool
}

// Delete webhook
type DeleteWebhookInput struct {
	UserId    string
	WebhookId string
}

// Get deliveries
type GetDeliveriesInput struct {
	UserId    string
	WebhookId string
	Limit     int
}

type DeliveryItem struct {
	Id         string
	Event      string
	StatusCode int
	Success    bool
	Duration   int
	CreatedAt  time.Time
}

type GetDeliveriesOutput struct {
	Deliveries []DeliveryItem
}

// Trigger webhook (internal)
type TriggerWebhookInput struct {
	Event   string
	SpaceId *string
	Payload map[string]interface{}
}
