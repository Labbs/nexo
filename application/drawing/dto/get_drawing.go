package dto

import "time"

type GetDrawingInput struct {
	UserId    string
	DrawingId string
}

type GetDrawingOutput struct {
	Id         string
	SpaceId    string
	DocumentId *string
	Name       string
	Icon       string
	Elements   []any
	AppState   map[string]any
	Files      map[string]any
	Thumbnail  string
	CreatedBy  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
