package dto

type DeleteDrawingUserPermissionInput struct {
	RequesterId  string
	DrawingId    string
	TargetUserId string
}
