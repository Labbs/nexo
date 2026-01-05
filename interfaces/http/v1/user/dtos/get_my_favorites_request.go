package dtos

type GetMyFavoritesResponse struct {
	Favorites []Favorite `json:"favorites"`
}
