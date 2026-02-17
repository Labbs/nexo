package apikey

import (
	"github.com/labbs/nexo/domain"
)

// GetAllApiKeys returns all API keys with pagination (admin only)
func (app *ApiKeyApplication) GetAllApiKeys(limit, offset int) ([]domain.ApiKey, int64, error) {
	logger := app.Logger.With().Str("component", "application.apikey.get_all_apikeys").Logger()

	apiKeys, total, err := app.ApiKeyPers.GetAll(limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get all API keys")
		return nil, 0, err
	}

	return apiKeys, total, nil
}
