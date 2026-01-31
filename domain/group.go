package domain

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	Id          string
	Name        string
	Description string

	Role Role `gorm:"type:role;default:'user'"`

	// Owner is the username of the user who owns the group
	OwnerId string
	// OwnerUser is the user who owns the group
	OwnerUser User `gorm:"foreignKey:OwnerId;references:Id"`

	// Members is the list of users who are members of the group
	Members []User `gorm:"many2many:group_members;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (g *Group) TableName() string {
	return "group"
}

// GroupPers defines the persistence interface for groups
type GroupPers interface {
	Create(group *Group) error
	GetById(groupId string) (*Group, error)
	GetAll(limit, offset int) ([]Group, int64, error)
	Update(group *Group) error
	Delete(groupId string) error

	// Member management
	AddMember(groupId, userId string) error
	RemoveMember(groupId, userId string) error
	GetMembers(groupId string) ([]User, error)
}
