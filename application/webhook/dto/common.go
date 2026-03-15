package dto

import "time"

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

type DeliveryItem struct {
	Id         string
	Event      string
	StatusCode int
	Success    bool
	Duration   int
	CreatedAt  time.Time
}
