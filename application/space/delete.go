package space

import (
	"fmt"

	"github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

func (c *SpaceApp) DeleteSpace(input dto.DeleteSpaceInput) error {
	logger := c.Logger.With().Str("component", "application.space.delete_space").Logger()

	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return fmt.Errorf("not_found")
	}

	// Only owner can delete (MVP policy)
	if !space.HasPermission(input.UserId, domain.SpaceRoleOwner) {
		return fmt.Errorf("forbidden")
	}

	// Guard: forbid delete if there are active documents in space (MVP: check root docs)
	docs, derr := c.DocumentPers.GetRootDocumentsFromSpaceWithUserPermissions(input.SpaceId, input.UserId)
	if derr == nil && len(docs) > 0 {
		return fmt.Errorf("conflict_children")
	}

	if err := c.SpacePres.Delete(input.SpaceId); err != nil {
		logger.Error().Err(err).Msg("failed to delete space")
		return err
	}

	return nil
}
