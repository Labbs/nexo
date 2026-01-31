package dto

import "github.com/labbs/nexo/domain"

type UpdateDocumentInput struct {
	UserId     string
	SpaceId    string
	DocumentId string
	Name       *string
	Content    *[]Block
	ParentId   *string
	Config     *domain.DocumentConfig
	Metadata   *domain.JSONB
}

type UpdateDocumentOutput struct {
	Document *domain.Document
}
