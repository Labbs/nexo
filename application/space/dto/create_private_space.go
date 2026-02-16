package dto

import "github.com/labbs/nexo/domain"

type CreatePrivateSpaceForUserInput struct {
	UserId string
}

type CreatePrivateSpaceForUserOutput struct {
	Space *domain.Space
}
