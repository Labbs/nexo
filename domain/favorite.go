package domain

import (
	"time"
)

type Favorite struct {
	Id string

	UserId string
	User   User `gorm:"foreignKey:UserId;references:Id"`

	DocumentId string
	Document   Document `gorm:"foreignKey:DocumentId;references:Id"`

	SpaceId string
	Space   Space `gorm:"foreignKey:SpaceId;references:Id"`

	Position int

	CreatedAt time.Time
}

func (f *Favorite) TableName() string {
	return "favorite"
}

type FavoritePers interface {
	GetLatestFavoritePositionByUser(userId string) (int, error)
	Create(favorite *Favorite) error
	Delete(documentId, userId string, spaceId string) error
	GetMyFavoritesWithMainDocumentInformations(userId string) ([]Favorite, error)
	UpdateFavoritePosition(favorite *Favorite) error
	GetFavoriteByIdAndUserId(favoriteId, userId string) (*Favorite, error)
}
