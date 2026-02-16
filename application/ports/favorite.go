package ports

import (
	"github.com/labbs/nexo/application/favorite/dto"
)

type FavoritePort interface {
	AddFavorite(input dto.AddFavoriteInput) (*dto.AddFavoriteOutput, error)
	RemoveFavorite(input dto.RemoveFavoriteInput) (*dto.RemoveFavoriteOutput, error)
	GetMyFavorites(input dto.GetMyFavoritesInput) (*dto.GetMyFavoritesOutput, error)
	UpdateFavoritePosition(input dto.UpdateFavoritePositionInput) (*dto.UpdateFavoritePositionOutput, error)
}
