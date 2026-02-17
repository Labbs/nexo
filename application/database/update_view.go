package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labbs/nexo/application/database/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DatabaseApplication) UpdateView(input dto.UpdateViewInput) error {
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

	// Parse existing views
	var views []dto.ViewConfig
	if database.Views != nil {
		viewsJSON, _ := json.Marshal(database.Views)
		json.Unmarshal(viewsJSON, &views)
	}

	// Find and update the view
	found := false
	for i, view := range views {
		if view.Id == input.ViewId {
			if input.Name != nil {
				views[i].Name = *input.Name
			}
			if input.Type != nil {
				views[i].Type = *input.Type
			}
			// Filter: nil means no change, empty map {} means clear, otherwise update
			if input.Filter != nil {
				if len(input.Filter) == 0 {
					// Empty object {} means clear the filter
					views[i].Filter = nil
				} else {
					views[i].Filter = input.Filter
				}
			}
			// Sort: nil means no change, empty slice [] means clear, otherwise update
			if input.Sort != nil {
				if len(input.Sort) == 0 {
					// Empty array [] means clear the sort
					views[i].Sort = nil
				} else {
					views[i].Sort = input.Sort
				}
			}
			if input.Columns != nil {
				views[i].Columns = input.Columns
			}
			// HiddenColumns: always update when provided (even empty array means "show all")
			if input.HiddenColumns != nil {
				views[i].HiddenColumns = input.HiddenColumns
			}
			// GroupBy: update when provided (used for board view grouping)
			if input.GroupBy != nil {
				views[i].GroupBy = *input.GroupBy
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("view not found")
	}

	// Save back to database
	viewsJSON, _ := json.Marshal(views)
	var viewsArray domain.JSONBArray
	json.Unmarshal(viewsJSON, &viewsArray)
	database.Views = viewsArray
	database.UpdatedAt = time.Now()

	if err := app.DatabasePers.Update(database); err != nil {
		return fmt.Errorf("failed to update view: %w", err)
	}

	return nil
}
