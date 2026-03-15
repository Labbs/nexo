package dto

type UpdateRowInput struct {
	UserId        string
	DatabaseId    string
	RowId         string
	Properties    map[string]any
	Content       map[string]any
	ShowInSidebar *bool
}
