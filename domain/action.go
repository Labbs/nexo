package domain

import (
	"time"

	"gorm.io/gorm"
)

// Action represents an automation rule
type Action struct {
	Id string

	UserId string
	User   User `gorm:"foreignKey:UserId;references:Id"`

	// Optional: scope action to a specific space or database
	SpaceId    *string
	Space      *Space `gorm:"foreignKey:SpaceId;references:Id"`
	DatabaseId *string

	Name        string
	Description string

	// Trigger configuration
	TriggerType   ActionTriggerType
	TriggerConfig JSONB // Trigger-specific configuration

	// Action steps to execute
	Steps JSONB // [{type, config}]

	// Status
	Active       bool
	LastRunAt    *time.Time
	LastError    string
	RunCount     int
	SuccessCount int
	FailureCount int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (a *Action) TableName() string {
	return "action"
}

// ActionTriggerType defines when an action should be triggered
type ActionTriggerType string

const (
	// Document triggers
	TriggerDocumentCreated  ActionTriggerType = "document.created"
	TriggerDocumentUpdated  ActionTriggerType = "document.updated"
	TriggerDocumentDeleted  ActionTriggerType = "document.deleted"
	TriggerDocumentMoved    ActionTriggerType = "document.moved"
	TriggerDocumentShared   ActionTriggerType = "document.shared"

	// Database triggers
	TriggerRowCreated       ActionTriggerType = "row.created"
	TriggerRowUpdated       ActionTriggerType = "row.updated"
	TriggerRowDeleted       ActionTriggerType = "row.deleted"
	TriggerPropertyChanged  ActionTriggerType = "property.changed"

	// Comment triggers
	TriggerCommentCreated   ActionTriggerType = "comment.created"
	TriggerCommentResolved  ActionTriggerType = "comment.resolved"

	// Schedule triggers
	TriggerSchedule         ActionTriggerType = "schedule"
)

// ActionStepType defines the types of actions that can be executed
type ActionStepType string

const (
	// Notification actions
	StepSendEmail         ActionStepType = "send_email"
	StepSendSlack         ActionStepType = "send_slack"
	StepSendWebhook       ActionStepType = "send_webhook"

	// Document actions
	StepCreateDocument    ActionStepType = "create_document"
	StepUpdateDocument    ActionStepType = "update_document"
	StepMoveDocument      ActionStepType = "move_document"
	StepDuplicateDocument ActionStepType = "duplicate_document"

	// Database actions
	StepCreateRow         ActionStepType = "create_row"
	StepUpdateRow         ActionStepType = "update_row"
	StepDeleteRow         ActionStepType = "delete_row"
	StepUpdateProperty    ActionStepType = "update_property"

	// Misc actions
	StepAddComment        ActionStepType = "add_comment"
	StepAssignUser        ActionStepType = "assign_user"
	StepSetReminder       ActionStepType = "set_reminder"
)

type ActionPers interface {
	Create(action *Action) error
	GetById(id string) (*Action, error)
	GetByUserId(userId string) ([]Action, error)
	GetActiveByTrigger(triggerType ActionTriggerType, spaceId *string, databaseId *string) ([]Action, error)
	Update(action *Action) error
	Delete(id string) error
	IncrementSuccess(id string) error
	RecordFailure(id string, errorMsg string) error
	UpdateLastRun(id string) error
}

// ActionRun records individual action execution
type ActionRun struct {
	Id string

	ActionId string
	Action   Action `gorm:"foreignKey:ActionId;references:Id"`

	TriggerData JSONB // Data that triggered the action
	StepsResult JSONB // Result of each step
	Success     bool
	Error       string
	Duration    int // milliseconds

	CreatedAt time.Time
}

func (r *ActionRun) TableName() string {
	return "action_run"
}

type ActionRunPers interface {
	Create(run *ActionRun) error
	GetByActionId(actionId string, limit int) ([]ActionRun, error)
}
