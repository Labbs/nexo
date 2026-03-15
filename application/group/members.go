package group

import (
	"github.com/labbs/nexo/application/group/dto"
)

// GetMembers retrieves all members of a group
func (app *GroupApplication) GetMembers(input dto.GetMembersInput) (*dto.GetMembersOutput, error) {
	members, err := app.GroupPers.GetMembers(input.GroupId)
	if err != nil {
		return nil, err
	}

	return &dto.GetMembersOutput{Members: members}, nil
}
