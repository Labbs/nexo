package dto

type ExecuteActionInput struct {
	TriggerType string
	SpaceId     *string
	DatabaseId  *string
	TriggerData map[string]any
}
