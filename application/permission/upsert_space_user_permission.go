package permission

import (
	"fmt"

	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// UpsertSpaceUserPermission adds or updates a user permission for a space.
// The requester must be an admin of the space.
func (app *PermissionApplication) UpsertSpaceUserPermission(input spaceDto.UpsertSpaceUserPermissionInput) error {
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: input.SpaceId})
	if err != nil || spaceResult.Space == nil {
		return fmt.Errorf("not_found")
	}
	if !spaceResult.Space.HasPermission(input.RequesterId, "admin") {
		return fmt.Errorf("forbidden")
	}

	// Prevent changing the owner's role on personal/private spaces
	if (spaceResult.Space.Type == "personal" || spaceResult.Space.Type == "private") &&
		spaceResult.Space.OwnerId != nil && *spaceResult.Space.OwnerId == input.TargetUserId && input.Role != "owner" {
		return fmt.Errorf("cannot_change_owner_role")
	}

	return app.PermissionPers.UpsertUser(domain.PermissionTypeSpace, input.SpaceId, input.TargetUserId, domain.PermissionRole(input.Role))
}
