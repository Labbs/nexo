package dto

type DeleteSpaceUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	TargetUserId string
}
