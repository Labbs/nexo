package dto

import "time"

// Step configuration
type ActionStep struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// Create action
type CreateActionInput struct {
	UserId        string
	SpaceId       *string
	DatabaseId    *string
	Name          string
	Description   string
	TriggerType   string
	TriggerConfig map[string]interface{}
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

// List actions
type ListActionsInput struct {
	UserId string
}

type ActionItem struct {
	Id           string
	Name         string
	Description  string
	SpaceId      *string
	SpaceName    *string
	DatabaseId   *string
	TriggerType  string
	Active       bool
	LastRunAt    *time.Time
	LastError    string
	RunCount     int
	SuccessCount int
	FailureCount int
	CreatedAt    time.Time
}

type ListActionsOutput struct {
	Actions []ActionItem
}

// Get action
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
	TriggerConfig map[string]interface{}
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

// Update action
type UpdateActionInput struct {
	UserId        string
	ActionId      string
	Name          *string
	Description   *string
	TriggerType   *string
	TriggerConfig map[string]interface{}
	Steps         *[]ActionStep
	Active        *bool
}

// Delete action
type DeleteActionInput struct {
	UserId   string
	ActionId string
}

// Get runs
type GetRunsInput struct {
	UserId   string
	ActionId string
	Limit    int
}

type RunItem struct {
	Id        string
	Success   bool
	Error     string
	Duration  int
	CreatedAt time.Time
}

type GetRunsOutput struct {
	Runs []RunItem
}

// Execute action (internal)
type ExecuteActionInput struct {
	TriggerType string
	SpaceId     *string
	DatabaseId  *string
	TriggerData map[string]interface{}
}
