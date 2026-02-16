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
	GroupBy       string                 `json:"group_by,omitempty"`
}

type SortConfig struct {
	PropertyId string `json:"property_id"`
	Direction  string `json:"direction"` // "asc" or "desc"
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

type RowItem struct {
	Id            string
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

// Permission item for listing
type DatabasePermissionItem struct {
	Id        string  `json:"id"`
	UserId    *string `json:"user_id,omitempty"`
	Username  *string `json:"username,omitempty"`
	GroupId   *string `json:"group_id,omitempty"`
	GroupName *string `json:"group_name,omitempty"`
	Role      string  `json:"role"`
}
