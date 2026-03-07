package permission

import (
	"fmt"

	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// DeleteSpaceUserPermission removes a user permission from a space.
// The requester must be an admin of the space.
func (app *PermissionApplication) DeleteSpaceUserPermission(input spaceDto.DeleteSpaceUserPermissionInput) error {
	spaceResult, err := app.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: input.SpaceId})
	if err != nil || spaceResult.Space == nil {
		return fmt.Errorf("not_found")
	}
	if !spaceResult.Space.HasPermission(input.RequesterId, "admin") {
		return fmt.Errorf("forbidden")
	}

	// Prevent removing the owner from personal/private spaces
	if (spaceResult.Space.Type == "personal" || spaceResult.Space.Type == "private") &&
		spaceResult.Space.OwnerId != nil && *spaceResult.Space.OwnerId == input.TargetUserId {
		return fmt.Errorf("cannot_remove_owner")
	}

	return app.PermissionPers.DeleteUser(domain.PermissionTypeSpace, input.SpaceId, input.TargetUserId)
}
