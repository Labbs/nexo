package dtos

import "time"

// Request DTOs

type EmptyRequest struct{}

type CreateApiKeyRequest struct {
	Name      string     `json:"name"`
	Scopes    []string   `json:"scopes"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type UpdateApiKeyRequest struct {
	ApiKeyId string    `path:"api_key_id"`
	Name     *string   `json:"name,omitempty"`
	Scopes   *[]string `json:"scopes,omitempty"`
}

type DeleteApiKeyRequest struct {
	ApiKeyId string `path:"api_key_id"`
}

// Response DTOs

type CreateApiKeyResponse struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Key       string     `json:"key"` // Only returned once
	KeyPrefix string     `json:"key_prefix"`
	Scopes    []string   `json:"scopes"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type ApiKeyItem struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	Scopes     []string   `json:"scopes"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type ListApiKeysResponse struct {
	ApiKeys []ApiKeyItem `json:"api_keys"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// Available scopes for reference
type AvailableScopesResponse struct {
	Scopes []ScopeInfo `json:"scopes"`
}

type ScopeInfo struct {
	Scope       string `json:"scope"`
	Description string `json:"description"`
}
