package dto

type RemoveFavoriteInput struct {
	FavoriteId string
	UserId     string
}

type RemoveFavoriteOutput struct {
	Message string
}
