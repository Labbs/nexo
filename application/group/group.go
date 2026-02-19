package group

import (
	"github.com/labbs/nexo/application/group/dto"
	"github.com/labbs/nexo/application/ports"
	userDto "github.com/labbs/nexo/application/user/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type GroupApp struct {
	Config    config.Config
	Logger    zerolog.Logger
	GroupPers domain.GroupPers
	UserApp   ports.UserPort
}

func NewGroupApp(config config.Config, logger zerolog.Logger, groupPers domain.GroupPers) *GroupApp {
	return &GroupApp{
		Config:    config,
		Logger:    logger,
		GroupPers: groupPers,
	}
}

// CreateGroup creates a new group
func (app *GroupApp) CreateGroup(input dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
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

// GetGroup retrieves a group by ID
func (app *GroupApp) GetGroup(input dto.GetGroupInput) (*dto.GetGroupOutput, error) {
	group, err := app.GroupPers.GetById(input.GroupId)
	if err != nil {
		return nil, err
	}

	return &dto.GetGroupOutput{Group: group}, nil
}

// GetAllGroups retrieves all groups with pagination
func (app *GroupApp) GetAllGroups(input dto.GetAllGroupsInput) (*dto.GetAllGroupsOutput, error) {
	logger := app.Logger.With().Str("component", "application.group.get_all").Logger()

	groups, total, err := app.GroupPers.GetAll(input.Limit, input.Offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get all groups")
		return nil, err
	}

	return &dto.GetAllGroupsOutput{Groups: groups, TotalCount: total}, nil
}

// UpdateGroup updates a group's name, description, or role
func (app *GroupApp) UpdateGroup(input dto.UpdateGroupInput) error {
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

// DeleteGroup deletes a group
func (app *GroupApp) DeleteGroup(input dto.DeleteGroupInput) error {
	logger := app.Logger.With().Str("component", "application.group.delete").Logger()

	if err := app.GroupPers.Delete(input.GroupId); err != nil {
		logger.Error().Err(err).Msg("failed to delete group")
		return err
	}

	return nil
}

// AddMember adds a user to a group
func (app *GroupApp) AddMember(input dto.AddMemberInput) error {
	logger := app.Logger.With().Str("component", "application.group.add_member").Logger()

	// Verify user exists
	_, err := app.UserApp.GetByUserId(userDto.GetByUserIdInput{UserId: input.UserId})
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

// RemoveMember removes a user from a group
func (app *GroupApp) RemoveMember(input dto.RemoveMemberInput) error {
	logger := app.Logger.With().Str("component", "application.group.remove_member").Logger()

	if err := app.GroupPers.RemoveMember(input.GroupId, input.UserId); err != nil {
		logger.Error().Err(err).Msg("failed to remove member from group")
		return err
	}

	return nil
}

// GetMembers retrieves all members of a group
func (app *GroupApp) GetMembers(input dto.GetMembersInput) (*dto.GetMembersOutput, error) {
	members, err := app.GroupPers.GetMembers(input.GroupId)
	if err != nil {
		return nil, err
	}

	return &dto.GetMembersOutput{Members: members}, nil
}
