package database

import (
	"fmt"

	"github.com/labbs/nexo/application/database/dto"
)

// GetDatabaseById returns a database by its ID without authorization checks.
// This is used internally by other applications that need to look up a database.
func (app *DatabaseApplication) GetDatabaseById(input dto.GetDatabaseByIdInput) (*dto.GetDatabaseByIdOutput, error) {
	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return nil, fmt.Errorf("database not found: %w", err)
	}

	detail := &dto.DatabaseDetail{
		Id:        database.Id,
		SpaceId:   database.SpaceId,
		CreatedBy: database.CreatedBy,
	}

	return &dto.GetDatabaseByIdOutput{Database: detail}, nil
}
