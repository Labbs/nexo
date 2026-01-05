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

	name := "Private Space"

	space := &domain.Space{
		Id:      utils.UUIDv4(),
		Slug:    slug.Make(name + "-" + shortuuid.GenerateShortUUID()),
		Name:    name,
		Type:    domain.SpaceTypePrivate,
		OwnerId: &input.UserId,
	}

	err := c.SpacePres.Create(space)
	if err != nil {
		logger.Error().Err(err).Str("userId", input.UserId).Msg("failed to create private space for user")
		return nil, fmt.Errorf("failed to create private space for user: %w", err)
	}

	return &dto.CreatePrivateSpaceForUserOutput{Space: space}, nil
}

func (c *SpaceApp) CreatePublicSpace(input dto.CreatePublicSpaceInput) (*dto.CreatePublicSpaceOutput, error) {
	logger := c.Logger.With().Str("component", "application.space.createPublicSpace").Logger()

	space := &domain.Space{
		Id:        utils.UUIDv4(),
		Slug:      slug.Make(input.Name + "-" + shortuuid.GenerateShortUUID()),
		Name:      input.Name,
		Type:      domain.SpaceTypePublic,
		OwnerId:   input.OwnerId,
		Icon:      *input.Icon,
		IconColor: *input.IconColor,
	}

	err := c.SpacePres.Create(space)
	if err != nil {
		logger.Error().Err(err).Str("name", input.Name).Msg("failed to create public space")
		return nil, fmt.Errorf("failed to create public space: %w", err)
	}

	return &dto.CreatePublicSpaceOutput{Space: space}, nil
}
