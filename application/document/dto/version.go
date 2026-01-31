package dto

import "time"

// List versions
type ListVersionsInput struct {
	UserId     string
	SpaceId    string
	DocumentId string
	Limit      int
	Offset     int
}

type VersionItem struct {
	Id          string
	Version     int
	Name        string
	Description string
	UserId      string
	UserName    string
	CreatedAt   time.Time
}

type ListVersionsOutput struct {
	Versions   []VersionItem
	TotalCount int64
}

// Get version
type GetVersionInput struct {
	UserId    string
	VersionId string
}

type GetVersionOutput struct {
	Id          string
	Version     int
	DocumentId  string
	Name        string
	Content     []Block
	Config      DocumentConfig
	Description string
	UserId      string
	UserName    string
	CreatedAt   time.Time
}

// Restore version
type RestoreVersionInput struct {
	UserId    string
	VersionId string
}

// Create version manually
type CreateVersionInput struct {
	UserId      string
	DocumentId  string
	Description string
}

type CreateVersionOutput struct {
	VersionId string
	Version   int
}
