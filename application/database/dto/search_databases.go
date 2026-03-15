package dto

import "time"

type SearchDatabasesInput struct {
	UserId  string
	Query   string
	SpaceId *string
	Limit   int
}

type SearchDatabaseResultItem struct {
	Id          string
	Name        string
	Description string
	Icon        string
	Type        string
	SpaceId     string
	SpaceName   string
	UpdatedAt   time.Time
}

type SearchDatabasesOutput struct {
	Results []SearchDatabaseResultItem
}
