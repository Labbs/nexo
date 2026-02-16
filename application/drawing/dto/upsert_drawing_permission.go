package dto

import "github.com/labbs/nexo/domain"

type UpsertDrawingUserPermissionInput struct {
	RequesterId  string
	DrawingId    string
	TargetUserId string
	Role         domain.PermissionRole
}
