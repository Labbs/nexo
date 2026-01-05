package dto

type GetMyFavoritesInput struct {
	UserId string
}

type GetMyFavoritesOutput struct {
	Favorites []Favorite
}
