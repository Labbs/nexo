package dto

type GetDocumentWithSpaceInput struct {
	UserId     string
	SpaceId    string
	DocumentId *string
	Slug       *string
}

type GetDocumentWithSpaceOutput struct {
	Document *Document
}
