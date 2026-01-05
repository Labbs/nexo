package persistence

import (
	"time"

	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type apiKeyPers struct {
	db *gorm.DB
}

func NewApiKeyPers(db *gorm.DB) *apiKeyPers {
	return &apiKeyPers{db: db}
}

func (p *apiKeyPers) Create(apiKey *domain.ApiKey) error {
	return p.db.Create(apiKey).Error
}

func (p *apiKeyPers) GetById(id string) (*domain.ApiKey, error) {
	var apiKey domain.ApiKey
	err := p.db.Where("id = ?", id).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (p *apiKeyPers) GetByKeyHash(keyHash string) (*domain.ApiKey, error) {
	var apiKey domain.ApiKey
	err := p.db.
		Preload("User").
		Where("key_hash = ? AND deleted_at IS NULL", keyHash).
		First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (p *apiKeyPers) GetByUserId(userId string) ([]domain.ApiKey, error) {
	var apiKeys []domain.ApiKey
	err := p.db.
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Order("created_at DESC").
		Find(&apiKeys).Error
	if err != nil {
		return nil, err
	}
	return apiKeys, nil
}

func (p *apiKeyPers) Update(apiKey *domain.ApiKey) error {
	return p.db.Save(apiKey).Error
}

func (p *apiKeyPers) Delete(id string) error {
	return p.db.Where("id = ?", id).Delete(&domain.ApiKey{}).Error
}

func (p *apiKeyPers) UpdateLastUsed(id string) error {
	now := time.Now()
	return p.db.Model(&domain.ApiKey{}).Where("id = ?", id).Update("last_used_at", &now).Error
}

// Admin methods

func (p *apiKeyPers) GetAll(limit, offset int) ([]domain.ApiKey, int64, error) {
	var apiKeys []domain.ApiKey
	var total int64

	// Count total (excluding soft deleted)
	if err := p.db.Model(&domain.ApiKey{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated list with user preloaded
	query := p.db.Model(&domain.ApiKey{}).
		Preload("User").
		Where("deleted_at IS NULL").
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err := query.Find(&apiKeys).Error; err != nil {
		return nil, 0, err
	}

	return apiKeys, total, nil
}
