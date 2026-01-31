package space

import (
	"fmt"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/gosimple/slug"
	"github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/shortuuid"
)

func (c *SpaceApp) CreatePrivateSpaceForUser(input dto.CreatePrivateSpaceForUserInput) (*dto.CreatePrivateSpaceForUserOutput, error) {
	logger := c.Logger.With().Str("component", "application.space.createPrivateSpaceForUser").Logger()

	name := "Personal Space"

	space := &domain.Space{
		Id:      utils.UUIDv4(),
		Slug:    slug.Make(name + "-" + shortuuid.GenerateShortUUID()),
		Name:    name,
		Type:    domain.SpaceTypePersonal,
		OwnerId: &input.UserId,
	}

	err := c.SpacePres.Create(space)
	if err != nil {
		logger.Error().Err(err).Str("userId", input.UserId).Msg("failed to create private space for user")
		return nil, fmt.Errorf("failed to create private space for user: %w", err)
	}

	// Auto-create owner permission for the creator
	if err := c.PermissionPers.UpsertUser(domain.PermissionTypeSpace, space.Id, input.UserId, domain.PermissionRoleOwner); err != nil {
		logger.Warn().Err(err).Str("space_id", space.Id).Str("user_id", input.UserId).Msg("failed to create owner permission")
	}

	return &dto.CreatePrivateSpaceForUserOutput{Space: space}, nil
}

func (c *SpaceApp) CreateSpace(input dto.CreateSpaceInput) (*dto.CreateSpaceOutput, error) {
	logger := c.Logger.With().Str("component", "application.space.createSpace").Logger()

	space := &domain.Space{
		Id:        utils.UUIDv4(),
		Slug:      slug.Make(input.Name + "-" + shortuuid.GenerateShortUUID()),
		Name:      input.Name,
		Type:      input.Type,
		OwnerId:   input.OwnerId,
		Icon:      *input.Icon,
		IconColor: *input.IconColor,
	}

	err := c.SpacePres.Create(space)
	if err != nil {
		logger.Error().Err(err).Str("name", input.Name).Msg("failed to create space")
		return nil, fmt.Errorf("failed to create space: %w", err)
	}

	// Auto-create owner permission for the creator
	if input.OwnerId != nil {
		if err := c.PermissionPers.UpsertUser(domain.PermissionTypeSpace, space.Id, *input.OwnerId, domain.PermissionRoleOwner); err != nil {
			logger.Warn().Err(err).Str("space_id", space.Id).Str("user_id", *input.OwnerId).Msg("failed to create owner permission")
		}
	}

	return &dto.CreateSpaceOutput{Space: space}, nil
}
