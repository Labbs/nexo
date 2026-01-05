package persistence

import (
	"github.com/gofiber/fiber/v2/utils"
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
		Preload("Permissions").
		Where(
			s.db.Where("owner_id = ?", userId).
				Or("id IN (?)",
					s.db.Table("space_permission").
						Select("space_id").
						Where("user_id = ? AND deleted_at IS NULL", userId),
				).
				Or("type = ?", domain.SpaceTypePublic),
		).
		Find(&spaces).Error

	return spaces, err
}

func (s *spacePers) GetSpaceById(spaceId string) (*domain.Space, error) {
	var space domain.Space

	err := s.db.Preload("Owner").
		Preload("Permissions").
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

func (s *spacePers) ListPermissions(spaceId string) ([]domain.SpacePermission, error) {
    var perms []domain.SpacePermission
    err := s.db.Where("space_id = ? AND deleted_at IS NULL", spaceId).Find(&perms).Error
    return perms, err
}

func (s *spacePers) UpsertUserPermission(spaceId string, userId string, role domain.SpaceRole) error {
    var perm domain.SpacePermission
    err := s.db.Where("space_id = ? AND user_id = ?", spaceId, userId).First(&perm).Error
    if err != nil {
        // create if not exists
        perm = domain.SpacePermission{
            Id:      utils.UUIDv4(),
            SpaceId: spaceId,
            UserId:  &userId,
            Role:    role,
        }
        return s.db.Create(&perm).Error
    }
    perm.Role = role
    return s.db.Save(&perm).Error
}

func (s *spacePers) DeleteUserPermission(spaceId string, userId string) error {
	return s.db.Where("space_id = ? AND user_id = ?", spaceId, userId).Delete(&domain.SpacePermission{}).Error
}

// Group permissions

func (s *spacePers) UpsertGroupPermission(spaceId string, groupId string, role domain.SpaceRole) error {
	var perm domain.SpacePermission
	err := s.db.Where("space_id = ? AND group_id = ?", spaceId, groupId).First(&perm).Error
	if err != nil {
		// create if not exists
		perm = domain.SpacePermission{
			Id:      utils.UUIDv4(),
			SpaceId: spaceId,
			GroupId: &groupId,
			Role:    role,
		}
		return s.db.Create(&perm).Error
	}
	perm.Role = role
	return s.db.Save(&perm).Error
}

func (s *spacePers) DeleteGroupPermission(spaceId string, groupId string) error {
	return s.db.Where("space_id = ? AND group_id = ?", spaceId, groupId).Delete(&domain.SpacePermission{}).Error
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

func (s *spacePers) ListPermissionsWithDetails(spaceId string) ([]domain.SpacePermission, error) {
	var perms []domain.SpacePermission
	err := s.db.Preload("User").Preload("Group").
		Where("space_id = ? AND deleted_at IS NULL", spaceId).
		Find(&perms).Error
	return perms, err
}
