package permission

import (
	"fmt"

	drawingDto "github.com/labbs/nexo/application/drawing/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// UpsertDrawingUserPermission adds or updates a user permission for a drawing.
// The requester must be an owner or admin of the parent space.
func (app *PermissionApplication) UpsertDrawingUserPermission(input drawingDto.UpsertDrawingUserPermissionInput) error {
	// Get the drawing
	drawingResult, err := app.DrawingApplication.GetDrawingById(drawingDto.GetDrawingByIdInput{DrawingId: input.DrawingId})
	if err != nil {
		return fmt.Errorf("not_found")
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: drawingResult.Drawing.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	// User must be admin/owner of space to manage permissions
	role := spaceResult.Space.GetUserRole(input.RequesterId)
	if role == nil || (*role != "owner" && *role != "admin") {
		return fmt.Errorf("forbidden")
	}

	// Prevent changing the owner's role in personal/private spaces
	if (spaceResult.Space.Type == "personal" || spaceResult.Space.Type == "private") &&
		spaceResult.Space.OwnerId != nil && *spaceResult.Space.OwnerId == input.TargetUserId && input.Role != "owner" {
		return fmt.Errorf("cannot_change_owner_role")
	}

	return app.PermissionPers.UpsertUser(domain.PermissionTypeDrawing, input.DrawingId, input.TargetUserId, domain.PermissionRole(input.Role))
}
