package permission

import (
	"fmt"

	dto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// UpsertSpaceUserPermission adds or updates a user permission for a space.
// The requester must be an admin of the space.
func (app *PermissionApplication) UpsertSpaceUserPermission(input dto.UpsertSpaceUserPermissionInput) error {
	space, err := app.SpacePers.GetSpaceById(input.SpaceId)
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

	return app.PermissionPers.UpsertUser(domain.PermissionTypeSpace, input.SpaceId, input.TargetUserId, input.Role)
}
