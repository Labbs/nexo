package dtos

type MoveDocumentRequest struct {
	SpaceId  string  `path:"space_id" validate:"required,uuid4"`
	Id       string  `path:"id" validate:"required"`
	ParentId *string `json:"parent_id,omitempty" validate:"omitempty,uuid4"`
}

type MoveDocumentResponse struct {
	Id       string  `json:"id"`
	ParentId *string `json:"parent_id,omitempty"`
	SpaceId  string  `json:"space_id"`
}




