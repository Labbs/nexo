package dto

import "github.com/labbs/nexo/domain"

type ListDrawingPermissionsInput struct {
	RequesterId string
	DrawingId   string
}

type ListDrawingPermissionsOutput struct {
	Permissions []domain.Permission
}
