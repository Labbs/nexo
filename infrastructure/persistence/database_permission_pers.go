package persistence

import (
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type databasePermissionPers struct {
	db *gorm.DB
}

func NewDatabasePermissionPers(db *gorm.DB) *databasePermissionPers {
	return &databasePermissionPers{db: db}
}

func (p *databasePermissionPers) ListByDatabaseId(databaseId string) ([]domain.DatabasePermission, error) {
	var perms []domain.DatabasePermission
	err := p.db.Preload("User").Preload("Group").
		Where("database_id = ? AND deleted_at IS NULL", databaseId).
		Find(&perms).Error
	return perms, err
}

func (p *databasePermissionPers) GetByDatabaseAndUser(databaseId, userId string) (*domain.DatabasePermission, error) {
	var perm domain.DatabasePermission
	err := p.db.Where("database_id = ? AND user_id = ? AND deleted_at IS NULL", databaseId, userId).
		First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (p *databasePermissionPers) GetByDatabaseAndGroup(databaseId, groupId string) (*domain.DatabasePermission, error) {
	var perm domain.DatabasePermission
	err := p.db.Where("database_id = ? AND group_id = ? AND deleted_at IS NULL", databaseId, groupId).
		First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (p *databasePermissionPers) Upsert(perm *domain.DatabasePermission) error {
	var existing domain.DatabasePermission
	var err error

	if perm.UserId != nil {
		err = p.db.Where("database_id = ? AND user_id = ?", perm.DatabaseId, *perm.UserId).First(&existing).Error
	} else if perm.GroupId != nil {
		err = p.db.Where("database_id = ? AND group_id = ?", perm.DatabaseId, *perm.GroupId).First(&existing).Error
	} else {
		return gorm.ErrRecordNotFound
	}

	now := time.Now()
	if err != nil {
		// Create if not exists
		perm.Id = utils.UUIDv4()
		perm.CreatedAt = now
		perm.UpdatedAt = now
		return p.db.Create(perm).Error
	}

	// Update existing
	existing.Role = perm.Role
	existing.UpdatedAt = now
	existing.DeletedAt = gorm.DeletedAt{} // Restore if soft-deleted
	return p.db.Save(&existing).Error
}

func (p *databasePermissionPers) Delete(id string) error {
	return p.db.Where("id = ?", id).Delete(&domain.DatabasePermission{}).Error
}

func (p *databasePermissionPers) DeleteByDatabaseAndUser(databaseId, userId string) error {
	return p.db.Where("database_id = ? AND user_id = ?", databaseId, userId).
		Delete(&domain.DatabasePermission{}).Error
}

func (p *databasePermissionPers) DeleteByDatabaseAndGroup(databaseId, groupId string) error {
	return p.db.Where("database_id = ? AND group_id = ?", databaseId, groupId).
		Delete(&domain.DatabasePermission{}).Error
}
