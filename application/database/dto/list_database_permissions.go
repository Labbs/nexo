package dto

type ListDatabasePermissionsInput struct {
	UserId     string
	DatabaseId string
}

type ListDatabasePermissionsOutput struct {
	Permissions []DatabasePermissionItem `json:"permissions"`
}
