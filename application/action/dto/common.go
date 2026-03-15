package dto

import "time"

// Step configuration
type ActionStep struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
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

type RunItem struct {
	Id        string
	Success   bool
	Error     string
	Duration  int
	CreatedAt time.Time
}
