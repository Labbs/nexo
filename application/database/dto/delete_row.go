package dto

type DeleteRowInput struct {
	UserId     string
	DatabaseId string
	RowId      string
}

type BulkDeleteRowsInput struct {
	UserId     string
	DatabaseId string
	RowIds     []string
}
