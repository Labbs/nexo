package user

import (
	"fmt"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/nexo/application/user/dto"
	helperError "github.com/labbs/nexo/infrastructure/helpers/error"
	"gorm.io/gorm"
)

func (c *UserApp) Create(input dto.CreateUserInput) (*dto.CreateUserOutput, error) {
	logger := c.Logger.With().Str("component", "application.user.create").Logger()

	// Generate UUID for user
	input.User.Id = utils.UUIDv4()

	createdUser, err := c.UserPres.Create(input.User)
	if helperError.Catch(err) == gorm.ErrDuplicatedKey {
		logger.Warn().Str("email", input.User.Email).Msg("user with the same email already exists")
		return nil, fmt.Errorf("user with the same email already exists")
	} else if err != nil {
		logger.Error().Err(err).Str("email", input.User.Email).Msg("failed to create user")
		return nil, err
	}

	return &dto.CreateUserOutput{User: &createdUser}, nil
}
