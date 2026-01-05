package space

import (
	"fmt"

	"github.com/gosimple/slug"
	"github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/shortuuid"
)

func (c *SpaceApp) UpdateSpace(input dto.UpdateSpaceInput) (*dto.UpdateSpaceOutput, error) {
	logger := c.Logger.With().Str("component", "application.space.update_space").Logger()

	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil || space == nil {
		return nil, fmt.Errorf("not_found")
	}

	// Require admin to update
	if !space.HasPermission(input.UserId, domain.SpaceRoleAdmin) {
		return nil, fmt.Errorf("forbidden")
	}

	// Apply updates
	if input.Name != nil && *input.Name != "" && space.Name != *input.Name {
		space.Name = *input.Name
		space.Slug = slug.Make(space.Name + "-" + shortuuid.GenerateShortUUID())
	}
	if input.Icon != nil {
		space.Icon = *input.Icon
	}
	if input.IconColor != nil {
		space.IconColor = *input.IconColor
	}

	if err := c.SpacePres.Update(space); err != nil {
		logger.Error().Err(err).Msg("failed to update space")
		return nil, err
	}

	return &dto.UpdateSpaceOutput{Space: space}, nil
}
