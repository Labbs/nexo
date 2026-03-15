package group

import (
	"github.com/labbs/nexo/application/group/dto"
	userDto "github.com/labbs/nexo/application/user/dto"
)

// AddMember adds a user to a group
func (app *GroupApplication) AddMember(input dto.AddMemberInput) error {
	logger := app.Logger.With().Str("component", "application.group.add_member").Logger()

	// Verify user exists
	_, err := app.UserApplication.GetByUserId(userDto.GetByUserIdInput{UserId: input.UserId})
	if err != nil {
		logger.Error().Err(err).Str("user_id", input.UserId).Msg("user not found")
		return err
	}

	if err := app.GroupPers.AddMember(input.GroupId, input.UserId); err != nil {
		logger.Error().Err(err).Msg("failed to add member to group")
		return err
	}

	return nil
}
