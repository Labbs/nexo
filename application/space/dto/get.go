package dto

import "github.com/labbs/nexo/domain"

type GetSpacesForUserInput struct {
	UserId string
}

type GetSpacesForUserOutput struct {
	Spaces []domain.Space
}
