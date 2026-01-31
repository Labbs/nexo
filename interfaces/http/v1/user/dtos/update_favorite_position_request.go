package dtos

type UpdateFavoritePositionRequest struct {
    FavoriteId string `path:"favorite_id" validate:"required,uuid4"`
    Position   int    `json:"position" validate:"required,min=0"`
}

type UpdateFavoritePositionResponse struct {
    FavoriteId string `json:"favorite_id"`
    Position   int    `json:"position"`
}




