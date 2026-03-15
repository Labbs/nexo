package apikey

import (
	"fmt"
	"time"

	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/application/apikey/dto"
	"github.com/labbs/nexo/domain"
)

func (app *ApiKeyApplication) UpdateApiKey(input dto.UpdateApiKeyInput) error {
	apiKey, err := app.ApiKeyPers.GetById(input.ApiKeyId)
	if err != nil {
		return fmt.Errorf("API key not found: %w", err)
	}

	// Verify ownership
	if apiKey.UserId != input.UserId {
		return apperrors.ErrAccessDenied
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
