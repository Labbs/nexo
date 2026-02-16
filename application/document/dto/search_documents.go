package dto

type SearchInput struct {
	UserId  string
	Query   string
	SpaceId *string
	Limit   int
}

type SearchOutput struct {
	Results []SearchResultItem
}
