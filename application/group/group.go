package group

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type GroupApp struct {
	Config    config.Config
	Logger    zerolog.Logger
	GroupPers domain.GroupPers
	UserPers  domain.UserPers
}

func NewGroupApp(config config.Config, logger zerolog.Logger, groupPers domain.GroupPers, userPers domain.UserPers) *GroupApp {
	return &GroupApp{
		Config:    config,
		Logger:    logger,
		GroupPers: groupPers,
		UserPers:  userPers,
	}
}

// CreateGroup creates a new group
func (app *GroupApp) CreateGroup(name, description, ownerId string, role domain.Role) (*domain.Group, error) {
	logger := app.Logger.With().Str("component", "application.group.create").Logger()

	group := &domain.Group{
		Name:        name,
		Description: description,
		OwnerId:     ownerId,
		Role:        role,
	}

	if err := app.GroupPers.Create(group); err != nil {
		logger.Error().Err(err).Msg("failed to create group")
		return nil, err
	}

	return group, nil
}

// GetGroup retrieves a group by ID
func (app *GroupApp) GetGroup(groupId string) (*domain.Group, error) {
	return app.GroupPers.GetById(groupId)
}

// GetAllGroups retrieves all groups with pagination
func (app *GroupApp) GetAllGroups(limit, offset int) ([]domain.Group, int64, error) {
	logger := app.Logger.With().Str("component", "application.group.get_all").Logger()

	groups, total, err := app.GroupPers.GetAll(limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get all groups")
		return nil, 0, err
	}

	return groups, total, nil
}

// UpdateGroup updates a group's name, description, or role
func (app *GroupApp) UpdateGroup(groupId, name, description string, role domain.Role) error {
	logger := app.Logger.With().Str("component", "application.group.update").Logger()

	group, err := app.GroupPers.GetById(groupId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get group")
		return err
	}

	group.Name = name
	group.Description = description
	group.Role = role

	if err := app.GroupPers.Update(group); err != nil {
		logger.Error().Err(err).Msg("failed to update group")
		return err
	}

	return nil
}

// DeleteGroup deletes a group
func (app *GroupApp) DeleteGroup(groupId string) error {
	logger := app.Logger.With().Str("component", "application.group.delete").Logger()

	if err := app.GroupPers.Delete(groupId); err != nil {
		logger.Error().Err(err).Msg("failed to delete group")
		return err
	}

	return nil
}

// AddMember adds a user to a group
func (app *GroupApp) AddMember(groupId, userId string) error {
	logger := app.Logger.With().Str("component", "application.group.add_member").Logger()

	// Verify user exists
	_, err := app.UserPers.GetById(userId)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userId).Msg("user not found")
		return err
	}

	if err := app.GroupPers.AddMember(groupId, userId); err != nil {
		logger.Error().Err(err).Msg("failed to add member to group")
		return err
	}

	return nil
}

// RemoveMember removes a user from a group
func (app *GroupApp) RemoveMember(groupId, userId string) error {
	logger := app.Logger.With().Str("component", "application.group.remove_member").Logger()

	if err := app.GroupPers.RemoveMember(groupId, userId); err != nil {
		logger.Error().Err(err).Msg("failed to remove member from group")
		return err
	}

	return nil
}

// GetMembers retrieves all members of a group
func (app *GroupApp) GetMembers(groupId string) ([]domain.User, error) {
	return app.GroupPers.GetMembers(groupId)
}
