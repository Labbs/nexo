package dtos

type SetPublicRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	DocumentId string `path:"document_id" validate:"required,uuid4"`
	Public     bool   `json:"public"`
}

type SetPublicResponse struct {
	Message string `json:"message"`
	Public  bool   `json:"public"`
}

type GetPublicDocumentRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	Identifier string `path:"identifier" validate:"required"`
}
