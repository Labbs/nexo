package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/database/dto"
	permissionDto "github.com/labbs/nexo/application/permission/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DatabaseApplication) CreateDatabase(input dto.CreateDatabaseInput) (*dto.CreateDatabaseOutput, error) {
	// Verify user has access to the space
	spaceResult, err := app.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: input.SpaceId})
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	// Determine database type (default to spreadsheet)
	dbType := domain.DatabaseTypeSpreadsheet
	if input.Type == string(domain.DatabaseTypeDocument) {
		dbType = domain.DatabaseTypeDocument
	}

	// Use default schema if none provided
	schemaToUse := input.Schema
	if len(schemaToUse) == 0 {
		if dbType == domain.DatabaseTypeDocument {
			// Document databases only need a title field by default
			schemaToUse = []dto.PropertySchema{
				{Id: "title", Name: "Title", Type: "title"},
			}
		} else {
			// Spreadsheet databases get name and date
			schemaToUse = []dto.PropertySchema{
				{Id: "name", Name: "Name", Type: "title"},
				{Id: "date", Name: "Date", Type: "date"},
			}
		}
	}

	// Build schema JSONBArray
	schemaJSON, err := json.Marshal(schemaToUse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}
	var schema domain.JSONBArray
	json.Unmarshal(schemaJSON, &schema)

	// Create default view (always table)
	defaultViewType := domain.ViewTypeTable
	defaultView := dto.ViewConfig{
		Id:   uuid.New().String(),
		Name: "Default",
		Type: string(defaultViewType),
	}
	viewsJSON, _ := json.Marshal([]dto.ViewConfig{defaultView})
	var views domain.JSONBArray
	json.Unmarshal(viewsJSON, &views)

	database := &domain.Database{
		Id:          uuid.New().String(),
		SpaceId:     input.SpaceId,
		DocumentId:  input.DocumentId,
		Name:        input.Name,
		Description: input.Description,
		Icon:        input.Icon,
		Schema:      schema,
		Views:       views,
		DefaultView: string(defaultViewType),
		Type:        dbType,
		CreatedBy:   input.UserId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := app.DatabasePers.Create(database); err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Auto-create editor permission for the creator
	// This ensures they retain access even if their space role is downgraded
	if err := app.PermissionApplication.AssignOwnerPermission(permissionDto.AssignOwnerPermissionInput{
		ResourceType: "database",
		ResourceId:   database.Id,
		UserId:       input.UserId,
		Role:         "editor",
	}); err != nil {
		// Log but don't fail - the database is already created
		app.Logger.Warn().Err(err).Str("database_id", database.Id).Str("user_id", input.UserId).Msg("failed to create creator permission")
	}

	// Create sample rows for new spreadsheet databases only (not for document databases)
	if dbType == domain.DatabaseTypeSpreadsheet {
		// Find the first "title" type column to use for sample data names
		now := time.Now()
		var titleColumnId string
		for _, prop := range schemaToUse {
			if prop.Type == "title" {
				titleColumnId = prop.Id
				break
			}
		}
		if titleColumnId == "" && len(schemaToUse) > 0 {
			titleColumnId = schemaToUse[0].Id // fallback to first column
		}

		if titleColumnId != "" {
			sampleNames := []string{"Data 1", "Data 2", "Data 3"}
			for _, name := range sampleNames {
				row := &domain.DatabaseRow{
					Id:         uuid.New().String(),
					DatabaseId: database.Id,
					Properties: domain.JSONB{
						titleColumnId: name,
					},
					CreatedBy: input.UserId,
					CreatedAt: now,
					UpdatedAt: now,
				}
				// Ignore errors for sample data - not critical
				app.DatabaseRowPers.Create(row)
			}
		}
	}

	return &dto.CreateDatabaseOutput{
		Id:          database.Id,
		Name:        database.Name,
		Description: database.Description,
		Icon:        database.Icon,
		Schema:      schemaToUse,
		DefaultView: database.DefaultView,
		Type:        string(database.Type),
		CreatedAt:   database.CreatedAt,
	}, nil
}
