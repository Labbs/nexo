package ports

import (
	"github.com/labbs/nexo/application/user/dto"
)

type UserPort interface {
	Create(input dto.CreateUserInput) (*dto.CreateUserOutput, error)
	GetByEmail(input dto.GetByEmailInput) (*dto.GetByEmailOutput, error)
	GetByUserId(input dto.GetByUserIdInput) (*dto.GetByUserIdOutput, error)
	CreateFavorite(input dto.CreateFavoriteInput) error
	DeleteFavorite(input dto.DeleteFavoriteInput) error
	GetMyFavorites(input dto.GetMyFavoritesInput) (*dto.GetMyFavoritesOutput, error)
	GetFavoriteByIdAndUserId(input dto.GetFavoriteByIdAndUserIdInput) (*dto.GetFavoriteByIdAndUserIdOutput, error)
	UpdateFavoritePosition(input dto.UpdateFavoritePositionInput) (*dto.UpdateFavoritePositionOutput, error)
}
