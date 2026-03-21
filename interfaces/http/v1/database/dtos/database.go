package dtos

import "time"

// Request DTOs

type EmptyRequest struct{}

type PropertySchema struct {
	Id      string         `json:"id"`
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	Options map[string]any `json:"options,omitempty"`
}

type ViewConfig struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Type          string         `json:"type"`
	Filter        map[string]any `json:"filter,omitempty"`
	Sort          []SortConfig   `json:"sort,omitempty"`
	Columns       []string       `json:"columns,omitempty"`
	HiddenColumns []string       `json:"hidden_columns,omitempty"`
	GroupBy       string         `json:"group_by,omitempty"`
}

type SortConfig struct {
	PropertyId string `json:"property_id"`
	Direction  string `json:"direction"`
}

type CreateDatabaseRequest struct {
	SpaceId     string           `json:"space_id" resource:"space" action:"write"`
	DocumentId  *string          `json:"document_id,omitempty"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Icon        string           `json:"icon,omitempty"`
	Schema      []PropertySchema `json:"schema"`
	Type        string           `json:"type,omitempty"` // "spreadsheet" or "document", defaults to "spreadsheet"
}

type ListDatabasesRequest struct {
	SpaceId string `query:"space_id" resource:"space" action:"read"`
}

type GetDatabaseRequest struct {
	DatabaseId string `path:"database_id" resource:"database" action:"read"`
}

type UpdateDatabaseRequest struct {
	DatabaseId  string            `path:"database_id" resource:"database" action:"write"`
	Name        *string           `json:"name,omitempty"`
	Description *string           `json:"description,omitempty"`
	Icon        *string           `json:"icon,omitempty"`
	Schema      *[]PropertySchema `json:"schema,omitempty"`
	DefaultView *string           `json:"default_view,omitempty"`
}

type DeleteDatabaseRequest struct {
	DatabaseId string `path:"database_id" resource:"database"`
}

// Row requests
type CreateRowRequest struct {
	DatabaseId    string         `path:"database_id" resource:"database" action:"write"`
	Properties    map[string]any `json:"properties"`
	Content       map[string]any `json:"content,omitempty"`
	ShowInSidebar bool           `json:"show_in_sidebar,omitempty"`
}

type ListRowsRequest struct {
	DatabaseId string `path:"database_id" resource:"database" action:"read"`
	ViewId     string `query:"view_id"`
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
}

type GetRowRequest struct {
	DatabaseId string `path:"database_id" resource:"database" action:"read"`
	RowId      string `path:"row_id"`
}

type UpdateRowRequest struct {
	DatabaseId    string         `path:"database_id" resource:"database" action:"write"`
	RowId         string         `path:"row_id"`
	Properties    map[string]any `json:"properties,omitempty"`
	Content       map[string]any `json:"content,omitempty"`
	ShowInSidebar *bool          `json:"show_in_sidebar,omitempty"`
}

type DeleteRowRequest struct {
	DatabaseId string `path:"database_id" resource:"database" action:"write"`
	RowId      string `path:"row_id"`
}

type BulkDeleteRowsRequest struct {
	DatabaseId string   `path:"database_id" resource:"database" action:"write"`
	RowIds     []string `json:"row_ids"`
}

// View requests
type CreateViewRequest struct {
	DatabaseId string         `path:"database_id" resource:"database" action:"write"`
	Name       string         `json:"name" validate:"required"`
	Type       string         `json:"type" validate:"required,oneof=table board calendar gallery list timeline"`
	Filter     map[string]any `json:"filter,omitempty"`
	Sort       []SortConfig   `json:"sort,omitempty"`
	Columns    []string       `json:"columns,omitempty"`
}

type UpdateViewRequest struct {
	DatabaseId    string         `path:"database_id" resource:"database" action:"write"`
	ViewId        string         `path:"view_id"`
	Name          *string        `json:"name,omitempty"`
	Type          *string        `json:"type,omitempty"`
	Filter        map[string]any `json:"filter,omitempty"`
	Sort          []SortConfig   `json:"sort,omitempty"`
	Columns       []string       `json:"columns,omitempty"`
	HiddenColumns []string       `json:"hidden_columns,omitempty"`
	GroupBy       *string        `json:"group_by,omitempty"`
}

type DeleteViewRequest struct {
	DatabaseId string `path:"database_id" resource:"database" action:"write"`
	ViewId     string `path:"view_id"`
}

// Move database
type MoveDatabaseRequest struct {
	DatabaseId string  `path:"database_id" resource:"database" action:"write"`
	DocumentId *string `json:"document_id"`
}

type MoveDatabaseResponse struct {
	Id         string  `json:"id"`
	DocumentId *string `json:"document_id,omitempty"`
}

// Filter rule for querying rows
type FilterRule struct {
	Property  string `json:"property"`
	Condition string `json:"condition"` // eq, neq, gt, lt, gte, lte, contains, is_empty, is_not_empty
	Value     any    `json:"value,omitempty"`
}

type FilterConfig struct {
	And []FilterRule `json:"and,omitempty"`
	Or  []FilterRule `json:"or,omitempty"`
}

// Response DTOs

type MessageResponse struct {
	Message string `json:"message"`
}

// UserInfo contains basic user information for display
type UserInfo struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}

type CreateDatabaseResponse struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Icon        string           `json:"icon"`
	Schema      []PropertySchema `json:"schema"`
	DefaultView string           `json:"default_view"`
	Type        string           `json:"type"`
	CreatedAt   time.Time        `json:"created_at"`
}

type DatabaseItem struct {
	Id          string    `json:"id"`
	DocumentId  *string   `json:"document_id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Type        string    `json:"type"`
	RowCount    int64     `json:"row_count"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListDatabasesResponse struct {
	Databases []DatabaseItem `json:"databases"`
}

type GetDatabaseResponse struct {
	Id          string           `json:"id"`
	SpaceId     string           `json:"space_id"`
	DocumentId  *string          `json:"document_id,omitempty"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Icon        string           `json:"icon"`
	Schema      []PropertySchema `json:"schema"`
	Views       []ViewConfig     `json:"views"`
	DefaultView string           `json:"default_view"`
	Type        string           `json:"type"`
	CreatedBy   string           `json:"created_by"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type CreateRowResponse struct {
	Id         string         `json:"id"`
	Properties map[string]any `json:"properties"`
	CreatedAt  time.Time      `json:"created_at"`
}

type RowItem struct {
	Id            string         `json:"id"`
	Properties    map[string]any `json:"properties"`
	Content       map[string]any `json:"content,omitempty"`
	ShowInSidebar bool           `json:"show_in_sidebar"`
	CreatedBy     string         `json:"created_by"`
	CreatedByUser *UserInfo      `json:"created_by_user,omitempty"`
	UpdatedBy     string         `json:"updated_by,omitempty"`
	UpdatedByUser *UserInfo      `json:"updated_by_user,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type ListRowsResponse struct {
	Rows       []RowItem `json:"rows"`
	TotalCount int64     `json:"total_count"`
}

