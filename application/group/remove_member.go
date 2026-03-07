package group

import (
	"github.com/labbs/nexo/application/group/dto"
)

// RemoveMember removes a user from a group
func (app *GroupApplication) RemoveMember(input dto.RemoveMemberInput) error {
	logger := app.Logger.With().Str("component", "application.group.remove_member").Logger()

	if err := app.GroupPers.RemoveMember(input.GroupId, input.UserId); err != nil {
		logger.Error().Err(err).Msg("failed to remove member from group")
		return err
	}

	return nil
}
