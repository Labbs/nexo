package dto

import "github.com/labbs/nexo/domain"

type GetAllGroupsInput struct {
	Limit  int
	Offset int
}

type GetAllGroupsOutput struct {
	Groups     []domain.Group
	TotalCount int64
}
