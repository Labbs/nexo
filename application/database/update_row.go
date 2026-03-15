package database

import (
	"fmt"
	"time"

	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DatabaseApplication) UpdateRow(input dto.UpdateRowInput) error {
	row, err := app.DatabaseRowPers.GetById(input.RowId)
	if err != nil {
		return fmt.Errorf("row not found: %w", err)
	}

	if row.DatabaseId != input.DatabaseId {
		return apperrors.ErrRowNotFound
	}

	database, err := app.DatabasePers.GetById(input.DatabaseId)
	if err != nil {
		return fmt.Errorf("database not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: database.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return apperrors.ErrAccessDenied
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

	row.UpdatedBy = input.UserId
	row.UpdatedAt = time.Now()

	if err := app.DatabaseRowPers.Update(row); err != nil {
		return fmt.Errorf("failed to update row: %w", err)
	}

	return nil
}
