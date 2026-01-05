package domain

import (
	"time"

	"gorm.io/gorm"
)

type DatabasePermission struct {
	Id         string
	DatabaseId string
	Database   Database `gorm:"foreignKey:DatabaseId;references:Id"`

	UserId *string
	User   *User `gorm:"foreignKey:UserId;references:Id"`

	GroupId *string
	Group   *Group `gorm:"foreignKey:GroupId;references:Id"`

	Role DatabaseRole

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (dp *DatabasePermission) TableName() string {
	return "database_permission"
}

type DatabaseRole string

const (
	DatabaseRoleEditor DatabaseRole = "editor"
	DatabaseRoleViewer DatabaseRole = "viewer"
	DatabaseRoleDenied DatabaseRole = "denied" // Explicitly deny access
)

type DatabasePermissionPers interface {
	ListByDatabaseId(databaseId string) ([]DatabasePermission, error)
	GetByDatabaseAndUser(databaseId, userId string) (*DatabasePermission, error)
	GetByDatabaseAndGroup(databaseId, groupId string) (*DatabasePermission, error)
	Upsert(perm *DatabasePermission) error
	Delete(id string) error
	DeleteByDatabaseAndUser(databaseId, userId string) error
	DeleteByDatabaseAndGroup(databaseId, groupId string) error
}
