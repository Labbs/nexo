package dto

import "github.com/labbs/nexo/domain"

type GetPublicDocumentInput struct {
	SpaceId    string
	DocumentId *string
	Slug       *string
}

type GetPublicDocumentOutput struct {
	Document *domain.Document
}
