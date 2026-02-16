package dto

import "github.com/labbs/nexo/domain"

type UpsertDocumentUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	DocumentId   string
	TargetUserId string
	Role         domain.PermissionRole
}
