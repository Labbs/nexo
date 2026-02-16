package dto

import "time"

type GetActionInput struct {
	UserId   string
	ActionId string
}

type GetActionOutput struct {
	Id            string
	Name          string
	Description   string
	SpaceId       *string
	SpaceName     *string
	DatabaseId    *string
	TriggerType   string
	TriggerConfig map[string]any
	Steps         []ActionStep
	Active        bool
	LastRunAt     *time.Time
	LastError     string
	RunCount      int
	SuccessCount  int
	FailureCount  int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
