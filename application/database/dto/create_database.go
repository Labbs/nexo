package dto

import "time"

type CreateDatabaseInput struct {
	UserId      string
	SpaceId     string
	DocumentId  *string
	Name        string
	Description string
	Icon        string
	Schema      []PropertySchema
	Type        string // "spreadsheet" or "document", defaults to "spreadsheet"
}

type CreateDatabaseOutput struct {
	Id          string
	Name        string
	Description string
	Icon        string
	Schema      []PropertySchema
	DefaultView string
	Type        string
	CreatedAt   time.Time
}
