package dto

import "github.com/labbs/nexo/domain"

type ListDocumentPermissionsInput struct {
	RequesterId string
	SpaceId     string
	DocumentId  string
}

type ListDocumentPermissionsOutput struct {
	Permissions []domain.DocumentPermission
}

type UpsertDocumentUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	DocumentId   string
	TargetUserId string
	Role         domain.DocumentRole
}

type DeleteDocumentUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	DocumentId   string
	TargetUserId string
}
