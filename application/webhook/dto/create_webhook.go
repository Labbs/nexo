package dto

type CreateWebhookInput struct {
	UserId  string
	SpaceId *string
	Name    string
	Url     string
	Events  []string
}

type CreateWebhookOutput struct {
	Id     string
	Name   string
	Url    string
	Secret string
	Events []string
	Active bool
}
