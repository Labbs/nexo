package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DatabaseApplication) UpdateDatabase(input dto.UpdateDatabaseInput) error {
	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return fmt.Errorf("database not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: database.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return apperrors.ErrAccessDenied
	}

	if input.Name != nil {
		database.Name = *input.Name
	}

	if input.Description != nil {
		database.Description = *input.Description
	}

	if input.Icon != nil {
		database.Icon = *input.Icon
	}

	if input.Schema != nil {
		schemaJSON, err := json.Marshal(input.Schema)
		if err != nil {
			return fmt.Errorf("failed to marshal schema: %w", err)
		}
		var schema domain.JSONBArray
		json.Unmarshal(schemaJSON, &schema)
		database.Schema = schema
	}

	if input.DefaultView != nil {
		database.DefaultView = *input.DefaultView
	}

	database.UpdatedAt = time.Now()

	if err := app.DatabasePers.Update(database); err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}
