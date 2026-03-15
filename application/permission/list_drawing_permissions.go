package permission

import (
	"fmt"

	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	drawingDto "github.com/labbs/nexo/application/drawing/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// ListDrawingPermissions returns all permissions for a drawing.
// The requester must have access to the parent space.
func (app *PermissionApplication) ListDrawingPermissions(input drawingDto.ListDrawingPermissionsInput) (*drawingDto.ListDrawingPermissionsOutput, error) {
	// Get the drawing
	drawingResult, err := app.DrawingApplication.GetDrawingById(drawingDto.GetDrawingByIdInput{DrawingId: input.DrawingId})
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: drawingResult.Drawing.SpaceId})
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.RequesterId) == nil {
		return nil, apperrors.ErrForbidden
	}

	permissions, err := app.PermissionPers.ListByResource(domain.PermissionTypeDrawing, input.DrawingId)
	if err != nil {
		return nil, err
	}

	return &drawingDto.ListDrawingPermissionsOutput{Permissions: permissions}, nil
}
