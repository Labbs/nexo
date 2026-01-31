package dto

import "github.com/labbs/nexo/domain"

type ListDocumentPermissionsInput struct {
	RequesterId string
	SpaceId     string
	DocumentId  string
}

type ListDocumentPermissionsOutput struct {
	Permissions []domain.Permission
}

type UpsertDocumentUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	DocumentId   string
	TargetUserId string
	Role         domain.PermissionRole
}

type DeleteDocumentUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	DocumentId   string
	TargetUserId string
}
