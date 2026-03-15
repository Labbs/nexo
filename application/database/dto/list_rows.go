package dto

type ListRowsInput struct {
	UserId     string
	DatabaseId string
	ViewId     string
	Limit      int
	Offset     int
}

type ListRowsOutput struct {
	Rows       []RowItem
	TotalCount int64
}
