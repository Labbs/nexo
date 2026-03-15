package dto

type ValidateApiKeyInput struct {
	Key string
}

type ValidateApiKeyOutput struct {
	Valid   bool
	UserId  string
	Scopes  []string
	Expired bool
}
