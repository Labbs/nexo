package space

import (
	"fmt"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/gosimple/slug"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/shortuuid"
)

// GetAllSpaces returns all spaces with pagination (admin only)
func (c *SpaceApp) GetAllSpaces(limit, offset int) ([]domain.Space, int64, error) {
	logger := c.Logger.With().Str("component", "application.space.get_all_spaces").Logger()

	spaces, total, err := c.SpacePres.GetAll(limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get all spaces")
		return nil, 0, err
	}

	return spaces, total, nil
}

// AdminCreateSpace creates a space without requiring an owner (admin only)
func (c *SpaceApp) AdminCreateSpace(name, icon, iconColor string, spaceType domain.SpaceType, ownerId *string) (*domain.Space, error) {
	logger := c.Logger.With().Str("component", "application.space.admin_create_space").Logger()

	space := &domain.Space{
		Id:        utils.UUIDv4(),
		Slug:      slug.Make(name + "-" + shortuuid.GenerateShortUUID()),
		Name:      name,
		Type:      spaceType,
		OwnerId:   ownerId,
		Icon:      icon,
		IconColor: iconColor,
	}

	err := c.SpacePres.Create(space)
	if err != nil {
		logger.Error().Err(err).Str("name", name).Msg("failed to create space")
		return nil, fmt.Errorf("failed to create space: %w", err)
	}

	return space, nil
}

// AdminUpdateSpace updates a space (admin only)
func (c *SpaceApp) AdminUpdateSpace(spaceId, name, icon, iconColor string, spaceType domain.SpaceType, ownerId *string) error {
	logger := c.Logger.With().Str("component", "application.space.admin_update_space").Logger()

	space, err := c.SpacePres.GetSpaceById(spaceId)
	if err != nil {
		return fmt.Errorf("space not found")
	}

	space.Name = name
	space.Icon = icon
	space.IconColor = iconColor
	space.Type = spaceType
	space.OwnerId = ownerId

	err = c.SpacePres.Update(space)
	if err != nil {
		logger.Error().Err(err).Str("spaceId", spaceId).Msg("failed to update space")
		return fmt.Errorf("failed to update space: %w", err)
	}

	return nil
}

// AdminDeleteSpace deletes a space (admin only)
func (c *SpaceApp) AdminDeleteSpace(spaceId string) error {
	logger := c.Logger.With().Str("component", "application.space.admin_delete_space").Logger()

	err := c.SpacePres.Delete(spaceId)
	if err != nil {
		logger.Error().Err(err).Str("spaceId", spaceId).Msg("failed to delete space")
		return fmt.Errorf("failed to delete space: %w", err)
	}

	return nil
}

// AdminListSpacePermissions lists all permissions for a space with user/group details (admin only)
func (c *SpaceApp) AdminListSpacePermissions(spaceId string) ([]domain.Permission, error) {
	logger := c.Logger.With().Str("component", "application.space.admin_list_permissions").Logger()

	permissions, err := c.PermissionPers.ListByResource(domain.PermissionTypeSpace, spaceId)
	if err != nil {
		logger.Error().Err(err).Str("spaceId", spaceId).Msg("failed to list permissions")
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	return permissions, nil
}

// AdminAddSpaceUserPermission adds a user permission to a space (admin only)
func (c *SpaceApp) AdminAddSpaceUserPermission(spaceId, userId string, role domain.PermissionRole) error {
	logger := c.Logger.With().Str("component", "application.space.admin_add_user_permission").Logger()

	err := c.PermissionPers.UpsertUser(domain.PermissionTypeSpace, spaceId, userId, role)
	if err != nil {
		logger.Error().Err(err).Str("spaceId", spaceId).Str("userId", userId).Msg("failed to add user permission")
		return fmt.Errorf("failed to add user permission: %w", err)
	}

	return nil
}

// AdminRemoveSpaceUserPermission removes a user permission from a space (admin only)
func (c *SpaceApp) AdminRemoveSpaceUserPermission(spaceId, userId string) error {
	logger := c.Logger.With().Str("component", "application.space.admin_remove_user_permission").Logger()

	err := c.PermissionPers.DeleteUser(domain.PermissionTypeSpace, spaceId, userId)
	if err != nil {
		logger.Error().Err(err).Str("spaceId", spaceId).Str("userId", userId).Msg("failed to remove user permission")
		return fmt.Errorf("failed to remove user permission: %w", err)
	}

	return nil
}

// AdminAddSpaceGroupPermission adds a group permission to a space (admin only)
func (c *SpaceApp) AdminAddSpaceGroupPermission(spaceId, groupId string, role domain.PermissionRole) error {
	logger := c.Logger.With().Str("component", "application.space.admin_add_group_permission").Logger()

	err := c.PermissionPers.UpsertGroup(domain.PermissionTypeSpace, spaceId, groupId, role)
	if err != nil {
		logger.Error().Err(err).Str("spaceId", spaceId).Str("groupId", groupId).Msg("failed to add group permission")
		return fmt.Errorf("failed to add group permission: %w", err)
	}

	return nil
}

// AdminRemoveSpaceGroupPermission removes a group permission from a space (admin only)
func (c *SpaceApp) AdminRemoveSpaceGroupPermission(spaceId, groupId string) error {
	logger := c.Logger.With().Str("component", "application.space.admin_remove_group_permission").Logger()

	err := c.PermissionPers.DeleteGroup(domain.PermissionTypeSpace, spaceId, groupId)
	if err != nil {
		logger.Error().Err(err).Str("spaceId", spaceId).Str("groupId", groupId).Msg("failed to remove group permission")
		return fmt.Errorf("failed to remove group permission: %w", err)
	}

	return nil
}
