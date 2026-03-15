package dto

import "time"

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
