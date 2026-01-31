package dto

import "github.com/labbs/nexo/domain"

type MoveDocumentInput struct {
	UserId      string
	SpaceId     string
	DocumentId  string
	NewParentId *string
}

type MoveDocumentOutput struct {
	Document *domain.Document
}
