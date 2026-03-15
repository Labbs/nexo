package apikey

import (
	"fmt"

	"github.com/labbs/nexo/application/apikey/dto"
)

func (app *ApiKeyApplication) ListApiKeys(input dto.ListApiKeysInput) (*dto.ListApiKeysOutput, error) {
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
			if s, ok := k.Permissions["scopes"].([]any); ok {
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
