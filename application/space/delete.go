package space

import (
	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	docDto "github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

func (c *SpaceApplication) DeleteSpace(input dto.DeleteSpaceInput) error {
	logger := c.Logger.With().Str("component", "application.space.delete_space").Logger()

	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return apperrors.ErrSpaceNotFound
	}

	// Only owner can delete (MVP policy)
	if !space.HasPermission(input.UserId, domain.PermissionRoleOwner) {
		return apperrors.ErrForbidden
	}

	// Guard: forbid delete if there are active documents in space (MVP: check root docs)
	docResult, derr := c.DocumentApplication.HasDocumentsInSpace(docDto.HasDocumentsInSpaceInput{
		SpaceId: input.SpaceId,
		UserId:  input.UserId,
	})
	if derr == nil && docResult.HasDocuments {
		return apperrors.ErrConflictChildren
	}

	if err := c.SpacePres.Delete(input.SpaceId); err != nil {
		logger.Error().Err(err).Msg("failed to delete space")
		return err
	}

	return nil
}
