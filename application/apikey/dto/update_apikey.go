package dto

type UpdateApiKeyInput struct {
	UserId   string
	ApiKeyId string
	Name     *string
	Scopes   *[]string
}
