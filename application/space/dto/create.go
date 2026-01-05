package dto

import "github.com/labbs/nexo/domain"

type CreatePrivateSpaceForUserInput struct {
	UserId string
}

type CreatePrivateSpaceForUserOutput struct {
	Space *domain.Space
}

type CreatePublicSpaceInput struct {
	Name      string
	Icon      *string
	IconColor *string
	OwnerId   *string
}

type CreatePublicSpaceOutput struct {
	Space *domain.Space
}
