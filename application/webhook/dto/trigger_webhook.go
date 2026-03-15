package dto

type TriggerWebhookInput struct {
	Event   string
	SpaceId *string
	Payload map[string]any
}
