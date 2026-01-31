package domain

import (
	"time"

	"gorm.io/gorm"
)

// PermissionType defines the type of resource the permission applies to
type PermissionType string

const (
	PermissionTypeSpace    PermissionType = "space"
	PermissionTypeDocument PermissionType = "document"
	PermissionTypeDatabase PermissionType = "database"
	PermissionTypeDrawing  PermissionType = "drawing"
)

// PermissionRole defines the role/access level
type PermissionRole string

const (
	PermissionRoleOwner  PermissionRole = "owner"  // Full control, can manage permissions
	PermissionRoleAdmin  PermissionRole = "admin"  // Admin access (for spaces)
	PermissionRoleEditor PermissionRole = "editor" // Can edit
	PermissionRoleViewer PermissionRole = "viewer" // Read-only
	PermissionRoleDenied PermissionRole = "denied" // Explicitly deny access
)

// Permission represents a unified permission entry for any resource type
type Permission struct {
	Id   string
	Type PermissionType // "space", "document", "database", "drawing"

	// Resource IDs - only one should be set based on Type
	SpaceId    *string
	Space      *Space `gorm:"foreignKey:SpaceId;references:Id"`
	DocumentId *string
	Document   *Document `gorm:"foreignKey:DocumentId;references:Id"`
	DatabaseId *string
	Database   *Database `gorm:"foreignKey:DatabaseId;references:Id"`
	DrawingId  *string
	Drawing    *Drawing `gorm:"foreignKey:DrawingId;references:Id"`

	// Target - either a user or a group
	UserId *string
	User   *User `gorm:"foreignKey:UserId;references:Id"`

	GroupId *string
	Group   *Group `gorm:"foreignKey:GroupId;references:Id"`

	Role PermissionRole

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (p *Permission) TableName() string {
	return "permission"
}

// PermissionPers is the persistence interface for permissions
type PermissionPers interface {
	// Generic methods
	ListByResource(resourceType PermissionType, resourceId string) ([]Permission, error)
	GetByResourceAndUser(resourceType PermissionType, resourceId, userId string) (*Permission, error)
	GetByResourceAndGroup(resourceType PermissionType, resourceId, groupId string) (*Permission, error)
	UpsertUser(resourceType PermissionType, resourceId, userId string, role PermissionRole) error
	UpsertGroup(resourceType PermissionType, resourceId, groupId string, role PermissionRole) error
	DeleteUser(resourceType PermissionType, resourceId, userId string) error
	DeleteGroup(resourceType PermissionType, resourceId, groupId string) error
}

