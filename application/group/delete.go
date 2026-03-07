package group

import (
	"github.com/labbs/nexo/application/group/dto"
)

// DeleteGroup deletes a group
func (app *GroupApplication) DeleteGroup(input dto.DeleteGroupInput) error {
	logger := app.Logger.With().Str("component", "application.group.delete").Logger()

	if err := app.GroupPers.Delete(input.GroupId); err != nil {
		logger.Error().Err(err).Msg("failed to delete group")
		return err
	}

	return nil
}
