package dtos

type AddFavoriteRequest struct {
	DocumentId string `path:"document_id" validate:"required,uuid4"`
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
}

type AddFavoriteResponse struct {
	Message string `json:"message"`
}
