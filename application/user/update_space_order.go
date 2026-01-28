package user

import (
	"fmt"

	"github.com/labbs/nexo/application/user/dto"
	"github.com/labbs/nexo/domain"
)

func (c *UserApp) UpdateSpaceOrder(input dto.UpdateSpaceOrderInput) (*dto.UpdateSpaceOrderOutput, error) {
	logger := c.Logger.With().Str("component", "application.user.update_space_order").Logger()

	user, err := c.UserPres.GetById(input.UserId)
	if err != nil {
		logger.Error().Err(err).Str("user_id", input.UserId).Msg("failed to get user")
		return nil, fmt.Errorf("user not found")
	}

	// Merge space_order into existing preferences (preserve other keys)
	if user.Preferences == nil {
		user.Preferences = domain.JSONB{}
	}

	// Convert []string to []any for JSONB storage
	spaceOrderAny := make([]any, len(input.SpaceIds))
	for i, id := range input.SpaceIds {
		spaceOrderAny[i] = id
	}
	user.Preferences["space_order"] = spaceOrderAny

	err = c.UserPres.Update(&user)
	if err != nil {
		logger.Error().Err(err).Str("user_id", input.UserId).Msg("failed to update space order")
		return nil, fmt.Errorf("failed to update space order")
	}

	return &dto.UpdateSpaceOrderOutput{SpaceIds: input.SpaceIds}, nil
}
