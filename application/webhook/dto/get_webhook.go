package dto

import "time"

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
