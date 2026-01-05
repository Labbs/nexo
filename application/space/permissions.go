package space

import (
	"fmt"

	"github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// Permissions (MVP: user-level only)
func (c *SpaceApp) ListSpacePermissions(input dto.ListSpacePermissionsInput) (*dto.ListSpacePermissionsOutput, error) {
	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return nil, fmt.Errorf("not_found")
	}
	if !space.HasPermission(input.UserId, domain.SpaceRoleAdmin) {
		return nil, fmt.Errorf("forbidden")
	}
	permissions, err := c.SpacePres.ListPermissions(input.SpaceId)
	if err != nil {
		return nil, err
	}
	return &dto.ListSpacePermissionsOutput{Permissions: permissions}, nil
}

func (c *SpaceApp) UpsertSpaceUserPermission(input dto.UpsertSpaceUserPermissionInput) error {
	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return fmt.Errorf("not_found")
	}
	if !space.HasPermission(input.RequesterId, domain.SpaceRoleAdmin) {
		return fmt.Errorf("forbidden")
	}

	// Prevent changing the owner's role on personal/private spaces
	if (space.Type == domain.SpaceTypePersonal || space.Type == domain.SpaceTypePrivate) &&
		space.OwnerId != nil && *space.OwnerId == input.TargetUserId && input.Role != domain.SpaceRoleOwner {
		return fmt.Errorf("cannot_change_owner_role")
	}

	return c.SpacePres.UpsertUserPermission(input.SpaceId, input.TargetUserId, input.Role)
}

func (c *SpaceApp) DeleteSpaceUserPermission(input dto.DeleteSpaceUserPermissionInput) error {
	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return fmt.Errorf("not_found")
	}
	if !space.HasPermission(input.RequesterId, domain.SpaceRoleAdmin) {
		return fmt.Errorf("forbidden")
	}

	// Prevent removing the owner from personal/private spaces
	if (space.Type == domain.SpaceTypePersonal || space.Type == domain.SpaceTypePrivate) &&
		space.OwnerId != nil && *space.OwnerId == input.TargetUserId {
		return fmt.Errorf("cannot_remove_owner")
	}

	return c.SpacePres.DeleteUserPermission(input.SpaceId, input.TargetUserId)
}
