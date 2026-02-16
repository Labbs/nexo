package dto

type ListDatabasesInput struct {
	UserId  string
	SpaceId string
}

type ListDatabasesOutput struct {
	Databases []DatabaseItem
}
