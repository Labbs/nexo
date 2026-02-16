package dto

type UpsertDatabasePermissionInput struct {
	UserId       string  // The user making the request
	DatabaseId   string  // The database to add permission to
	TargetUserId *string // The user to add permission for (mutually exclusive with GroupId)
	GroupId      *string // The group to add permission for (mutually exclusive with TargetUserId)
	Role         string  // editor, viewer, denied
}
