package dto

type UpdateFavoritePositionInput struct {
	FavoriteId string
	UserId     string
	Position   int
}

type UpdateFavoritePositionOutput struct {
	Favorite *Favorite
}
