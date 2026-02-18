package permission

import (
	"fmt"

	dto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// DeleteSpaceUserPermission removes a user permission from a space.
// The requester must be an admin of the space.
func (app *PermissionApplication) DeleteSpaceUserPermission(input dto.DeleteSpaceUserPermissionInput) error {
	space, err := app.SpacePers.GetSpaceById(input.SpaceId)
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

	return app.PermissionPers.DeleteUser(domain.PermissionTypeSpace, input.SpaceId, input.TargetUserId)
}
