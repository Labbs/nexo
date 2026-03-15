package dto

type UpsertSpaceUserPermissionInput struct {
	RequesterId  string
	SpaceId      string
	TargetUserId string
	Role         string
}
