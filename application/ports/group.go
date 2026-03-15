package ports

import (
	"github.com/labbs/nexo/application/group/dto"
)

type GroupPort interface {
	CreateGroup(input dto.CreateGroupInput) (*dto.CreateGroupOutput, error)
	GetGroup(input dto.GetGroupInput) (*dto.GetGroupOutput, error)
	GetAllGroups(input dto.GetAllGroupsInput) (*dto.GetAllGroupsOutput, error)
	UpdateGroup(input dto.UpdateGroupInput) error
	DeleteGroup(input dto.DeleteGroupInput) error
	AddMember(input dto.AddMemberInput) error
	RemoveMember(input dto.RemoveMemberInput) error
	GetMembers(input dto.GetMembersInput) (*dto.GetMembersOutput, error)
}
