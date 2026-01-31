package dtos

type GetDocumentsFromSpaceRequest struct {
	SpaceId  string `path:"space_id" validate:"required,uuid4"`
	ParentId string `query:"parent_id" validate:"omitempty,uuid4"`
}

type GetDocumentsFromSpaceResponse struct {
	Documents []Document `json:"documents"`
}
