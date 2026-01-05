package domain

import (
	"time"

	"gorm.io/gorm"
)

// DatabaseType defines the types of databases
type DatabaseType string

const (
	DatabaseTypeSpreadsheet DatabaseType = "spreadsheet"
	DatabaseTypeDocument    DatabaseType = "document"
)

// Database represents a Notion-like database (table)
type Database struct {
	Id string

	SpaceId string
	Space   Space `gorm:"foreignKey:SpaceId;references:Id"`

	// Optional: database can be inline in a document
	DocumentId *string
	Document   *Document `gorm:"foreignKey:DocumentId;references:Id"`

	Name        string
	Description string
	Icon        string

	// Schema defines the columns/properties of the database
	Schema JSONBArray // [{id, name, type, options}]

	// Views configuration (table, board, calendar, etc.)
	Views JSONBArray // [{id, name, type, filter, sort, columns}]

	// Default view type
	DefaultView string

	// Type of database: "spreadsheet" or "document"
	Type DatabaseType

	CreatedBy string
	User      User `gorm:"foreignKey:CreatedBy;references:Id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (d *Database) TableName() string {
	return "database"
}

// PropertyType defines the types of properties/columns
type PropertyType string

const (
	PropertyTypeTitle       PropertyType = "title"
	PropertyTypeText        PropertyType = "text"
	PropertyTypeNumber      PropertyType = "number"
	PropertyTypeSelect      PropertyType = "select"
	PropertyTypeMultiSelect PropertyType = "multi_select"
	PropertyTypeDate        PropertyType = "date"
	PropertyTypeCheckbox    PropertyType = "checkbox"
	PropertyTypeUrl         PropertyType = "url"
	PropertyTypeEmail       PropertyType = "email"
	PropertyTypePhone       PropertyType = "phone"
	PropertyTypeRelation    PropertyType = "relation"
	PropertyTypeRollup      PropertyType = "rollup"
	PropertyTypeFormula     PropertyType = "formula"
	PropertyTypeCreatedTime PropertyType = "created_time"
	PropertyTypeUpdatedTime PropertyType = "updated_time"
	PropertyTypeCreatedBy   PropertyType = "created_by"
	PropertyTypeUpdatedBy   PropertyType = "updated_by"
	PropertyTypeFiles       PropertyType = "files"
	PropertyTypePerson      PropertyType = "person"
)

// ViewType defines the types of database views
type ViewType string

const (
	ViewTypeTable    ViewType = "table"
	ViewTypeBoard    ViewType = "board"
	ViewTypeCalendar ViewType = "calendar"
	ViewTypeGallery  ViewType = "gallery"
	ViewTypeList     ViewType = "list"
	ViewTypeTimeline ViewType = "timeline"
)

type DatabasePers interface {
	Create(database *Database) error
	GetById(id string) (*Database, error)
	GetBySpaceId(spaceId string) ([]Database, error)
	GetByDocumentId(documentId string) ([]Database, error)
	Update(database *Database) error
	Delete(id string) error
}

// DatabaseRow represents a row/page in a database
type DatabaseRow struct {
	Id string

	DatabaseId string
	Database   Database `gorm:"foreignKey:DatabaseId;references:Id"`

	// Properties holds the values for each column
	Properties JSONB // {propertyId: value}

	// Row can optionally have page content
	Content JSONB

	// ShowInSidebar indicates if this row should appear in sidebar (for document databases)
	ShowInSidebar bool

	CreatedBy string
	User      User `gorm:"foreignKey:CreatedBy;references:Id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (r *DatabaseRow) TableName() string {
	return "database_row"
}

// FilterRule defines a single filter condition
type FilterRule struct {
	Property  string      `json:"property"`
	Condition string      `json:"condition"` // eq, neq, gt, lt, gte, lte, contains, is_empty, is_not_empty
	Value     interface{} `json:"value,omitempty"`
}

// FilterConfig defines the filter configuration with AND/OR groups
type FilterConfig struct {
	And []FilterRule `json:"and,omitempty"`
	Or  []FilterRule `json:"or,omitempty"`
}

// SortRule defines a sort condition
type SortRule struct {
	PropertyId string `json:"property_id"`
	Direction  string `json:"direction"` // "asc" or "desc"
}

// RowQueryOptions contains filter and sort options for row queries
type RowQueryOptions struct {
	Filter *FilterConfig
	Sort   []SortRule
	Limit  int
	Offset int
}

type DatabaseRowPers interface {
	Create(row *DatabaseRow) error
	GetById(id string) (*DatabaseRow, error)
	GetByDatabaseId(databaseId string, limit, offset int) ([]DatabaseRow, error)
	GetByDatabaseIdWithOptions(databaseId string, options RowQueryOptions) ([]DatabaseRow, error)
	GetRowCount(databaseId string) (int64, error)
	GetRowCountWithFilter(databaseId string, filter *FilterConfig) (int64, error)
	Update(row *DatabaseRow) error
	Delete(id string) error
	BulkDelete(ids []string) error
}
