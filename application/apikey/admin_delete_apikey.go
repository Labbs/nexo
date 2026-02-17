package apikey

import (
	"fmt"
)

// AdminDeleteApiKey deletes an API key without ownership check (admin only)
func (app *ApiKeyApplication) AdminDeleteApiKey(apiKeyId string) error {
	if err := app.ApiKeyPers.Delete(apiKeyId); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}
