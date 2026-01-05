package dtos

import "time"

// Request DTOs

type EmptyRequest struct{}

type CreateWebhookRequest struct {
	Name    string   `json:"name"`
	Url     string   `json:"url"`
	SpaceId *string  `json:"space_id,omitempty"`
	Events  []string `json:"events"`
}

type UpdateWebhookRequest struct {
	WebhookId string    `path:"webhook_id"`
	Name      *string   `json:"name,omitempty"`
	Url       *string   `json:"url,omitempty"`
	Events    *[]string `json:"events,omitempty"`
	Active    *bool     `json:"active,omitempty"`
}

type DeleteWebhookRequest struct {
	WebhookId string `path:"webhook_id"`
}

type GetWebhookRequest struct {
	WebhookId string `path:"webhook_id"`
}

type GetDeliveriesRequest struct {
	WebhookId string `path:"webhook_id"`
	Limit     int    `query:"limit"`
}

// Response DTOs

type MessageResponse struct {
	Message string `json:"message"`
}

type CreateWebhookResponse struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Url    string   `json:"url"`
	Secret string   `json:"secret"`
	Events []string `json:"events"`
	Active bool     `json:"active"`
}

type WebhookItem struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Url          string     `json:"url"`
	SpaceId      *string    `json:"space_id,omitempty"`
	SpaceName    *string    `json:"space_name,omitempty"`
	Events       []string   `json:"events"`
	Active       bool       `json:"active"`
	LastError    string     `json:"last_error,omitempty"`
	LastErrorAt  *time.Time `json:"last_error_at,omitempty"`
	SuccessCount int        `json:"success_count"`
	FailureCount int        `json:"failure_count"`
	CreatedAt    time.Time  `json:"created_at"`
}

type ListWebhooksResponse struct {
	Webhooks []WebhookItem `json:"webhooks"`
}

type GetWebhookResponse struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Url          string     `json:"url"`
	Secret       string     `json:"secret"`
	SpaceId      *string    `json:"space_id,omitempty"`
	SpaceName    *string    `json:"space_name,omitempty"`
	Events       []string   `json:"events"`
	Active       bool       `json:"active"`
	LastError    string     `json:"last_error,omitempty"`
	LastErrorAt  *time.Time `json:"last_error_at,omitempty"`
	SuccessCount int        `json:"success_count"`
	FailureCount int        `json:"failure_count"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type DeliveryItem struct {
	Id         string    `json:"id"`
	Event      string    `json:"event"`
	StatusCode int       `json:"status_code"`
	Success    bool      `json:"success"`
	Duration   int       `json:"duration_ms"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetDeliveriesResponse struct {
	Deliveries []DeliveryItem `json:"deliveries"`
}

// Available events for reference
type AvailableEventsResponse struct {
	Events []EventInfo `json:"events"`
}

type EventInfo struct {
	Event       string `json:"event"`
	Description string `json:"description"`
}
