package persistence

import (
	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type spacePers struct {
	db *gorm.DB
}

func NewSpacePers(db *gorm.DB) *spacePers {
	return &spacePers{db: db}
}

func (s *spacePers) Create(space *domain.Space) error {
	return s.db.Create(space).Error
}

func (s *spacePers) GetSpacesForUser(userId string) ([]domain.Space, error) {
	var spaces []domain.Space

	err := s.db.Preload("Owner").
		Preload("Permissions", "type = ?", domain.PermissionTypeSpace).
		Preload("Permissions.User").
		Preload("Permissions.Group").
		Where(
			s.db.Where("owner_id = ?", userId).
				Or("id IN (?)",
					s.db.Table("permission").
						Select("space_id").
						Where("type = ? AND user_id = ? AND deleted_at IS NULL", domain.PermissionTypeSpace, userId),
				).
				Or("type = ?", domain.SpaceTypePublic),
		).
		Find(&spaces).Error

	return spaces, err
}

func (s *spacePers) GetSpaceById(spaceId string) (*domain.Space, error) {
	var space domain.Space

	err := s.db.Preload("Owner").
		Preload("Permissions", "type = ?", domain.PermissionTypeSpace).
		Preload("Permissions.User").
		Preload("Permissions.Group").
		First(&space, "id = ?", spaceId).Error

	if err != nil {
		return nil, err
	}

	return &space, nil
}

func (s *spacePers) Update(space *domain.Space) error {
	return s.db.Save(space).Error
}

func (s *spacePers) Delete(spaceId string) error {
	return s.db.Where("id = ?", spaceId).Delete(&domain.Space{}).Error
}

// Admin methods

func (s *spacePers) GetAll(limit, offset int) ([]domain.Space, int64, error) {
	var spaces []domain.Space
	var total int64

	// Count total
	if err := s.db.Model(&domain.Space{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated list with owner preloaded
	query := s.db.Model(&domain.Space{}).Preload("Owner").Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err := query.Find(&spaces).Error; err != nil {
		return nil, 0, err
	}

	return spaces, total, nil
}
