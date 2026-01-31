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
	if !space.HasPermission(input.UserId, domain.PermissionRoleAdmin) {
		return nil, fmt.Errorf("forbidden")
	}
	permissions, err := c.PermissionPers.ListByResource(domain.PermissionTypeSpace, input.SpaceId)
	if err != nil {
		return nil, err
	}

	// Include the space owner if not already in the permissions list
	if space.OwnerId != nil {
		ownerFound := false
		for _, p := range permissions {
			if p.UserId != nil && *p.UserId == *space.OwnerId {
				ownerFound = true
				break
			}
		}
		if !ownerFound {
			ownerPerm := domain.Permission{
				Id:      "owner-" + space.Id,
				Type:    domain.PermissionTypeSpace,
				SpaceId: &space.Id,
				UserId:  space.OwnerId,
				Role:    domain.PermissionRoleOwner,
				User:    space.Owner,
			}
			permissions = append([]domain.Permission{ownerPerm}, permissions...)
		}
	}

	return &dto.ListSpacePermissionsOutput{Permissions: permissions}, nil
}

func (c *SpaceApp) UpsertSpaceUserPermission(input dto.UpsertSpaceUserPermissionInput) error {
	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return fmt.Errorf("not_found")
	}
	if !space.HasPermission(input.RequesterId, domain.PermissionRoleAdmin) {
		return fmt.Errorf("forbidden")
	}

	// Prevent changing the owner's role on personal/private spaces
	if (space.Type == domain.SpaceTypePersonal || space.Type == domain.SpaceTypePrivate) &&
		space.OwnerId != nil && *space.OwnerId == input.TargetUserId && input.Role != domain.PermissionRoleOwner {
		return fmt.Errorf("cannot_change_owner_role")
	}

	return c.PermissionPers.UpsertUser(domain.PermissionTypeSpace, input.SpaceId, input.TargetUserId, input.Role)
}

func (c *SpaceApp) DeleteSpaceUserPermission(input dto.DeleteSpaceUserPermissionInput) error {
	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return fmt.Errorf("not_found")
	}
	if !space.HasPermission(input.RequesterId, domain.PermissionRoleAdmin) {
		return fmt.Errorf("forbidden")
	}

	// Prevent removing the owner from personal/private spaces
	if (space.Type == domain.SpaceTypePersonal || space.Type == domain.SpaceTypePrivate) &&
		space.OwnerId != nil && *space.OwnerId == input.TargetUserId {
		return fmt.Errorf("cannot_remove_owner")
	}

	return c.PermissionPers.DeleteUser(domain.PermissionTypeSpace, input.SpaceId, input.TargetUserId)
}
