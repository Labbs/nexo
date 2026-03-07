package permission

import (
	"fmt"

	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// DeleteSpaceGroupPermission removes a group permission from a space.
// The requester must be an admin or owner of the space.
func (app *PermissionApplication) DeleteSpaceGroupPermission(input spaceDto.DeleteSpaceGroupPermissionInput) error {
	space, err := app.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: input.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	role := space.Space.GetUserRole(input.RequesterId)
	if role == nil || (*role != "owner" && *role != "admin") {
		return fmt.Errorf("forbidden")
	}

	return app.PermissionPers.DeleteGroup(domain.PermissionTypeSpace, input.SpaceId, input.GroupId)
}
