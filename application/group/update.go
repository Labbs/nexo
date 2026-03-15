package group

import (
	"github.com/labbs/nexo/application/group/dto"
)

// UpdateGroup updates a group's name, description, or role
func (app *GroupApplication) UpdateGroup(input dto.UpdateGroupInput) error {
	logger := app.Logger.With().Str("component", "application.group.update").Logger()

	group, err := app.GroupPers.GetById(input.GroupId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get group")
		return err
	}

	group.Name = input.Name
	group.Description = input.Description
	group.Role = input.Role

	if err := app.GroupPers.Update(group); err != nil {
		logger.Error().Err(err).Msg("failed to update group")
		return err
	}

	return nil
}
