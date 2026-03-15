package dto

type DeleteDocumentUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	DocumentId   string
	TargetUserId string
}
