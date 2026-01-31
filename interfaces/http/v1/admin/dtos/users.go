package dtos

import "time"

// User list

type ListUsersRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type UserItem struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"avatar_url"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersResponse struct {
	Users      []UserItem `json:"users"`
	TotalCount int64      `json:"total_count"`
}

// Update user role

type UpdateUserRoleRequest struct {
	UserId string `path:"user_id"`
	Role   string `json:"role" validate:"required,oneof=user admin guest"`
}

type UpdateUserRoleResponse struct {
	Message string `json:"message"`
}

// Update user active status

type UpdateUserActiveRequest struct {
	UserId string `path:"user_id"`
	Active bool   `json:"active"`
}

type UpdateUserActiveResponse struct {
	Message string `json:"message"`
}

// Delete user

type DeleteUserRequest struct {
	UserId string `path:"user_id"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

// Invite user

type InviteUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=user admin guest"`
}

type InviteUserResponse struct {
	Message string `json:"message"`
	UserId  string `json:"user_id,omitempty"`
}
