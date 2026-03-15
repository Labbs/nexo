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
