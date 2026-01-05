package dtos

import "time"

type SearchRequest struct {
	Query   string  `query:"q"`
	SpaceId *string `query:"space_id"`
	Limit   int     `query:"limit"`
}

type SearchResultItem struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	SpaceId   string    `json:"space_id"`
	SpaceName string    `json:"space_name"`
	Icon      string    `json:"icon,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchResponse struct {
	Results []SearchResultItem `json:"results"`
}
