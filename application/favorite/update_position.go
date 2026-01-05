package favorite

import (
	"fmt"

	"github.com/labbs/nexo/application/favorite/dto"
)

func (a *FavoriteApp) UpdateFavoritePosition(input dto.UpdateFavoritePositionInput) (*dto.UpdateFavoritePositionOutput, error) {
	logger := a.Logger.With().Str("component", "application.favorite.update_position").Logger()

	// Verify the favorite belongs to the user
	favorite, err := a.FavoritePers.GetFavoriteByIdAndUserId(input.FavoriteId, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("favorite not found or does not belong to user")
		return nil, fmt.Errorf("favorite not found: %w", err)
	}

	// Update the position
	favorite.Position = input.Position

	err = a.FavoritePers.UpdateFavoritePosition(favorite)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update favorite position")
		return nil, fmt.Errorf("failed to update favorite position: %w", err)
	}

	return &dto.UpdateFavoritePositionOutput{
		Favorite: &dto.Favorite{
			Id:         favorite.Id,
			UserId:     favorite.UserId,
			DocumentId: favorite.DocumentId,
			SpaceId:    favorite.SpaceId,
			Position:   favorite.Position,
			CreatedAt:  favorite.CreatedAt,
		},
	}, nil
}
