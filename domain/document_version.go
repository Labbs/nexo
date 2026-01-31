package domain

import (
	"time"

	"gorm.io/datatypes"
)

type DocumentVersion struct {
	Id string

	DocumentId string
	Document   Document `gorm:"foreignKey:DocumentId;references:Id"`

	// User who created this version
	UserId string
	User   User `gorm:"foreignKey:UserId;references:Id"`

	// Version number (auto-incremented per document)
	Version int

	// Snapshot of document at this version
	Name    string
	Content datatypes.JSON
	Config  DocumentConfig

	// Optional description of changes
	Description string

	CreatedAt time.Time
}

func (v *DocumentVersion) TableName() string {
	return "document_version"
}

type DocumentVersionPers interface {
	Create(version *DocumentVersion) error
	GetByDocumentId(documentId string, limit int, offset int) ([]DocumentVersion, error)
	GetById(versionId string) (*DocumentVersion, error)
	GetLatestVersion(documentId string) (*DocumentVersion, error)
	GetVersionCount(documentId string) (int64, error)
	DeleteOldVersions(documentId string, keepCount int) error
}
