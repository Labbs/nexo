package dto

type UpdateViewInput struct {
	UserId        string
	DatabaseId    string
	ViewId        string
	Name          *string
	Type          *string
	Filter        map[string]interface{}
	Sort          []SortConfig
	Columns       []string
	HiddenColumns []string
	GroupBy       *string
}
