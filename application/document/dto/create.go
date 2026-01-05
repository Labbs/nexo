package dto

import "github.com/labbs/nexo/domain"

type CreateDocumentInput struct {
	Name     string
	UserId   string
	SpaceId  string
	Content  []Block
	ParentId *string
}

type CreateDocumentOutput struct {
	Document *domain.Document
}
