package group

import (
	"github.com/labbs/nexo/application/group/dto"
	"github.com/labbs/nexo/domain"
)

// CreateGroup creates a new group
func (app *GroupApplication) CreateGroup(input dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
	logger := app.Logger.With().Str("component", "application.group.create").Logger()

	group := &domain.Group{
		Name:        input.Name,
		Description: input.Description,
		OwnerId:     input.OwnerId,
		Role:        input.Role,
	}

	if err := app.GroupPers.Create(group); err != nil {
		logger.Error().Err(err).Msg("failed to create group")
		return nil, err
	}

	return &dto.CreateGroupOutput{Group: group}, nil
}
