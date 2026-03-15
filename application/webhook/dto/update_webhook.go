package dto

type UpdateWebhookInput struct {
	UserId    string
	WebhookId string
	Name      *string
	Url       *string
	Events    *[]string
	Active    *bool
}
