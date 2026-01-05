package domain

import (
	"time"
)

type User struct {
	Id       string
	Username string
	Email    string
	Password string

	AvatarUrl   string
	Preferences JSONB
	Active      bool

	Role Role `gorm:"type:role;default:'user'"`

	Favorites []Favorite `gorm:"foreignKey:UserId;references:Id"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) TableName() string {
	return "user"
}

type UserPers interface {
	GetByUsername(username string) (User, error)
	GetByEmail(email string) (User, error)
	GetById(id string) (User, error)
	Create(user User) (User, error)
	Update(user *User) error
	UpdatePassword(userId, hashedPassword string) error
	// Admin methods
	GetAll(limit, offset int) ([]User, int64, error)
	UpdateRole(userId string, role Role) error
	UpdateActive(userId string, active bool) error
	Delete(userId string) error
}
