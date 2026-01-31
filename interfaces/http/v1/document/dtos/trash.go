package dtos

import "time"

type GetTrashRequest struct {
	SpaceId string `path:"space_id" validate:"required,uuid4"`
}

type TrashDocument struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	DeletedAt time.Time `json:"deleted_at"`
}

type GetTrashResponse struct {
	Documents []TrashDocument `json:"documents"`
}

type RestoreDocumentRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	DocumentId string `path:"document_id" validate:"required,uuid4"`
}

type RestoreDocumentResponse struct {
	Message string `json:"message"`
}
