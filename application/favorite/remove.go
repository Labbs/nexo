package favorite

import (
	"fmt"

	"github.com/labbs/nexo/application/favorite/dto"
)

func (a *FavoriteApp) RemoveFavorite(input dto.RemoveFavoriteInput) (*dto.RemoveFavoriteOutput, error) {
	logger := a.Logger.With().Str("component", "application.favorite.remove_favorite").Logger()

	// Verify the favorite belongs to the user
	favorite, err := a.FavoritePers.GetFavoriteByIdAndUserId(input.FavoriteId, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("favorite not found or does not belong to user")
		return nil, fmt.Errorf("favorite not found: %w", err)
	}

	err = a.FavoritePers.Delete(favorite.DocumentId, favorite.UserId, favorite.SpaceId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete favorite")
		return nil, fmt.Errorf("failed to delete favorite: %w", err)
	}

	return &dto.RemoveFavoriteOutput{
		Message: "Favorite removed successfully",
	}, nil
}
