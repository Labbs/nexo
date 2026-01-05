package persistence

import (
	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type favoritePers struct {
	db *gorm.DB
}

func NewFavoritePers(db *gorm.DB) *favoritePers {
	return &favoritePers{db: db}
}

func (p *favoritePers) GetLatestFavoritePositionByUser(userId string) (int, error) {
	var position int
	err := p.db.Model(&domain.Favorite{}).Select("position").Where("user_id = ?", userId).Order("position desc").Limit(1).Scan(&position).Error
	if err != nil {
		return 0, err
	}
	return position, nil
}

func (p *favoritePers) Create(favorite *domain.Favorite) error {
	return p.db.Create(favorite).Error
}

func (p *favoritePers) Delete(documentId, userId string, spaceId string) error {
	return p.db.Where("document_id = ? AND user_id = ? AND space_id = ?", documentId, userId, spaceId).Delete(&domain.Favorite{}).Error
}

func (f *favoritePers) GetMyFavoritesWithMainDocumentInformations(userId string) ([]domain.Favorite, error) {
	var favorites []domain.Favorite
	err := f.db.
		Preload("Document", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, slug, config")
		}).
		Where("user_id = ?", userId).
		Order("position asc").
		Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	return favorites, nil
}

func (f *favoritePers) UpdateFavoritePosition(favorite *domain.Favorite) error {
	return f.db.Save(favorite).Error
}

func (f *favoritePers) GetFavoriteByIdAndUserId(favoriteId, userId string) (*domain.Favorite, error) {
	var favorite domain.Favorite
	err := f.db.Where("id = ? AND user_id = ?", favoriteId, userId).First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}
