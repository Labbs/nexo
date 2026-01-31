package dtos

import "time"

// Request DTOs

type ListVersionsRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	DocumentId string `path:"document_id" validate:"required"`
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
}

type GetVersionRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	DocumentId string `path:"document_id" validate:"required"`
	VersionId  string `path:"version_id" validate:"required,uuid4"`
}

type RestoreVersionRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	DocumentId string `path:"document_id" validate:"required"`
	VersionId  string `path:"version_id" validate:"required,uuid4"`
}

type CreateVersionRequest struct {
	SpaceId     string `path:"space_id" validate:"required,uuid4"`
	DocumentId  string `path:"document_id" validate:"required"`
	Description string `json:"description"`
}

// Response DTOs

type VersionItem struct {
	Id          string    `json:"id"`
	Version     int       `json:"version"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	UserId      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListVersionsResponse struct {
	Versions   []VersionItem `json:"versions"`
	TotalCount int64         `json:"total_count"`
}

type VersionConfig struct {
	FullWidth        bool   `json:"full_width"`
	Icon             string `json:"icon,omitempty"`
	Lock             bool   `json:"lock"`
	HeaderBackground string `json:"header_background,omitempty"`
}

type GetVersionResponse struct {
	Id          string        `json:"id"`
	Version     int           `json:"version"`
	DocumentId  string        `json:"document_id"`
	Name        string        `json:"name"`
	Content     []Block       `json:"content"`
	Config      VersionConfig `json:"config"`
	Description string        `json:"description,omitempty"`
	UserId      string        `json:"user_id"`
	UserName    string        `json:"user_name"`
	CreatedAt   time.Time     `json:"created_at"`
}

type CreateVersionResponse struct {
	VersionId string `json:"version_id"`
	Version   int    `json:"version"`
}
