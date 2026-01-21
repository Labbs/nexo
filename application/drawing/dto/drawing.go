package dto

import "time"

// Create drawing
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

// List drawings
type ListDrawingsInput struct {
	UserId  string
	SpaceId string
}

type DrawingItem struct {
	Id         string
	DocumentId *string
	Name       string
	Icon       string
	Thumbnail  string
	CreatedBy  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ListDrawingsOutput struct {
	Drawings []DrawingItem
}

// Get drawing
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

// Update drawing
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

// Delete drawing
type DeleteDrawingInput struct {
	UserId    string
	DrawingId string
}
