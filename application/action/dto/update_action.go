package dto

type UpdateActionInput struct {
	UserId        string
	ActionId      string
	Name          *string
	Description   *string
	TriggerType   *string
	TriggerConfig map[string]any
	Steps         *[]ActionStep
	Active        *bool
}
