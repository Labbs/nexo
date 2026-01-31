package dtos

type CreateDocumentRequest struct {
	SpaceId  string  `path:"space_id" validate:"required,uuid4"`
	ParentId *string `json:"parent_id,omitempty" validate:"omitempty"`
}

type CreateDocumentResponse struct {
	Document
}
