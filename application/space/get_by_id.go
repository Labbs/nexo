package space

import (
	"fmt"

	"github.com/labbs/nexo/application/space/dto"
)

func (c *SpaceApplication) GetSpaceById(input dto.GetSpaceByIdInput) (*dto.GetSpaceByIdOutput, error) {
	logger := c.Logger.With().Str("component", "application.space.get_space_by_id").Logger()

	space, err := c.SpacePres.GetSpaceById(input.SpaceId)
	if err != nil {
		logger.Error().Err(err).Str("spaceId", input.SpaceId).Msg("failed to get space by id")
		return nil, fmt.Errorf("space not found: %w", err)
	}

	// Map domain.Space → SpaceDetail DTO
	permissions := make([]dto.SpacePermission, len(space.Permissions))
	for i, p := range space.Permissions {
		permissions[i] = dto.SpacePermission{
			UserId:  p.UserId,
			GroupId: p.GroupId,
			Role:    string(p.Role),
		}
	}

	detail := &dto.SpaceDetail{
		Id:          space.Id,
		Name:        space.Name,
		Slug:        space.Slug,
		Icon:        space.Icon,
		IconColor:   space.IconColor,
		Type:        string(space.Type),
		OwnerId:     space.OwnerId,
		Permissions: permissions,
		CreatedAt:   space.CreatedAt,
		UpdatedAt:   space.UpdatedAt,
	}

	return &dto.GetSpaceByIdOutput{Space: detail}, nil
}
