package database

import (
	"encoding/json"
	"fmt"

	"github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DatabaseApplication) ListRows(input dto.ListRowsInput) (*dto.ListRowsOutput, error) {
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
		rowItem := dto.RowItem{
			Id:            row.Id,
			Properties:    map[string]interface{}(row.Properties),
			Content:       map[string]interface{}(row.Content),
			ShowInSidebar: row.ShowInSidebar,
			CreatedBy:     row.CreatedBy,
			UpdatedBy:     row.UpdatedBy,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}
		if row.CreatedUser.Id != "" {
			rowItem.CreatedByUser = &dto.UserInfo{
				Id:        row.CreatedUser.Id,
				Username:  row.CreatedUser.Username,
				AvatarUrl: row.CreatedUser.AvatarUrl,
			}
		}
		if row.UpdatedUser.Id != "" {
			rowItem.UpdatedByUser = &dto.UserInfo{
				Id:        row.UpdatedUser.Id,
				Username:  row.UpdatedUser.Username,
				AvatarUrl: row.UpdatedUser.AvatarUrl,
			}
		}
		output.Rows[i] = rowItem
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
