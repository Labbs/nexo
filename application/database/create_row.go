package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/database/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DatabaseApplication) CreateRow(input dto.CreateRowInput) (*dto.CreateRowOutput, error) {
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
