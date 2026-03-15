package group

import (
	"github.com/labbs/nexo/application/group/dto"
)

// GetAllGroups retrieves all groups with pagination
func (app *GroupApplication) GetAllGroups(input dto.GetAllGroupsInput) (*dto.GetAllGroupsOutput, error) {
	logger := app.Logger.With().Str("component", "application.group.get_all").Logger()

	groups, total, err := app.GroupPers.GetAll(input.Limit, input.Offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get all groups")
		return nil, err
	}

	return &dto.GetAllGroupsOutput{Groups: groups, TotalCount: total}, nil
}
