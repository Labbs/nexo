package ports

import (
	"github.com/labbs/nexo/application/apikey/dto"
	"github.com/labbs/nexo/domain"
)

type ApiKeyPort interface {
	CreateApiKey(input dto.CreateApiKeyInput) (*dto.CreateApiKeyOutput, error)
	ListApiKeys(input dto.ListApiKeysInput) (*dto.ListApiKeysOutput, error)
	DeleteApiKey(input dto.DeleteApiKeyInput) error
	UpdateApiKey(input dto.UpdateApiKeyInput) error
	ValidateApiKey(input dto.ValidateApiKeyInput) (*dto.ValidateApiKeyOutput, error)
	GetAllApiKeys(limit, offset int) ([]domain.ApiKey, int64, error)
	AdminDeleteApiKey(apiKeyId string) error
}
