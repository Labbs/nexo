package dtos

import "time"

// ListGroupsRequest is the request to list all groups
type ListGroupsRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// ListGroupsResponse is the response for listing all groups
type ListGroupsResponse struct {
	Groups     []GroupItem `json:"groups"`
	TotalCount int64       `json:"total_count"`
}

// GroupItem represents a group in list responses
type GroupItem struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Role        string       `json:"role"`
	OwnerId     string       `json:"owner_id"`
	OwnerName   string       `json:"owner_name,omitempty"`
	MemberCount int          `json:"member_count"`
	Members     []MemberItem `json:"members,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// MemberItem represents a group member
type MemberItem struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}

// CreateGroupRequest is the request to create a group
type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
	Role        string `json:"role" validate:"required,oneof=user admin guest"`
}

// CreateGroupResponse is the response for creating a group
type CreateGroupResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// UpdateGroupRequest is the request to update a group
type UpdateGroupRequest struct {
	GroupId     string `path:"group_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
	Role        string `json:"role" validate:"required,oneof=user admin guest"`
}

// UpdateGroupResponse is the response for updating a group
type UpdateGroupResponse struct {
	Message string `json:"message"`
}

// DeleteGroupRequest is the request to delete a group
type DeleteGroupRequest struct {
	GroupId string `path:"group_id" validate:"required"`
}

// DeleteGroupResponse is the response for deleting a group
type DeleteGroupResponse struct {
	Message string `json:"message"`
}

// AddGroupMemberRequest is the request to add a member to a group
type AddGroupMemberRequest struct {
	GroupId string `path:"group_id" validate:"required"`
	UserId  string `json:"user_id" validate:"required"`
}

// AddGroupMemberResponse is the response for adding a member
type AddGroupMemberResponse struct {
	Message string `json:"message"`
}

// RemoveGroupMemberRequest is the request to remove a member from a group
type RemoveGroupMemberRequest struct {
	GroupId string `path:"group_id" validate:"required"`
	UserId  string `path:"user_id" validate:"required"`
}

// RemoveGroupMemberResponse is the response for removing a member
type RemoveGroupMemberResponse struct {
	Message string `json:"message"`
}

// GetGroupMembersRequest is the request to get group members
type GetGroupMembersRequest struct {
	GroupId string `path:"group_id" validate:"required"`
}

// GetGroupMembersResponse is the response for getting group members
type GetGroupMembersResponse struct {
	Members []MemberItem `json:"members"`
}
