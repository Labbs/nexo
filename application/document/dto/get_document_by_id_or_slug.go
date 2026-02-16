package dto

type GetDocumentByIdOrSlugWithUserPermissionsInput struct {
	DocumentId *string
	Slug       *string
	SpaceId    string
	UserId     string
}

type GetDocumentByIdOrSlugWithUserPermissionsOutput struct {
	Document *Document
}
