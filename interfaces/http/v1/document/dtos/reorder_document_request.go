package dtos

type ReorderItem struct {
	Id       string `json:"id" validate:"required"`
	Position int    `json:"position" validate:"min=0"`
}

type ReorderDocumentsRequest struct {
	SpaceId string        `path:"space_id" validate:"required,uuid4"`
	Items   []ReorderItem `json:"items" validate:"required,min=1"`
}

type ReorderDocumentsResponse struct {
	Message string `json:"message"`
}
