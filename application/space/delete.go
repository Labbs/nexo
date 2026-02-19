package space

import (
	"fmt"

	docDto "github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
)

func (c *SpaceApplication) DeleteSpace(input dto.DeleteSpaceInput) error {
	logger := c.Logger.With().Str("component", "application.space.delete_space").Logger()

	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return fmt.Errorf("not_found")
	}

	// Only owner can delete (MVP policy)
	if !space.HasPermission(input.UserId, domain.PermissionRoleOwner) {
		return fmt.Errorf("forbidden")
	}

	// Guard: forbid delete if there are active documents in space (MVP: check root docs)
	docResult, derr := c.DocumentApp.HasDocumentsInSpace(docDto.HasDocumentsInSpaceInput{
		SpaceId: input.SpaceId,
		UserId:  input.UserId,
	})
	if derr == nil && docResult.HasDocuments {
		return fmt.Errorf("conflict_children")
	}

	if err := c.SpacePres.Delete(input.SpaceId); err != nil {
		logger.Error().Err(err).Msg("failed to delete space")
		return err
	}

	return nil
}
