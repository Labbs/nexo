package database

import (
	"fmt"

	"github.com/labbs/nexo/application/database/dto"
)

func (app *DatabaseApplication) ListDatabases(input dto.ListDatabasesInput) (*dto.ListDatabasesOutput, error) {
	// Verify user has access to the space
	space, err := app.SpacePers.GetSpaceById(input.SpaceId)
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	databases, err := app.DatabasePers.GetBySpaceId(input.SpaceId)
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}

	output := &dto.ListDatabasesOutput{
		Databases: make([]dto.DatabaseItem, len(databases)),
	}

	for i, db := range databases {
		rowCount, _ := app.DatabaseRowPers.GetRowCount(db.Id)
		output.Databases[i] = dto.DatabaseItem{
			Id:          db.Id,
			DocumentId:  db.DocumentId,
			Name:        db.Name,
			Description: db.Description,
			Icon:        db.Icon,
			Type:        string(db.Type),
			RowCount:    rowCount,
			CreatedBy:   db.User.Username,
			CreatedAt:   db.CreatedAt,
			UpdatedAt:   db.UpdatedAt,
		}
	}

	return output, nil
}
