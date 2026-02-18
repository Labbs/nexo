package permission

import (
	"fmt"

	dto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// ListSpacePermissions returns all permissions for a space.
// The requester must be an admin of the space.
func (app *PermissionApplication) ListSpacePermissions(input dto.ListSpacePermissionsInput) (*dto.ListSpacePermissionsOutput, error) {
	space, err := app.SpacePers.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return nil, fmt.Errorf("not_found")
	}
	if !space.HasPermission(input.UserId, domain.PermissionRoleAdmin) {
		return nil, fmt.Errorf("forbidden")
	}
	permissions, err := app.PermissionPers.ListByResource(domain.PermissionTypeSpace, input.SpaceId)
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
