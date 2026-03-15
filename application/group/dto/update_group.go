package dto

import "github.com/labbs/nexo/domain"

type UpdateGroupInput struct {
	GroupId     string
	Name        string
	Description string
	Role        domain.Role
}
