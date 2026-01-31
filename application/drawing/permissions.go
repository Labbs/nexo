package drawing

import (
	"fmt"

	"github.com/labbs/nexo/application/drawing/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DrawingApp) ListDrawingPermissions(input dto.ListDrawingPermissionsInput) (*dto.ListDrawingPermissionsOutput, error) {
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

func (app *DrawingApp) UpsertDrawingUserPermission(input dto.UpsertDrawingUserPermissionInput) error {
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

	// Prevent changing the owner's role in personal/private spaces
	if (space.Type == domain.SpaceTypePersonal || space.Type == domain.SpaceTypePrivate) &&
		space.OwnerId != nil && *space.OwnerId == input.TargetUserId && input.Role != domain.PermissionRoleOwner {
		return fmt.Errorf("cannot_change_owner_role")
	}

	return app.PermissionPers.UpsertUser(domain.PermissionTypeDrawing, input.DrawingId, input.TargetUserId, input.Role)
}

func (app *DrawingApp) DeleteDrawingUserPermission(input dto.DeleteDrawingUserPermissionInput) error {
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
