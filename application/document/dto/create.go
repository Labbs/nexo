package dto

import "github.com/labbs/nexo/domain"

type CreateDocumentInput struct {
	Name       string
	UserId     string
	SpaceId    string
	Content    []Block
	ParentId   *string
	TemplateId *string // if set, content/config/name are cloned from this template
}

type CreateDocumentOutput struct {
	Document *domain.Document
}
