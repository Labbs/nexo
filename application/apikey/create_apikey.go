package apikey

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/apikey/dto"
	"github.com/labbs/nexo/domain"
)

func (app *ApiKeyApplication) CreateApiKey(input dto.CreateApiKeyInput) (*dto.CreateApiKeyOutput, error) {
	// Generate the API key
	plainKey, err := generateApiKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Hash the key for storage
	keyHash := hashApiKey(plainKey)

	// Create the key prefix for identification (first 11 chars: "zk_" + 8 chars)
	keyPrefix := plainKey[:11]

	// Build permissions JSON
	permissions := domain.JSONB{
		"scopes": input.Scopes,
	}

	apiKey := &domain.ApiKey{
		Id:          uuid.New().String(),
		UserId:      input.UserId,
		Name:        input.Name,
		KeyHash:     keyHash,
		KeyPrefix:   keyPrefix,
		Permissions: permissions,
		ExpiresAt:   input.ExpiresAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := app.ApiKeyPers.Create(apiKey); err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return &dto.CreateApiKeyOutput{
		Id:        apiKey.Id,
		Name:      apiKey.Name,
		Key:       plainKey, // Only returned once
		KeyPrefix: keyPrefix,
		Scopes:    input.Scopes,
		ExpiresAt: input.ExpiresAt,
		CreatedAt: apiKey.CreatedAt,
	}, nil
}
