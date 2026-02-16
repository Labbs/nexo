package dto

import "github.com/labbs/nexo/domain"

type CreateSpaceInput struct {
	Name      string
	Icon      *string
	IconColor *string
	OwnerId   *string
	Type      domain.SpaceType
}

type CreateSpaceOutput struct {
	Space *domain.Space
}
