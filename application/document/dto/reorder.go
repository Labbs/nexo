package dto

type ReorderItem struct {
	Id       string
	Position int
}

type ReorderDocumentsInput struct {
	UserId  string
	SpaceId string
	Items   []ReorderItem
}
