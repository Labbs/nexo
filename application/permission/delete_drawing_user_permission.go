package permission

import (
	"fmt"

	dto "github.com/labbs/nexo/application/drawing/dto"
	"github.com/labbs/nexo/domain"
)

// DeleteDrawingUserPermission removes a user permission from a drawing.
// The requester must be an owner or admin of the parent space.
func (app *PermissionApplication) DeleteDrawingUserPermission(input dto.DeleteDrawingUserPermissionInput) error {
	// Get the drawing
	drawing, err := app.DrawingPers.GetById(input.DrawingId)
	if err != nil {
		return fmt.Errorf("not_found")
	}

	// Verify user has access to the space
	space, err := app.SpacePers.GetSpaceById(drawing.SpaceId)
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	// User must be admin/owner of space to manage permissions
	role := space.GetUserRole(input.RequesterId)
	if role == nil || (*role != domain.PermissionRoleOwner && *role != domain.PermissionRoleAdmin) {
		return fmt.Errorf("forbidden")
	}

	// Prevent removing the space owner in personal/private spaces
	if (space.Type == domain.SpaceTypePersonal || space.Type == domain.SpaceTypePrivate) &&
		space.OwnerId != nil && *space.OwnerId == input.TargetUserId {
		return fmt.Errorf("cannot_remove_owner")
	}

	return app.PermissionPers.DeleteUser(domain.PermissionTypeDrawing, input.DrawingId, input.TargetUserId)
}
