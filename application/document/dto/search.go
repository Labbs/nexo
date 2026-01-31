package dto

import "time"

type SearchInput struct {
	UserId  string
	Query   string
	SpaceId *string
	Limit   int
}

type SearchResultItem struct {
	Id        string
	Name      string
	Slug      string
	SpaceId   string
	SpaceName string
	Icon      string
	UpdatedAt time.Time
}

type SearchOutput struct {
	Results []SearchResultItem
}
