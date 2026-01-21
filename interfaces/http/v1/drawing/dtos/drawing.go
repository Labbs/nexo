package dtos

import "time"

// Request DTOs

type CreateDrawingRequest struct {
	SpaceId    string                 `json:"space_id"`
	DocumentId *string                `json:"document_id,omitempty"`
	Name       string                 `json:"name"`
	Icon       string                 `json:"icon,omitempty"`
	Elements   []interface{}          `json:"elements,omitempty"`
	AppState   map[string]interface{} `json:"app_state,omitempty"`
	Files      map[string]interface{} `json:"files,omitempty"`
	Thumbnail  string                 `json:"thumbnail,omitempty"`
}

type ListDrawingsRequest struct {
	SpaceId string `query:"space_id"`
}

type GetDrawingRequest struct {
	DrawingId string `path:"drawing_id"`
}

type UpdateDrawingRequest struct {
	DrawingId string                 `path:"drawing_id"`
	Name      *string                `json:"name,omitempty"`
	Icon      *string                `json:"icon,omitempty"`
	Elements  []interface{}          `json:"elements,omitempty"`
	AppState  map[string]interface{} `json:"app_state,omitempty"`
	Files     map[string]interface{} `json:"files,omitempty"`
	Thumbnail *string                `json:"thumbnail,omitempty"`
}

type DeleteDrawingRequest struct {
	DrawingId string `path:"drawing_id"`
}

// Response DTOs

type MessageResponse struct {
	Message string `json:"message"`
}

type CreateDrawingResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type DrawingItem struct {
	Id         string    `json:"id"`
	DocumentId *string   `json:"document_id,omitempty"`
	Name       string    `json:"name"`
	Icon       string    `json:"icon,omitempty"`
	Thumbnail  string    `json:"thumbnail,omitempty"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ListDrawingsResponse struct {
	Drawings []DrawingItem `json:"drawings"`
}

type GetDrawingResponse struct {
	Id         string                 `json:"id"`
	SpaceId    string                 `json:"space_id"`
	DocumentId *string                `json:"document_id,omitempty"`
	Name       string                 `json:"name"`
	Icon       string                 `json:"icon,omitempty"`
	Elements   []interface{}          `json:"elements"`
	AppState   map[string]interface{} `json:"app_state"`
	Files      map[string]interface{} `json:"files"`
	Thumbnail  string                 `json:"thumbnail,omitempty"`
	CreatedBy  string                 `json:"created_by"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}
