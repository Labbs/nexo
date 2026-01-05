package dto

import "github.com/labbs/nexo/domain"

type CreateUserInput struct {
	User domain.User
}

type CreateUserOutput struct {
	User *domain.User
}
