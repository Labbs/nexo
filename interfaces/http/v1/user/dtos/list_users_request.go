package dtos

// List users (simplified - for use in person picker)

type ListUsersRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type UserListItem struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

type ListUsersResponse struct {
	Users      []UserListItem `json:"users"`
	TotalCount int64          `json:"total_count"`
}
