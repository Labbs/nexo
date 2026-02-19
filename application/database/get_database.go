package database

import (
	"encoding/json"
	"fmt"

	"github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
)

func (app *DatabaseApplication) GetDatabase(input dto.GetDatabaseInput) (*dto.GetDatabaseOutput, error) {
	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return nil, fmt.Errorf("database not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: database.SpaceId})
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	// Parse schema
	var schema []dto.PropertySchema
	if database.Schema != nil {
		schemaJSON, _ := json.Marshal(database.Schema)
		json.Unmarshal(schemaJSON, &schema)
	}

	// Parse views
	var views []dto.ViewConfig
	if database.Views != nil {
		viewsJSON, _ := json.Marshal(database.Views)
		json.Unmarshal(viewsJSON, &views)
	}

	return &dto.GetDatabaseOutput{
		Id:          database.Id,
		SpaceId:     database.SpaceId,
		DocumentId:  database.DocumentId,
		Name:        database.Name,
		Description: database.Description,
		Icon:        database.Icon,
		Schema:      schema,
		Views:       views,
		DefaultView: database.DefaultView,
		Type:        string(database.Type),
		CreatedBy:   database.User.Username,
		CreatedAt:   database.CreatedAt,
		UpdatedAt:   database.UpdatedAt,
	}, nil
}
