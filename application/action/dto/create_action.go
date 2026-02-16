package dto

import "time"

type CreateActionInput struct {
	UserId        string
	SpaceId       *string
	DatabaseId    *string
	Name          string
	Description   string
	TriggerType   string
	TriggerConfig map[string]any
	Steps         []ActionStep
}

type CreateActionOutput struct {
	Id          string
	Name        string
	Description string
	TriggerType string
	Active      bool
	CreatedAt   time.Time
}
