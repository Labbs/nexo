package domain

import (
	"time"

	"gorm.io/gorm"
)

type ApiKey struct {
	Id string

	UserId string
	User   User `gorm:"foreignKey:UserId;references:Id"`

	Name        string
	KeyHash     string // Hashed API key (never store plain text)
	KeyPrefix   string // First 8 chars for identification (e.g., "zk_abc123")
	Permissions JSONB  // Scopes: ["read:documents", "write:documents", etc.]

	LastUsedAt *time.Time
	ExpiresAt  *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (k *ApiKey) TableName() string {
	return "api_key"
}

// ApiKeyScope defines available permission scopes
type ApiKeyScope string

const (
	ApiKeyScopeReadDocuments   ApiKeyScope = "read:documents"
	ApiKeyScopeWriteDocuments  ApiKeyScope = "write:documents"
	ApiKeyScopeReadSpaces      ApiKeyScope = "read:spaces"
	ApiKeyScopeWriteSpaces     ApiKeyScope = "write:spaces"
	ApiKeyScopeReadComments    ApiKeyScope = "read:comments"
	ApiKeyScopeWriteComments   ApiKeyScope = "write:comments"
	ApiKeyScopeManageWebhooks  ApiKeyScope = "manage:webhooks"
	ApiKeyScopeManageDatabases ApiKeyScope = "manage:databases"
)

func (k *ApiKey) HasScope(scope ApiKeyScope) bool {
	if k.Permissions == nil {
		return false
	}
	scopes, ok := k.Permissions["scopes"].([]interface{})
	if !ok {
		return false
	}
	for _, s := range scopes {
		if str, ok := s.(string); ok && str == string(scope) {
			return true
		}
	}
	return false
}

type ApiKeyPers interface {
	Create(apiKey *ApiKey) error
	GetById(id string) (*ApiKey, error)
	GetByKeyHash(keyHash string) (*ApiKey, error)
	GetByUserId(userId string) ([]ApiKey, error)
	Update(apiKey *ApiKey) error
	Delete(id string) error
	UpdateLastUsed(id string) error
	// Admin methods
	GetAll(limit, offset int) ([]ApiKey, int64, error)
}
