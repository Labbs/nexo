package dto

import "github.com/labbs/nexo/domain"

type GetByUserIdInput struct {
	UserId string
}

type GetByUserIdOutput struct {
	User *domain.User
}