type GetRowResponse struct {
	Id            string         `json:"id"`
	DatabaseId    string         `json:"database_id"`
	Properties    map[string]any `json:"properties"`
	Content       map[string]any `json:"content,omitempty"`
	ShowInSidebar bool           `json:"show_in_sidebar"`
	CreatedBy     string         `json:"created_by"`
	CreatedByUser *UserInfo      `json:"created_by_user,omitempty"`
	UpdatedBy     string         `json:"updated_by,omitempty"`
	UpdatedByUser *UserInfo      `json:"updated_by_user,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// View responses
type CreateViewResponse struct {
	Id      string         `json:"id"`
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	Filter  map[string]any `json:"filter,omitempty"`
	Sort    []SortConfig   `json:"sort,omitempty"`
	Columns []string       `json:"columns,omitempty"`
}

// Available property types
type AvailableTypesResponse struct {
	Types []TypeInfo `json:"types"`
}

type TypeInfo struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Search
type SearchDatabasesRequest struct {
	Query   string  `query:"q"`
	SpaceId *string `query:"space_id"`
	Limit   int     `query:"limit"`
}

type SearchDatabaseResultItem struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon,omitempty"`
	Type        string    `json:"type"`
	SpaceId     string    `json:"space_id"`
	SpaceName   string    `json:"space_name"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchDatabasesResponse struct {
	Results []SearchDatabaseResultItem `json:"results"`
}
