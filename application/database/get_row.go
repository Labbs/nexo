package database

import (
	"fmt"

	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
)

func (app *DatabaseApplication) GetRow(input dto.GetRowInput) (*dto.GetRowOutput, error) {
	row, err := app.DatabaseRowPers.GetById(input.RowId)
	if err != nil {
		return nil, fmt.Errorf("row not found: %w", err)
	}

	if row.DatabaseId != input.DatabaseId {
		return nil, apperrors.ErrRowNotFound
	}

	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return nil, fmt.Errorf("database not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: database.SpaceId})
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return nil, apperrors.ErrAccessDenied
	}

	output := &dto.GetRowOutput{
		Id:            row.Id,
		DatabaseId:    row.DatabaseId,
		Properties:    map[string]any(row.Properties),
		Content:       map[string]any(row.Content),
		ShowInSidebar: row.ShowInSidebar,
		CreatedBy:     row.CreatedBy,
		UpdatedBy:     row.UpdatedBy,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
	if row.CreatedUser.Id != "" {
		output.CreatedByUser = &dto.UserInfo{
			Id:        row.CreatedUser.Id,
			Username:  row.CreatedUser.Username,
			AvatarUrl: row.CreatedUser.AvatarUrl,
		}
	}
	if row.UpdatedUser.Id != "" {
		output.UpdatedByUser = &dto.UserInfo{
			Id:        row.UpdatedUser.Id,
			Username:  row.UpdatedUser.Username,
			AvatarUrl: row.UpdatedUser.AvatarUrl,
		}
	}
	return output, nil
}
