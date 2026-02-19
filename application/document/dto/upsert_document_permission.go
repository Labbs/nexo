package dto

type UpsertDocumentUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	DocumentId   string
	TargetUserId string
	Role         string
}
