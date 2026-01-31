package favorite

import (
	"fmt"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/nexo/application/favorite/dto"
	"github.com/labbs/nexo/domain"
)

func (a *FavoriteApp) AddFavorite(input dto.AddFavoriteInput) (*dto.AddFavoriteOutput, error) {
	logger := a.Logger.With().Str("component", "application.favorite.add_favorite").Logger()

	// Verify user has access to the document
	doc, err := a.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, nil, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get document")
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if !doc.HasPermission(input.UserId, domain.PermissionRoleViewer) {
		logger.Error().Msg("user does not have permission to view document")
		return nil, fmt.Errorf("user does not have permission to view document")
	}

	// Get the next position
	latestPosition, err := a.FavoritePers.GetLatestFavoritePositionByUser(input.UserId)
	if err != nil {
		// If no favorites exist, start at 0
		latestPosition = -1
	}

	favorite := &domain.Favorite{
		Id:         utils.UUIDv4(),
		UserId:     input.UserId,
		DocumentId: input.DocumentId,
		SpaceId:    input.SpaceId,
		Position:   latestPosition + 1,
	}

	err = a.FavoritePers.Create(favorite)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create favorite")
		return nil, fmt.Errorf("failed to create favorite: %w", err)
	}

	return &dto.AddFavoriteOutput{
		Favorite: &dto.Favorite{
			Id:         favorite.Id,
			UserId:     favorite.UserId,
			DocumentId: favorite.DocumentId,
			SpaceId:    favorite.SpaceId,
			Position:   favorite.Position,
			CreatedAt:  favorite.CreatedAt,
			Document: &dto.FavoriteDocument{
				Id:   doc.Id,
				Name: doc.Name,
				Slug: doc.Slug,
				Icon: doc.Config.Icon,
			},
		},
	}, nil
}
