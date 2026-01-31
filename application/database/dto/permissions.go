package dto

// Permission item for listing
type DatabasePermissionItem struct {
	Id        string  `json:"id"`
	UserId    *string `json:"user_id,omitempty"`
	Username  *string `json:"username,omitempty"`
	GroupId   *string `json:"group_id,omitempty"`
	GroupName *string `json:"group_name,omitempty"`
	Role      string  `json:"role"`
}

// List permissions
type ListDatabasePermissionsInput struct {
	UserId     string
	DatabaseId string
}

type ListDatabasePermissionsOutput struct {
	Permissions []DatabasePermissionItem `json:"permissions"`
}

// Upsert permission
type UpsertDatabasePermissionInput struct {
	UserId       string  // The user making the request
	DatabaseId   string  // The database to add permission to
	TargetUserId *string // The user to add permission for (mutually exclusive with GroupId)
	GroupId      *string // The group to add permission for (mutually exclusive with TargetUserId)
	Role         string  // editor, viewer, denied
}

// Delete permission
type DeleteDatabasePermissionInput struct {
	UserId       string  // The user making the request
	DatabaseId   string  // The database
	TargetUserId *string // The user to remove permission from (mutually exclusive with GroupId)
	GroupId      *string // The group to remove permission from (mutually exclusive with TargetUserId)
}
