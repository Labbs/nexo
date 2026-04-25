package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
)

func (a *DocumentApplication) ListTemplates(input dto.ListTemplatesInput) (*dto.ListTemplatesOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.list_templates").Logger()

	spaceIds := input.SpaceIds
	if len(spaceIds) == 0 {
		// Fall back to all spaces the user can access
		spacesResult, err := a.SpaceApplication.GetSpacesForUser(spaceDto.GetSpacesForUserInput{UserId: input.UserId})
		if err != nil {
			logger.Error().Err(err).Msg("failed to get spaces for user")
			return nil, fmt.Errorf("failed to get spaces: %w", err)
		}
		for _, s := range spacesResult.Spaces {
			spaceIds = append(spaceIds, s.Id)
		}
	}

	templates, err := a.DocumentPers.ListTemplates(spaceIds, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list templates")
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	return &dto.ListTemplatesOutput{Templates: templates}, nil
}
