package dto

type UpsertSpaceGroupPermissionInput struct {
	RequesterId string
	SpaceId     string
	GroupId     string
	Role        string
}
