package permission

import (
	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

// ListSpacePermissions returns all permissions for a space.
// The requester must be an admin of the space.
func (app *PermissionApplication) ListSpacePermissions(input spaceDto.ListSpacePermissionsInput) (*spaceDto.ListSpacePermissionsOutput, error) {
	spaceResult, err := app.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: input.SpaceId})
	if err != nil || spaceResult.Space == nil {
		return nil, apperrors.ErrNotFound
	}
	if !spaceResult.Space.HasPermission(input.UserId, "admin") {
		return nil, apperrors.ErrForbidden
	}
	permissions, err := app.PermissionPers.ListByResource(domain.PermissionTypeSpace, input.SpaceId)
	if err != nil {
		return nil, err
	}

	// Include the space owner if not already in the permissions list
	if spaceResult.Space.OwnerId != nil {
		ownerFound := false
		for _, p := range permissions {
			if p.UserId != nil && *p.UserId == *spaceResult.Space.OwnerId {
				ownerFound = true
				break
			}
		}
		if !ownerFound {
			spaceId := spaceResult.Space.Id
			ownerPerm := domain.Permission{
				Id:      "owner-" + spaceResult.Space.Id,
				Type:    domain.PermissionTypeSpace,
				SpaceId: &spaceId,
				UserId:  spaceResult.Space.OwnerId,
				Role:    domain.PermissionRoleOwner,
			}
			permissions = append([]domain.Permission{ownerPerm}, permissions...)
		}
	}

	return &spaceDto.ListSpacePermissionsOutput{Permissions: permissions}, nil
}
