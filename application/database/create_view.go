package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/database/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DatabaseApplication) CreateView(input dto.CreateViewInput) (*dto.CreateViewOutput, error) {
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

	// Parse existing views
	var views []dto.ViewConfig
	if database.Views != nil {
		viewsJSON, _ := json.Marshal(database.Views)
		json.Unmarshal(viewsJSON, &views)
	}

	// Create new view
	newView := dto.ViewConfig{
		Id:      uuid.New().String(),
		Name:    input.Name,
		Type:    input.Type,
		Filter:  input.Filter,
		Sort:    input.Sort,
		Columns: input.Columns,
	}

	// Add to views array
	views = append(views, newView)

	// Save back to database
	viewsJSON, _ := json.Marshal(views)
	var viewsArray domain.JSONBArray
	json.Unmarshal(viewsJSON, &viewsArray)
	database.Views = viewsArray
	database.UpdatedAt = time.Now()

	if err := app.DatabasePers.Update(database); err != nil {
		return nil, fmt.Errorf("failed to create view: %w", err)
	}

	return &dto.CreateViewOutput{
		Id:      newView.Id,
		Name:    newView.Name,
		Type:    newView.Type,
		Filter:  newView.Filter,
		Sort:    newView.Sort,
		Columns: newView.Columns,
	}, nil
}
