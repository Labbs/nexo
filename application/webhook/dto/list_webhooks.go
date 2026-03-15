package dto

type ListWebhooksInput struct {
	UserId string
}

type ListWebhooksOutput struct {
	Webhooks []WebhookItem
}
