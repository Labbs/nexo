package database

import (
	"fmt"

	"github.com/labbs/nexo/application/database/dto"
)

func (app *DatabaseApplication) Search(input dto.SearchDatabasesInput) (*dto.SearchDatabasesOutput, error) {
	if len(input.Query) < 2 {
		return nil, fmt.Errorf("query must be at least 2 characters")
	}

	databases, err := app.DatabasePers.Search(input.Query, input.UserId, input.SpaceId, input.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search databases: %w", err)
	}

	output := &dto.SearchDatabasesOutput{
		Results: make([]dto.SearchDatabaseResultItem, len(databases)),
	}

	for i, db := range databases {
		output.Results[i] = dto.SearchDatabaseResultItem{
			Id:          db.Id,
			Name:        db.Name,
			Description: db.Description,
			Icon:        db.Icon,
			Type:        string(db.Type),
			SpaceId:     db.SpaceId,
			SpaceName:   db.Space.Name,
			UpdatedAt:   db.UpdatedAt,
		}
	}

	return output, nil
}
