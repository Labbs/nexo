package dtos

type DeleteDocumentRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	Identifier string `path:"identifier" validate:"required"`
}

type DeleteDocumentResponse struct {
	Message string `json:"message"`
}
