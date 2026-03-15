package dto

import "github.com/labbs/nexo/domain"

type GetGroupInput struct {
	GroupId string
}

type GetGroupOutput struct {
	Group *domain.Group
}
