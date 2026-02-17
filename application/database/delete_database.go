package database

import (
	"fmt"

	"github.com/labbs/nexo/application/database/dto"
)

func (app *DatabaseApplication) DeleteDatabase(input dto.DeleteDatabaseInput) error {
	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return fmt.Errorf("database not found: %w", err)
	}

	// Verify user has access to the space
	space, err := app.SpacePers.GetSpaceById(database.SpaceId)
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	if space.GetUserRole(input.UserId) == nil {
		return fmt.Errorf("access denied")
	}

	if err := app.DatabasePers.Delete(input.DatabaseId); err != nil {
		return fmt.Errorf("failed to delete database: %w", err)
	}

	return nil
}
