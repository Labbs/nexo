package dto

type ListApiKeysInput struct {
	UserId string
}

type ListApiKeysOutput struct {
	ApiKeys []ApiKeyItem
}
