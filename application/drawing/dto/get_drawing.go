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
	Elements   []interface{}
	AppState   map[string]interface{}
	Files      map[string]interface{}
	Thumbnail  string
	CreatedBy  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
