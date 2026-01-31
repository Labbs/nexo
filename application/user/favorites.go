package user

import (
	"fmt"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/nexo/application/user/dto"
	"github.com/labbs/nexo/domain"
)

func (c *UserApp) CreateFavorite(input dto.CreateFavoriteInput) error {
	logger := c.Logger.With().Str("component", "application.user.create_favorite").Logger()

	latestPosition, err := c.FavoritePers.GetLatestFavoritePositionByUser(input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get latest favorite position")
		return err
	}

	favorite := &domain.Favorite{
		Id:         utils.UUIDv4(),
		DocumentId: input.DocumentId,
		SpaceId:    input.SpaceId,
		UserId:     input.UserId,
		Position:   latestPosition + 1,
	}

	err = c.FavoritePers.Create(favorite)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create favorite")
		return err
	}

	return nil
}

func (c *UserApp) DeleteFavorite(input dto.DeleteFavoriteInput) error {
	logger := c.Logger.With().Str("component", "application.user.delete_favorite").Logger()

	err := c.FavoritePers.Delete(input.DocumentId, input.UserId, input.SpaceId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete favorite")
		return err
	}

	return nil
}

func (c *UserApp) GetMyFavorites(input dto.GetMyFavoritesInput) (*dto.GetMyFavoritesOutput, error) {
	logger := c.Logger.With().Str("component", "application.user.get_my_favorites").Logger()

	favorites, err := c.FavoritePers.GetMyFavoritesWithMainDocumentInformations(input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get my favorites")
		return nil, err
	}

	return &dto.GetMyFavoritesOutput{Favorites: favorites}, nil
}

func (c *UserApp) GetFavoriteByIdAndUserId(input dto.GetFavoriteByIdAndUserIdInput) (*dto.GetFavoriteByIdAndUserIdOutput, error) {
	logger := c.Logger.With().Str("component", "application.user.get_favorite_by_id_and_user_id").Logger()

	favorite, err := c.FavoritePers.GetFavoriteByIdAndUserId(input.FavoriteId, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get favorite by id and user id")
		return nil, err
	}

	return &dto.GetFavoriteByIdAndUserIdOutput{Favorite: favorite}, nil
}

func (c *UserApp) UpdateFavoritePosition(input dto.UpdateFavoritePositionInput) (*dto.UpdateFavoritePositionOutput, error) {
	logger := c.Logger.With().Str("component", "application.user.update_favorite_position").Logger()

	favorite, err := c.FavoritePers.GetFavoriteByIdAndUserId(input.FavoriteId, input.UserId)
	if err != nil || favorite == nil {
		logger.Error().Err(err).Msg("failed to get favorite for update position")
		return nil, fmt.Errorf("not_found")
	}

	favorite.Position = input.NewPosition
	if err := c.FavoritePers.UpdateFavoritePosition(favorite); err != nil {
		logger.Error().Err(err).Msg("failed to update favorite position")
		return nil, err
	}

	return &dto.UpdateFavoritePositionOutput{Favorite: favorite}, nil
}
