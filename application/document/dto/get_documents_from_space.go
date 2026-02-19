package dto

import "github.com/labbs/nexo/domain"

type GetDocumentsFromSpaceInput struct {
	SpaceId  string
	UserId   string
	ParentId *string
}

type GetDocumentsFromSpaceOutput struct {
	Documents []domain.Document
}
