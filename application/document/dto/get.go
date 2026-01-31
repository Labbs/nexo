package dto

import "github.com/labbs/nexo/domain"

type GetDocumentWithSpaceInput struct {
	UserId     string
	SpaceId    string
	DocumentId *string
	Slug       *string
}

type GetDocumentWithSpaceOutput struct {
	Document *domain.Document
}

type GetDocumentsFromSpaceInput struct {
	SpaceId  string
	UserId   string
	ParentId *string
}

type GetDocumentsFromSpaceOutput struct {
	Documents []domain.Document
}
