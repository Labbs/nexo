package space

import (
	"fmt"

	"github.com/labbs/nexo/application/space/dto"
)

func (c *SpaceApp) GetSpacesForUser(input dto.GetSpacesForUserInput) (*dto.GetSpacesForUserOutput, error) {
	logger := c.Logger.With().Str("component", "application.space.getSpacesForUser").Logger()

	spaces, err := c.SpacePres.GetSpacesForUser(input.UserId)
	if err != nil {
		logger.Error().Err(err).Str("userId", input.UserId).Msg("failed to get spaces for user")
		return nil, fmt.Errorf("failed to get spaces for user: %w", err)
	}

	return &dto.GetSpacesForUserOutput{Spaces: spaces}, nil
}
