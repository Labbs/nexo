package domain

import (
	"time"

	"gorm.io/gorm"
)

type SpacePermission struct {
	Id string

	SpaceId string
	Space   Space `gorm:"foreignKey:SpaceId;references:Id"`

	UserId *string
	User   *User `gorm:"foreignKey:UserId;references:Id"`

	GroupId *string
	Group   *Group `gorm:"foreignKey:GroupId;references:Id"`

	Role SpaceRole

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type SpaceRole string

const (
	SpaceRoleOwner  SpaceRole = "owner"
	SpaceRoleAdmin  SpaceRole = "admin"
	SpaceRoleEditor SpaceRole = "editor"
	SpaceRoleViewer SpaceRole = "viewer"
)

func (sp *SpacePermission) TableName() string {
	return "space_permission"
}
