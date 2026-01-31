package dto

import "github.com/labbs/nexo/domain"

type GetByEmailInput struct {
	Email string
}

type GetByEmailOutput struct {
	User *domain.User
}

type GetByUserIdInput struct {
	UserId string
}

type GetByUserIdOutput struct {
	User *domain.User
}
