package dto

type ListVersionsInput struct {
	UserId     string
	SpaceId    string
	DocumentId string
	Limit      int
	Offset     int
}

type ListVersionsOutput struct {
	Versions   []VersionItem
	TotalCount int64
}
