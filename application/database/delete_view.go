package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DatabaseApplication) DeleteView(input dto.DeleteViewInput) error {
	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return fmt.Errorf("database not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: database.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return fmt.Errorf("access denied")
	}

	// Parse existing views
	var views []dto.ViewConfig
	if database.Views != nil {
		viewsJSON, _ := json.Marshal(database.Views)
		json.Unmarshal(viewsJSON, &views)
	}

	// Cannot delete the last view
	if len(views) <= 1 {
		return fmt.Errorf("cannot delete last view")
	}

	// Find and remove the view
	found := false
	newViews := make([]dto.ViewConfig, 0, len(views)-1)
	for _, view := range views {
		if view.Id == input.ViewId {
			found = true
			continue
		}
		newViews = append(newViews, view)
	}

	if !found {
		return fmt.Errorf("view not found")
	}

	// Save back to database
	viewsJSON, _ := json.Marshal(newViews)
	var viewsArray domain.JSONBArray
	json.Unmarshal(viewsJSON, &viewsArray)
	database.Views = viewsArray
	database.UpdatedAt = time.Now()

	if err := app.DatabasePers.Update(database); err != nil {
		return fmt.Errorf("failed to delete view: %w", err)
	}

	return nil
}
