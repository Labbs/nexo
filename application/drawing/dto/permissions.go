package dto

import "github.com/labbs/nexo/domain"

type ListDrawingPermissionsInput struct {
	RequesterId string
	DrawingId   string
}

type ListDrawingPermissionsOutput struct {
	Permissions []domain.Permission
}

type UpsertDrawingUserPermissionInput struct {
	RequesterId  string
	DrawingId    string
	TargetUserId string
	Role         domain.PermissionRole
}

type DeleteDrawingUserPermissionInput struct {
	RequesterId  string
	DrawingId    string
	TargetUserId string
}
