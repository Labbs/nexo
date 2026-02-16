package dto

type ReorderDocumentsInput struct {
	UserId  string
	SpaceId string
	Items   []ReorderItem
}
