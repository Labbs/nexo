package dto

import "time"

type CreateDrawingInput struct {
	UserId     string
	SpaceId    string
	DocumentId *string
	Name       string
	Icon       string
	Elements   []any
	AppState   map[string]any
	Files      map[string]any
	Thumbnail  string
}

type CreateDrawingOutput struct {
	Id        string
	Name      string
	CreatedAt time.Time
}
