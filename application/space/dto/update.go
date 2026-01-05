package dto

import "github.com/labbs/nexo/domain"

type UpdateSpaceInput struct {
	UserId    string
	SpaceId   string
	Name      *string
	Icon      *string
	IconColor *string
}

type UpdateSpaceOutput struct {
	Space *domain.Space
}
