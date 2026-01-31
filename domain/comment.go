package domain

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	Id string

	DocumentId string
	Document   Document `gorm:"foreignKey:DocumentId;references:Id"`

	UserId string
	User   User `gorm:"foreignKey:UserId;references:Id"`

	// For replies - optional parent comment
	ParentId *string
	Parent   *Comment `gorm:"foreignKey:ParentId;references:Id"`

	Content string

	// For inline comments - optional block reference
	BlockId *string

	Resolved  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (c *Comment) TableName() string {
	return "comment"
}

type CommentPers interface {
	Create(comment *Comment) error
	GetById(commentId string) (*Comment, error)
	GetByDocumentId(documentId string) ([]Comment, error)
	Update(comment *Comment) error
	Delete(commentId string) error
	Resolve(commentId string, resolved bool) error
}
