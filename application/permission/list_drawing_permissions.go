package permission

import (
	"fmt"

	dto "github.com/labbs/nexo/application/drawing/dto"
	"github.com/labbs/nexo/domain"
)

// ListDrawingPermissions returns all permissions for a drawing.
// The requester must have access to the parent space.
func (app *PermissionApplication) ListDrawingPermissions(input dto.ListDrawingPermissionsInput) (*dto.ListDrawingPermissionsOutput, error) {
	// Get the drawing
	drawing, err := app.DrawingPers.GetById(input.DrawingId)
	if err != nil {
		return nil, fmt.Errorf("not_found")
	}

	// Verify user has access to the space
	space, err := app.SpacePers.GetSpaceById(drawing.SpaceId)
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if space.GetUserRole(input.RequesterId) == nil {
		return nil, fmt.Errorf("forbidden")
	}

	permissions, err := app.PermissionPers.ListByResource(domain.PermissionTypeDrawing, input.DrawingId)
	if err != nil {
		return nil, err
	}

	return &dto.ListDrawingPermissionsOutput{Permissions: permissions}, nil
}
