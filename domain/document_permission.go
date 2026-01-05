package domain

import (
	"time"

	"gorm.io/gorm"
)

type DocumentPermission struct {
	Id         string
	DocumentId string
	Document   Document `gorm:"foreignKey:DocumentId;references:Id"`

	UserId *string
	User   *User `gorm:"foreignKey:UserId;references:Id"`

	GroupId *string
	Group   *Group `gorm:"foreignKey:GroupId;references:Id"`

	Role DocumentRole

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (dp *DocumentPermission) TableName() string {
	return "document_permission"
}

type DocumentRole string

const (
	DocumentRoleOwner  DocumentRole = "owner"  // Full control, can manage permissions
	DocumentRoleEditor DocumentRole = "editor"
	DocumentRoleViewer DocumentRole = "viewer"
	DocumentRoleDenied DocumentRole = "denied" // Explicitly deny access
)

type DocumentPermissionPers interface {
	ListByDocumentId(documentId string) ([]DocumentPermission, error)
	GetByDocumentAndUser(documentId, userId string) (*DocumentPermission, error)
	Upsert(documentId, userId string, role DocumentRole) error
	Delete(documentId, userId string) error
}
