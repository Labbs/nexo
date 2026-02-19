package permission

import (
	"fmt"

	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// UpsertSpaceGroupPermission adds or updates a group permission for a space.
// The requester must be an admin or owner of the space.
func (app *PermissionApplication) UpsertSpaceGroupPermission(input spaceDto.UpsertSpaceGroupPermissionInput) error {
	space, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: input.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	role := space.Space.GetUserRole(input.RequesterId)
	if role == nil || (*role != "owner" && *role != "admin") {
		return fmt.Errorf("forbidden")
	}

	return app.PermissionPers.UpsertGroup(domain.PermissionTypeSpace, input.SpaceId, input.GroupId, domain.PermissionRole(input.Role))
}
