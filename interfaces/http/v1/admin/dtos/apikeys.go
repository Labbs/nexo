package dtos

import "time"

// API Keys list for admin

type ListAllApiKeysRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type ApiKeyItem struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	KeyPrefix   string    `json:"key_prefix"` // First 8 chars of the key
	UserId      string    `json:"user_id"`
	Username    string    `json:"username"`
	Permissions []string  `json:"permissions"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	LastUsedAt  time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListAllApiKeysResponse struct {
	ApiKeys    []ApiKeyItem `json:"api_keys"`
	TotalCount int64        `json:"total_count"`
}

// Revoke API key

type RevokeApiKeyRequest struct {
	ApiKeyId string `path:"apikey_id"`
}

type RevokeApiKeyResponse struct {
	Message string `json:"message"`
}
