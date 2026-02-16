package dto

type UpdateDrawingInput struct {
	UserId    string
	DrawingId string
	Name      *string
	Icon      *string
	Elements  []interface{}
	AppState  map[string]interface{}
	Files     map[string]interface{}
	Thumbnail *string
}
