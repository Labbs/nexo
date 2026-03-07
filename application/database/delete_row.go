package database

import (
	"fmt"

	"github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
)

func (app *DatabaseApplication) DeleteRow(input dto.DeleteRowInput) error {
	row, err := app.DatabaseRowPers.GetById(input.RowId)
	if err != nil {
		return fmt.Errorf("row not found: %w", err)
	}

	if row.DatabaseId != input.DatabaseId {
		return fmt.Errorf("row not found in this database")
	}

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
		return fmt.Errorf("access denied")
	}

	if err := app.DatabaseRowPers.Delete(input.RowId); err != nil {
		return fmt.Errorf("failed to delete row: %w", err)
	}

	return nil
}
