package dtos

// Permission item
type DatabasePermissionItem struct {
	Id        string  `json:"id"`
	UserId    *string `json:"user_id,omitempty"`
	Username  *string `json:"username,omitempty"`
	GroupId   *string `json:"group_id,omitempty"`
	GroupName *string `json:"group_name,omitempty"`
	Role      string  `json:"role"`
}

// List permissions request/response
type ListDatabasePermissionsRequest struct {
	DatabaseId string `path:"database_id"`
}

type ListDatabasePermissionsResponse struct {
	Permissions []DatabasePermissionItem `json:"permissions"`
}

// Upsert permission request
type UpsertDatabasePermissionRequest struct {
	DatabaseId string  `path:"database_id"`
	UserId     *string `json:"user_id,omitempty"`
	GroupId    *string `json:"group_id,omitempty"`
	Role       string  `json:"role"`
}

type UpsertDatabasePermissionResponse struct {
	Success bool `json:"success"`
}

// Delete permission request
type DeleteDatabasePermissionRequest struct {
	DatabaseId string  `path:"database_id"`
	UserId     *string `query:"user_id"`
	GroupId    *string `query:"group_id"`
}

type DeleteDatabasePermissionResponse struct {
	Success bool `json:"success"`
}
