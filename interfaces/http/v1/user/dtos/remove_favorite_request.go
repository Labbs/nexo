package dtos

type RemoveFavoriteRequest struct {
	FavoriteId string `path:"favorite_id" validate:"required,uuid4"`
}

type RemoveFavoriteResponse struct {
	Message string `json:"message"`
}
