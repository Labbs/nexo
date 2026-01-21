package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/database/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type DatabaseApp struct {
	Config          config.Config
	Logger          zerolog.Logger
	DatabasePers    domain.DatabasePers
	DatabaseRowPers domain.DatabaseRowPers
	SpacePers       domain.SpacePers
	PermissionPers  domain.PermissionPers
}

func NewDatabaseApp(config config.Config, logger zerolog.Logger, databasePers domain.DatabasePers, databaseRowPers domain.DatabaseRowPers, spacePers domain.SpacePers, permissionPers domain.PermissionPers) *DatabaseApp {
	return &DatabaseApp{
		Config:          config,
		Logger:          logger,
		DatabasePers:    databasePers,
		DatabaseRowPers: databaseRowPers,
		SpacePers:       spacePers,
		PermissionPers:  permissionPers,
	}
}

func (app *DatabaseApp) CreateDatabase(input dto.CreateDatabaseInput) (*dto.CreateDatabaseOutput, error) {
	// Verify user has access to the space
	space, err := app.SpacePers.GetSpaceById(input.SpaceId)
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if space.GetUserRole(input.UserId) == nil {
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

	// Create default view (list for document databases, table for spreadsheets)
	defaultViewType := domain.ViewTypeTable
	if dbType == domain.DatabaseTypeDocument {
		defaultViewType = domain.ViewTypeList
	}
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
	if err := app.PermissionPers.UpsertUser(domain.PermissionTypeDatabase, database.Id, input.UserId, domain.PermissionRoleEditor); err != nil {
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

func (app *DatabaseApp) ListDatabases(input dto.ListDatabasesInput) (*dto.ListDatabasesOutput, error) {
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

func (app *DatabaseApp) GetDatabase(input dto.GetDatabaseInput) (*dto.GetDatabaseOutput, error) {
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

func (app *DatabaseApp) UpdateDatabase(input dto.UpdateDatabaseInput) error {
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

	if input.Name != nil {
		database.Name = *input.Name
	}

	if input.Description != nil {
		database.Description = *input.Description
	}

	if input.Icon != nil {
		database.Icon = *input.Icon
	}

	if input.Schema != nil {
		schemaJSON, err := json.Marshal(input.Schema)
		if err != nil {
			return fmt.Errorf("failed to marshal schema: %w", err)
		}
		var schema domain.JSONBArray
		json.Unmarshal(schemaJSON, &schema)
		database.Schema = schema
	}

	if input.DefaultView != nil {
		database.DefaultView = *input.DefaultView
	}

	database.UpdatedAt = time.Now()

	if err := app.DatabasePers.Update(database); err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}

func (app *DatabaseApp) DeleteDatabase(input dto.DeleteDatabaseInput) error {
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

	if err := app.DatabasePers.Delete(input.DatabaseId); err != nil {
		return fmt.Errorf("failed to delete database: %w", err)
	}

	return nil
}

// View operations

func (app *DatabaseApp) CreateView(input dto.CreateViewInput) (*dto.CreateViewOutput, error) {
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

func (app *DatabaseApp) UpdateView(input dto.UpdateViewInput) error {
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

func (app *DatabaseApp) DeleteView(input dto.DeleteViewInput) error {
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

// Row operations

func (app *DatabaseApp) CreateRow(input dto.CreateRowInput) (*dto.CreateRowOutput, error) {
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

	row := &domain.DatabaseRow{
		Id:            uuid.New().String(),
		DatabaseId:    input.DatabaseId,
		Properties:    domain.JSONB(input.Properties),
		Content:       domain.JSONB(input.Content),
		ShowInSidebar: input.ShowInSidebar,
		CreatedBy:     input.UserId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := app.DatabaseRowPers.Create(row); err != nil {
		return nil, fmt.Errorf("failed to create row: %w", err)
	}

	return &dto.CreateRowOutput{
		Id:         row.Id,
		Properties: input.Properties,
		CreatedAt:  row.CreatedAt,
	}, nil
}

func (app *DatabaseApp) ListRows(input dto.ListRowsInput) (*dto.ListRowsOutput, error) {
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

	limit := input.Limit
	if limit <= 0 {
		limit = 50
	}

	// Build query options from view if provided
	var queryOptions domain.RowQueryOptions
	queryOptions.Limit = limit
	queryOptions.Offset = input.Offset

	if input.ViewId != "" {
		// Parse views to find the requested view
		var views []dto.ViewConfig
		if database.Views != nil {
			viewsJSON, _ := json.Marshal(database.Views)
			json.Unmarshal(viewsJSON, &views)
		}

		for _, view := range views {
			if view.Id == input.ViewId {
				// Convert view filter to domain filter
				if view.Filter != nil {
					queryOptions.Filter = convertFilterConfigToDomain(view.Filter)
				}
				// Convert view sort to domain sort
				if len(view.Sort) > 0 {
					queryOptions.Sort = make([]domain.SortRule, len(view.Sort))
					for i, s := range view.Sort {
						queryOptions.Sort[i] = domain.SortRule{
							PropertyId: s.PropertyId,
							Direction:  s.Direction,
						}
					}
				}
				break
			}
		}
	}

	var rows []domain.DatabaseRow
	var totalCount int64

	if queryOptions.Filter != nil || len(queryOptions.Sort) > 0 {
		// Use filtered query
		rows, err = app.DatabaseRowPers.GetByDatabaseIdWithOptions(input.DatabaseId, queryOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list rows: %w", err)
		}
		totalCount, err = app.DatabaseRowPers.GetRowCountWithFilter(input.DatabaseId, queryOptions.Filter)
		if err != nil {
			return nil, fmt.Errorf("failed to get row count: %w", err)
		}
	} else {
		// Use simple query
		rows, err = app.DatabaseRowPers.GetByDatabaseId(input.DatabaseId, limit, input.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to list rows: %w", err)
		}
		totalCount, err = app.DatabaseRowPers.GetRowCount(input.DatabaseId)
		if err != nil {
			return nil, fmt.Errorf("failed to get row count: %w", err)
		}
	}

	output := &dto.ListRowsOutput{
		Rows:       make([]dto.RowItem, len(rows)),
		TotalCount: totalCount,
	}

	for i, row := range rows {
		createdBy := ""
		if row.User.Id != "" {
			createdBy = row.User.Username
		}
		output.Rows[i] = dto.RowItem{
			Id:            row.Id,
			Properties:    map[string]interface{}(row.Properties),
			ShowInSidebar: row.ShowInSidebar,
			CreatedBy:     createdBy,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}
	}

	return output, nil
}

// convertFilterConfigToDomain converts the DTO filter config to domain filter config
func convertFilterConfigToDomain(filter map[string]interface{}) *domain.FilterConfig {
	if filter == nil {
		return nil
	}

	result := &domain.FilterConfig{}

	// Handle "and" filters
	if andRules, ok := filter["and"].([]interface{}); ok {
		for _, r := range andRules {
			if rule, ok := r.(map[string]interface{}); ok {
				result.And = append(result.And, domain.FilterRule{
					Property:  getString(rule, "property"),
					Condition: getString(rule, "condition"),
					Value:     rule["value"],
				})
			}
		}
	}

	// Handle "or" filters
	if orRules, ok := filter["or"].([]interface{}); ok {
		for _, r := range orRules {
			if rule, ok := r.(map[string]interface{}); ok {
				result.Or = append(result.Or, domain.FilterRule{
					Property:  getString(rule, "property"),
					Condition: getString(rule, "condition"),
					Value:     rule["value"],
				})
			}
		}
	}

	return result
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func (app *DatabaseApp) GetRow(input dto.GetRowInput) (*dto.GetRowOutput, error) {
	row, err := app.DatabaseRowPers.GetById(input.RowId)
	if err != nil {
		return nil, fmt.Errorf("row not found: %w", err)
	}

	if row.DatabaseId != input.DatabaseId {
		return nil, fmt.Errorf("row not found in this database")
	}

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

	return &dto.GetRowOutput{
		Id:            row.Id,
		DatabaseId:    row.DatabaseId,
		Properties:    map[string]interface{}(row.Properties),
		Content:       map[string]interface{}(row.Content),
		ShowInSidebar: row.ShowInSidebar,
		CreatedBy:     row.User.Username,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (app *DatabaseApp) UpdateRow(input dto.UpdateRowInput) error {
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
	space, err := app.SpacePers.GetSpaceById(database.SpaceId)
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	if space.GetUserRole(input.UserId) == nil {
		return fmt.Errorf("access denied")
	}

	if input.Properties != nil {
		row.Properties = domain.JSONB(input.Properties)
	}

	if input.Content != nil {
		row.Content = domain.JSONB(input.Content)
	}

	if input.ShowInSidebar != nil {
		row.ShowInSidebar = *input.ShowInSidebar
	}

	row.UpdatedAt = time.Now()

	if err := app.DatabaseRowPers.Update(row); err != nil {
		return fmt.Errorf("failed to update row: %w", err)
	}

	return nil
}

func (app *DatabaseApp) DeleteRow(input dto.DeleteRowInput) error {
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
	space, err := app.SpacePers.GetSpaceById(database.SpaceId)
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	if space.GetUserRole(input.UserId) == nil {
		return fmt.Errorf("access denied")
	}

	if err := app.DatabaseRowPers.Delete(input.RowId); err != nil {
		return fmt.Errorf("failed to delete row: %w", err)
	}

	return nil
}

func (app *DatabaseApp) BulkDeleteRows(input dto.BulkDeleteRowsInput) error {
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

	if err := app.DatabaseRowPers.BulkDelete(input.RowIds); err != nil {
		return fmt.Errorf("failed to delete rows: %w", err)
	}

	return nil
}
