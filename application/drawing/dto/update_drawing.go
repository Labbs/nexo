package dto

type UpdateDrawingInput struct {
	UserId    string
	DrawingId string
	Name      *string
	Icon      *string
	Elements  []any
	AppState  map[string]any
	Files     map[string]any
	Thumbnail *string
}
