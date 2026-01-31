package dto

import "github.com/labbs/nexo/domain"

type ListSpacePermissionsInput struct {
	UserId  string
	SpaceId string
}

type ListSpacePermissionsOutput struct {
	Permissions []domain.Permission
}

type UpsertSpaceUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	TargetUserId string
	Role         domain.PermissionRole
}

type DeleteSpaceUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	TargetUserId string
}
