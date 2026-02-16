package dto

import "github.com/labbs/nexo/domain"

type GetMembersInput struct {
	GroupId string
}

type GetMembersOutput struct {
	Members []domain.User
}
