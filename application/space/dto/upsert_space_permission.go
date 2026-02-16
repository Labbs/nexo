package dto

import "github.com/labbs/nexo/domain"

type UpsertSpaceUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	TargetUserId string
	Role         domain.PermissionRole
}
