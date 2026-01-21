package domain

import (
	"time"

	"gorm.io/gorm"
)

// Drawing represents an Excalidraw drawing
type Drawing struct {
	Id string

	SpaceId string
	Space   Space `gorm:"foreignKey:SpaceId;references:Id"`

	// Optional: drawing can be inline in a document
	DocumentId *string
	Document   *Document `gorm:"foreignKey:DocumentId;references:Id"`

	Name string
	Icon string // Emoji icon for the drawing

	// Excalidraw data
	Elements JSONBArray // Excalidraw elements array
	AppState JSONB      // Excalidraw appState
	Files    JSONB      // Embedded files (images in base64)

	// Base64 PNG thumbnail for preview
	Thumbnail string

	CreatedBy string
	User      User `gorm:"foreignKey:CreatedBy;references:Id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (d *Drawing) TableName() string {
	return "drawing"
}

type DrawingPers interface {
	Create(drawing *Drawing) error
	GetById(id string) (*Drawing, error)
	GetBySpaceId(spaceId string) ([]Drawing, error)
	GetByDocumentId(documentId string) ([]Drawing, error)
	Update(drawing *Drawing) error
	Delete(id string) error
}
