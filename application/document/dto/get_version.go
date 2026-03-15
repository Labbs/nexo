package dto

import "time"

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
