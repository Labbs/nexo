package dto

type AddFavoriteInput struct {
	UserId     string
	DocumentId string
	SpaceId    string
}

type AddFavoriteOutput struct {
	Favorite *Favorite
}
