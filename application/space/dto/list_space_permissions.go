package dto

import "github.com/labbs/nexo/domain"

type ListSpacePermissionsInput struct {
	UserId  string
	SpaceId string
}

type ListSpacePermissionsOutput struct {
	Permissions []domain.Permission
}
