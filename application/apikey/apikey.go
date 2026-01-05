package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/apikey/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type ApiKeyApp struct {
	Config     config.Config
	Logger     zerolog.Logger
	ApiKeyPers domain.ApiKeyPers
}

func NewApiKeyApp(config config.Config, logger zerolog.Logger, apiKeyPers domain.ApiKeyPers) *ApiKeyApp {
	return &ApiKeyApp{
		Config:     config,
		Logger:     logger,
		ApiKeyPers: apiKeyPers,
	}
}

// generateApiKey generates a secure random API key
func generateApiKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "zk_" + hex.EncodeToString(bytes), nil
}

// hashApiKey creates a SHA-256 hash of the API key
func hashApiKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (app *ApiKeyApp) CreateApiKey(input dto.CreateApiKeyInput) (*dto.CreateApiKeyOutput, error) {
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

func (app *ApiKeyApp) ListApiKeys(input dto.ListApiKeysInput) (*dto.ListApiKeysOutput, error) {
	apiKeys, err := app.ApiKeyPers.GetByUserId(input.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	output := &dto.ListApiKeysOutput{
		ApiKeys: make([]dto.ApiKeyItem, len(apiKeys)),
	}

	for i, k := range apiKeys {
		var scopes []string
		if k.Permissions != nil {
			if s, ok := k.Permissions["scopes"].([]interface{}); ok {
				for _, scope := range s {
					if str, ok := scope.(string); ok {
						scopes = append(scopes, str)
					}
				}
			}
		}

		output.ApiKeys[i] = dto.ApiKeyItem{
			Id:         k.Id,
			Name:       k.Name,
			KeyPrefix:  k.KeyPrefix,
			Scopes:     scopes,
			LastUsedAt: k.LastUsedAt,
			ExpiresAt:  k.ExpiresAt,
			CreatedAt:  k.CreatedAt,
		}
	}

	return output, nil
}

func (app *ApiKeyApp) DeleteApiKey(input dto.DeleteApiKeyInput) error {
	apiKey, err := app.ApiKeyPers.GetById(input.ApiKeyId)
	if err != nil {
		return fmt.Errorf("API key not found: %w", err)
	}

	// Verify ownership
	if apiKey.UserId != input.UserId {
		return fmt.Errorf("access denied")
	}

	if err := app.ApiKeyPers.Delete(input.ApiKeyId); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}

func (app *ApiKeyApp) UpdateApiKey(input dto.UpdateApiKeyInput) error {
	apiKey, err := app.ApiKeyPers.GetById(input.ApiKeyId)
	if err != nil {
		return fmt.Errorf("API key not found: %w", err)
	}

	// Verify ownership
	if apiKey.UserId != input.UserId {
		return fmt.Errorf("access denied")
	}

	if input.Name != nil {
		apiKey.Name = *input.Name
	}

	if input.Scopes != nil {
		apiKey.Permissions = domain.JSONB{
			"scopes": *input.Scopes,
		}
	}

	apiKey.UpdatedAt = time.Now()

	if err := app.ApiKeyPers.Update(apiKey); err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
	}

	return nil
}

func (app *ApiKeyApp) ValidateApiKey(input dto.ValidateApiKeyInput) (*dto.ValidateApiKeyOutput, error) {
	keyHash := hashApiKey(input.Key)

	apiKey, err := app.ApiKeyPers.GetByKeyHash(keyHash)
	if err != nil {
		return &dto.ValidateApiKeyOutput{Valid: false}, nil
	}

	// Check if expired
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return &dto.ValidateApiKeyOutput{Valid: false, Expired: true}, nil
	}

	// Update last used timestamp asynchronously
	go func() {
		_ = app.ApiKeyPers.UpdateLastUsed(apiKey.Id)
	}()

	// Extract scopes
	var scopes []string
	if apiKey.Permissions != nil {
		if s, ok := apiKey.Permissions["scopes"].([]interface{}); ok {
			for _, scope := range s {
				if str, ok := scope.(string); ok {
					scopes = append(scopes, str)
				}
			}
		}
	}

	return &dto.ValidateApiKeyOutput{
		Valid:  true,
		UserId: apiKey.UserId,
		Scopes: scopes,
	}, nil
}
