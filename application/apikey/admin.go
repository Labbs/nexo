package apikey

import (
	"fmt"

	"github.com/labbs/nexo/domain"
)

// GetAllApiKeys returns all API keys with pagination (admin only)
func (app *ApiKeyApp) GetAllApiKeys(limit, offset int) ([]domain.ApiKey, int64, error) {
	logger := app.Logger.With().Str("component", "application.apikey.get_all_apikeys").Logger()

	apiKeys, total, err := app.ApiKeyPers.GetAll(limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get all API keys")
		return nil, 0, err
	}

	return apiKeys, total, nil
}

// AdminDeleteApiKey deletes an API key without ownership check (admin only)
func (app *ApiKeyApp) AdminDeleteApiKey(apiKeyId string) error {
	if err := app.ApiKeyPers.Delete(apiKeyId); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}
