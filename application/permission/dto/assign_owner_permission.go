package dto

type AssignOwnerPermissionInput struct {
	ResourceType string
	ResourceId   string
	UserId       string
	Role         string
}
