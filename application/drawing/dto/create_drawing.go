package dto

import "time"

type CreateDrawingInput struct {
	UserId     string
	SpaceId    string
	DocumentId *string
	Name       string
	Icon       string
	Elements   []interface{}
	AppState   map[string]interface{}
	Files      map[string]interface{}
	Thumbnail  string
}

type CreateDrawingOutput struct {
	Id        string
	Name      string
	CreatedAt time.Time
}
