package user

import (
	"fmt"

	"github.com/labbs/nexo/application/user/dto"
	"golang.org/x/crypto/bcrypt"
)

func (c *UserApp) UpdateProfile(input dto.UpdateProfileInput) (*dto.UpdateProfileOutput, error) {
	logger := c.Logger.With().Str("component", "application.user.update_profile").Logger()

	user, err := c.UserPres.GetById(input.UserId)
	if err != nil {
		logger.Error().Err(err).Str("user_id", input.UserId).Msg("failed to get user")
		return nil, fmt.Errorf("user not found")
	}

	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.AvatarUrl != nil {
		user.AvatarUrl = *input.AvatarUrl
	}
	if input.Preferences != nil {
		user.Preferences = *input.Preferences
	}

	err = c.UserPres.Update(&user)
	if err != nil {
		logger.Error().Err(err).Str("user_id", input.UserId).Msg("failed to update user")
		return nil, fmt.Errorf("failed to update profile")
	}

	return &dto.UpdateProfileOutput{User: &user}, nil
}

func (c *UserApp) ChangePassword(input dto.ChangePasswordInput) error {
	logger := c.Logger.With().Str("component", "application.user.change_password").Logger()

	user, err := c.UserPres.GetById(input.UserId)
	if err != nil {
		logger.Error().Err(err).Str("user_id", input.UserId).Msg("failed to get user")
		return fmt.Errorf("user not found")
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword))
	if err != nil {
		logger.Warn().Str("user_id", input.UserId).Msg("invalid current password")
		return fmt.Errorf("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("failed to hash new password")
		return fmt.Errorf("failed to process password")
	}

	err = c.UserPres.UpdatePassword(input.UserId, string(hashedPassword))
	if err != nil {
		logger.Error().Err(err).Str("user_id", input.UserId).Msg("failed to update password")
		return fmt.Errorf("failed to update password")
	}

	return nil
}
