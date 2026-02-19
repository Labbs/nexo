package dto

type UpsertDrawingUserPermissionInput struct {
	RequesterId  string
	DrawingId    string
	TargetUserId string
	Role         string
}
