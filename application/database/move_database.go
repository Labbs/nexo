package database

import (
	"fmt"
	"time"

	"github.com/labbs/nexo/application/database/dto"
)

func (app *DatabaseApplication) MoveDatabase(input dto.MoveDatabaseInput) (*dto.MoveDatabaseOutput, error) {
	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return nil, fmt.Errorf("database not found: %w", err)
	}

	// Verify user has access to the space
	space, err := app.SpacePers.GetSpaceById(database.SpaceId)
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	database.DocumentId = input.DocumentId
	database.UpdatedAt = time.Now()

	if err := app.DatabasePers.Update(database); err != nil {
		return nil, fmt.Errorf("failed to move database: %w", err)
	}

	return &dto.MoveDatabaseOutput{
		Id:         database.Id,
		DocumentId: database.DocumentId,
	}, nil
}
