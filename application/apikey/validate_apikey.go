package apikey

import (
	"time"

	"github.com/labbs/nexo/application/apikey/dto"
)

func (app *ApiKeyApplication) ValidateApiKey(input dto.ValidateApiKeyInput) (*dto.ValidateApiKeyOutput, error) {
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
