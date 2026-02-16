package dto

import "github.com/labbs/nexo/domain"

type CreateGroupInput struct {
	Name        string
	Description string
	OwnerId     string
	Role        domain.Role
}

type CreateGroupOutput struct {
	Group *domain.Group
}
