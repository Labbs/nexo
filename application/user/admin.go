package user

import (
	"github.com/labbs/nexo/domain"
)

// GetAllUsers returns all users with pagination (admin only)
func (c *UserApp) GetAllUsers(limit, offset int) ([]domain.User, int64, error) {
	logger := c.Logger.With().Str("component", "application.user.get_all_users").Logger()

	users, total, err := c.UserPres.GetAll(limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get all users")
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateRole updates a user's role (admin only)
func (c *UserApp) UpdateRole(userId string, role domain.Role) error {
	logger := c.Logger.With().Str("component", "application.user.update_role").Logger()

	err := c.UserPres.UpdateRole(userId, role)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userId).Str("role", string(role)).Msg("failed to update user role")
		return err
	}

	return nil
}

// UpdateActive updates a user's active status (admin only)
func (c *UserApp) UpdateActive(userId string, active bool) error {
	logger := c.Logger.With().Str("component", "application.user.update_active").Logger()

	err := c.UserPres.UpdateActive(userId, active)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userId).Bool("active", active).Msg("failed to update user active status")
		return err
	}

	return nil
}

// DeleteUser deletes a user (admin only)
func (c *UserApp) DeleteUser(userId string) error {
	logger := c.Logger.With().Str("component", "application.user.delete_user").Logger()

	err := c.UserPres.Delete(userId)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userId).Msg("failed to delete user")
		return err
	}

	return nil
}
