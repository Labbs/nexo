package ports

import (
	"github.com/labbs/nexo/application/space/dto"
)

type SpacePort interface {
	CreatePrivateSpaceForUser(input dto.CreatePrivateSpaceForUserInput) (*dto.CreatePrivateSpaceForUserOutput, error)
	CreateSpace(input dto.CreateSpaceInput) (*dto.CreateSpaceOutput, error)
	GetSpacesForUser(input dto.GetSpacesForUserInput) (*dto.GetSpacesForUserOutput, error)
	UpdateSpace(input dto.UpdateSpaceInput) (*dto.UpdateSpaceOutput, error)
	DeleteSpace(input dto.DeleteSpaceInput) error
	ListSpacePermissions(input dto.ListSpacePermissionsInput) (*dto.ListSpacePermissionsOutput, error)
	UpsertSpaceUserPermission(input dto.UpsertSpaceUserPermissionInput) error
	DeleteSpaceUserPermission(input dto.DeleteSpaceUserPermissionInput) error
}
