package dto

type UpdateRowInput struct {
	UserId        string
	DatabaseId    string
	RowId         string
	Properties    map[string]interface{}
	Content       map[string]interface{}
	ShowInSidebar *bool
}
