package dto

import "time"

// Create API key
type CreateApiKeyInput struct {
	UserId    string
	Name      string
	Scopes    []string
	ExpiresAt *time.Time
}

type CreateApiKeyOutput struct {
	Id        string
	Name      string
	Key       string // Plain text key - only returned once at creation
	KeyPrefix string
	Scopes    []string
	ExpiresAt *time.Time
	CreatedAt time.Time
}

// List API keys
type ListApiKeysInput struct {
	UserId string
}

type ApiKeyItem struct {
	Id         string
	Name       string
	KeyPrefix  string
	Scopes     []string
	LastUsedAt *time.Time
	ExpiresAt  *time.Time
	CreatedAt  time.Time
}

type ListApiKeysOutput struct {
	ApiKeys []ApiKeyItem
}

// Delete API key
type DeleteApiKeyInput struct {
	UserId   string
	ApiKeyId string
}

// Update API key
type UpdateApiKeyInput struct {
	UserId   string
	ApiKeyId string
	Name     *string
	Scopes   *[]string
}

// Validate API key (for authentication)
type ValidateApiKeyInput struct {
	Key string
}

type ValidateApiKeyOutput struct {
	Valid   bool
	UserId  string
	Scopes  []string
	Expired bool
}
