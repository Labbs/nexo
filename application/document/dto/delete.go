package dto

type DeleteDocumentInput struct {
	UserId     string
	SpaceId    string
	DocumentId *string
	Slug       *string
}
