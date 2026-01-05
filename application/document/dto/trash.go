package dto

import "github.com/labbs/nexo/domain"

type GetTrashInput struct {
	UserId  string
	SpaceId string
}

type GetTrashOutput struct {
	Documents []domain.Document
}

type RestoreDocumentInput struct {
	UserId     string
	SpaceId    string
	DocumentId string
}
