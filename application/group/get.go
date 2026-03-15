package group

import (
	"github.com/labbs/nexo/application/group/dto"
)

// GetGroup retrieves a group by ID
func (app *GroupApplication) GetGroup(input dto.GetGroupInput) (*dto.GetGroupOutput, error) {
	group, err := app.GroupPers.GetById(input.GroupId)
	if err != nil {
		return nil, err
	}

	return &dto.GetGroupOutput{Group: group}, nil
}
