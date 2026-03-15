package dto

import "time"

type GetDatabaseInput struct {
	UserId     string
	DatabaseId string
}

type GetDatabaseOutput struct {
	Id          string
	SpaceId     string
	DocumentId  *string
	Name        string
	Description string
	Icon        string
	Schema      []PropertySchema
	Views       []ViewConfig
	DefaultView string
	Type        string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
