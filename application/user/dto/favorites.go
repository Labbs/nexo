package dto

import "github.com/labbs/nexo/domain"

type CreateFavoriteInput struct {
	DocumentId string
	SpaceId    string
	UserId     string
}

type DeleteFavoriteInput struct {
	DocumentId string
	SpaceId    string
	UserId     string
}

type GetMyFavoritesInput struct {
	UserId string
}

type GetMyFavoritesOutput struct {
	Favorites []domain.Favorite
}

type GetFavoriteByIdAndUserIdInput struct {
	FavoriteId string
	UserId     string
}

type GetFavoriteByIdAndUserIdOutput struct {
	Favorite *domain.Favorite
}

type UpdateFavoritePositionInput struct {
	UserId      string
	FavoriteId  string
	NewPosition int
}

type UpdateFavoritePositionOutput struct {
	Favorite *domain.Favorite
}
