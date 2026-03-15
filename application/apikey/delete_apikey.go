package apikey

import (
	"fmt"

	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/application/apikey/dto"
)

func (app *ApiKeyApplication) DeleteApiKey(input dto.DeleteApiKeyInput) error {
	apiKey, err := app.ApiKeyPers.GetById(input.ApiKeyId)
	if err != nil {
		return fmt.Errorf("API key not found: %w", err)
	}

	// Verify ownership
	if apiKey.UserId != input.UserId {
		return apperrors.ErrAccessDenied
	}

	if err := app.ApiKeyPers.Delete(input.ApiKeyId); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}
