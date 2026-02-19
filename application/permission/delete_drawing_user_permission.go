package permission

import (
	"fmt"

	drawingDto "github.com/labbs/nexo/application/drawing/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// DeleteDrawingUserPermission removes a user permission from a drawing.
// The requester must be an owner or admin of the parent space.
func (app *PermissionApplication) DeleteDrawingUserPermission(input drawingDto.DeleteDrawingUserPermissionInput) error {
	// Get the drawing
	drawingResult, err := app.DrawingApp.GetDrawingById(drawingDto.GetDrawingByIdInput{DrawingId: input.DrawingId})
	if err != nil {
		return fmt.Errorf("not_found")
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: drawingResult.Drawing.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	// User must be admin/owner of space to manage permissions
	role := spaceResult.Space.GetUserRole(input.RequesterId)
	if role == nil || (*role != "owner" && *role != "admin") {
		return fmt.Errorf("forbidden")
	}

	// Prevent removing the space owner in personal/private spaces
	if (spaceResult.Space.Type == "personal" || spaceResult.Space.Type == "private") &&
		spaceResult.Space.OwnerId != nil && *spaceResult.Space.OwnerId == input.TargetUserId {
		return fmt.Errorf("cannot_remove_owner")
	}

	return app.PermissionPers.DeleteUser(domain.PermissionTypeDrawing, input.DrawingId, input.TargetUserId)
}
