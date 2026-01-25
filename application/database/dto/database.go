package dto

import "time"

// UserInfo contains basic user information
type UserInfo struct {
	Id        string
	Username  string
	AvatarUrl string
}

// Property schema
type PropertySchema struct {
	Id      string                 `json:"id"`
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// View configuration
type ViewConfig struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	Filter        map[string]interface{} `json:"filter,omitempty"`
	Sort          []SortConfig           `json:"sort,omitempty"`
	Columns       []string               `json:"columns,omitempty"`
	HiddenColumns []string               `json:"hidden_columns,omitempty"`
}

type SortConfig struct {
	PropertyId string `json:"property_id"`
	Direction  string `json:"direction"` // "asc" or "desc"
}

// Create database
type CreateDatabaseInput struct {
	UserId      string
	SpaceId     string
	DocumentId  *string
	Name        string
	Description string
	Icon        string
	Schema      []PropertySchema
	Type        string // "spreadsheet" or "document", defaults to "spreadsheet"
}

type CreateDatabaseOutput struct {
	Id          string
	Name        string
	Description string
	Icon        string
	Schema      []PropertySchema
	DefaultView string
	Type        string
	CreatedAt   time.Time
}

// List databases
type ListDatabasesInput struct {
	UserId  string
	SpaceId string
}

type DatabaseItem struct {
	Id          string
	DocumentId  *string
	Name        string
	Description string
	Icon        string
	Type        string
	RowCount    int64
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListDatabasesOutput struct {
	Databases []DatabaseItem
}

// Get database
type GetDatabaseInput struct {
	UserId     string
	DatabaseId string
}

type GetDatabaseOutput struct {
	Id          string
	SpaceId     string
	DocumentId  *string
	Name        string
	Description string
	Icon        string
	Schema      []PropertySchema
	Views       []ViewConfig
	DefaultView string
	Type        string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Update database
type UpdateDatabaseInput struct {
	UserId      string
	DatabaseId  string
	Name        *string
	Description *string
	Icon        *string
	Schema      *[]PropertySchema
	DefaultView *string
}

// Delete database
type DeleteDatabaseInput struct {
	UserId     string
	DatabaseId string
}

// Create view
type CreateViewInput struct {
	UserId     string
	DatabaseId string
	Name       string
	Type       string
	Filter     map[string]interface{}
	Sort       []SortConfig
	Columns    []string
}

type CreateViewOutput struct {
	Id      string
	Name    string
	Type    string
	Filter  map[string]interface{}
	Sort    []SortConfig
	Columns []string
}

// Update view
type UpdateViewInput struct {
	UserId        string
	DatabaseId    string
	ViewId        string
	Name          *string
	Filter        map[string]interface{}
	Sort          []SortConfig
	Columns       []string
	HiddenColumns []string
}

// Delete view
type DeleteViewInput struct {
	UserId     string
	DatabaseId string
	ViewId     string
}

// Row operations
type CreateRowInput struct {
	UserId        string
	DatabaseId    string
	Properties    map[string]interface{}
	Content       map[string]interface{}
	ShowInSidebar bool
}

type CreateRowOutput struct {
	Id         string
	Properties map[string]interface{}
	CreatedAt  time.Time
}

type ListRowsInput struct {
	UserId     string
	DatabaseId string
	ViewId     string
	Limit      int
	Offset     int
}

type RowItem struct {
	Id            string
	Properties    map[string]interface{}
	ShowInSidebar bool
	CreatedBy     string
	CreatedByUser *UserInfo
	UpdatedBy     string
	UpdatedByUser *UserInfo
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ListRowsOutput struct {
	Rows       []RowItem
	TotalCount int64
}

type GetRowInput struct {
	UserId     string
	DatabaseId string
	RowId      string
}

type GetRowOutput struct {
	Id            string
	DatabaseId    string
	Properties    map[string]interface{}
	Content       map[string]interface{}
	ShowInSidebar bool
	CreatedBy     string
	CreatedByUser *UserInfo
	UpdatedBy     string
	UpdatedByUser *UserInfo
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UpdateRowInput struct {
	UserId        string
	DatabaseId    string
	RowId         string
	Properties    map[string]interface{}
	Content       map[string]interface{}
	ShowInSidebar *bool
}

type DeleteRowInput struct {
	UserId     string
	DatabaseId string
	RowId      string
}

type BulkDeleteRowsInput struct {
	UserId     string
	DatabaseId string
	RowIds     []string
}
