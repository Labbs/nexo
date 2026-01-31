package dtos

import "time"

// Request DTOs

type EmptyRequest struct{}

type ActionStep struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type CreateActionRequest struct {
	SpaceId       *string                `json:"space_id,omitempty"`
	DatabaseId    *string                `json:"database_id,omitempty"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description,omitempty"`
	TriggerType   string                 `json:"trigger_type"`
	TriggerConfig map[string]interface{} `json:"trigger_config,omitempty"`
	Steps         []ActionStep           `json:"steps"`
}

type GetActionRequest struct {
	ActionId string `path:"action_id"`
}

type UpdateActionRequest struct {
	ActionId      string                 `path:"action_id"`
	Name          *string                `json:"name,omitempty"`
	Description   *string                `json:"description,omitempty"`
	TriggerType   *string                `json:"trigger_type,omitempty"`
	TriggerConfig map[string]interface{} `json:"trigger_config,omitempty"`
	Steps         *[]ActionStep          `json:"steps,omitempty"`
	Active        *bool                  `json:"active,omitempty"`
}

type DeleteActionRequest struct {
	ActionId string `path:"action_id"`
}

type GetRunsRequest struct {
	ActionId string `path:"action_id"`
	Limit    int    `query:"limit"`
}

// Response DTOs

type MessageResponse struct {
	Message string `json:"message"`
}

type CreateActionResponse struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TriggerType string    `json:"trigger_type"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

type ActionItem struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	SpaceId      *string    `json:"space_id,omitempty"`
	SpaceName    *string    `json:"space_name,omitempty"`
	DatabaseId   *string    `json:"database_id,omitempty"`
	TriggerType  string     `json:"trigger_type"`
	Active       bool       `json:"active"`
	LastRunAt    *time.Time `json:"last_run_at,omitempty"`
	LastError    string     `json:"last_error,omitempty"`
	RunCount     int        `json:"run_count"`
	SuccessCount int        `json:"success_count"`
	FailureCount int        `json:"failure_count"`
	CreatedAt    time.Time  `json:"created_at"`
}

type ListActionsResponse struct {
	Actions []ActionItem `json:"actions"`
}

type GetActionResponse struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	SpaceId       *string                `json:"space_id,omitempty"`
	SpaceName     *string                `json:"space_name,omitempty"`
	DatabaseId    *string                `json:"database_id,omitempty"`
	TriggerType   string                 `json:"trigger_type"`
	TriggerConfig map[string]interface{} `json:"trigger_config,omitempty"`
	Steps         []ActionStep           `json:"steps"`
	Active        bool                   `json:"active"`
	LastRunAt     *time.Time             `json:"last_run_at,omitempty"`
	LastError     string                 `json:"last_error,omitempty"`
	RunCount      int                    `json:"run_count"`
	SuccessCount  int                    `json:"success_count"`
	FailureCount  int                    `json:"failure_count"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type RunItem struct {
	Id        string    `json:"id"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	Duration  int       `json:"duration_ms"`
	CreatedAt time.Time `json:"created_at"`
}

type GetRunsResponse struct {
	Runs []RunItem `json:"runs"`
}

// Available triggers and steps
type AvailableTriggersResponse struct {
	Triggers []TriggerInfo `json:"triggers"`
}

type TriggerInfo struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type AvailableStepsResponse struct {
	Steps []StepInfo `json:"steps"`
}

type StepInfo struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
