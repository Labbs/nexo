package dto

import "github.com/labbs/nexo/domain"

type SetPublicInput struct {
	UserId     string
	SpaceId    string
	DocumentId string
	Public     bool
}

type GetPublicDocumentInput struct {
	SpaceId    string
	DocumentId *string
	Slug       *string
}

type GetPublicDocumentOutput struct {
	Document *domain.Document
}
