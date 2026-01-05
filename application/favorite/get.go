package favorite

import (
	"fmt"

	"github.com/labbs/nexo/application/favorite/dto"
)

func (a *FavoriteApp) GetMyFavorites(input dto.GetMyFavoritesInput) (*dto.GetMyFavoritesOutput, error) {
	logger := a.Logger.With().Str("component", "application.favorite.get_my_favorites").Logger()

	favorites, err := a.FavoritePers.GetMyFavoritesWithMainDocumentInformations(input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get favorites")
		return nil, fmt.Errorf("failed to get favorites: %w", err)
	}

	result := make([]dto.Favorite, len(favorites))
	for i, fav := range favorites {
		result[i] = dto.Favorite{
			Id:         fav.Id,
			UserId:     fav.UserId,
			DocumentId: fav.DocumentId,
			SpaceId:    fav.SpaceId,
			Position:   fav.Position,
			CreatedAt:  fav.CreatedAt,
		}

		// Add document info if available
		if fav.Document.Id != "" {
			result[i].Document = &dto.FavoriteDocument{
				Id:   fav.Document.Id,
				Name: fav.Document.Name,
				Slug: fav.Document.Slug,
				Icon: fav.Document.Config.Icon,
			}
		}
	}

	return &dto.GetMyFavoritesOutput{
		Favorites: result,
	}, nil
}
