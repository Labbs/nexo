package dto

type CreateViewInput struct {
	UserId     string
	DatabaseId string
	Name       string
	Type       string
	Filter     map[string]any
	Sort       []SortConfig
	Columns    []string
}

type CreateViewOutput struct {
	Id      string
	Name    string
	Type    string
	Filter  map[string]any
	Sort    []SortConfig
	Columns []string
}
