package dto

type DeleteDatabasePermissionInput struct {
	UserId       string  // The user making the request
	DatabaseId   string  // The database
	TargetUserId *string // The user to remove permission from (mutually exclusive with GroupId)
	GroupId      *string // The group to remove permission from (mutually exclusive with TargetUserId)
}
