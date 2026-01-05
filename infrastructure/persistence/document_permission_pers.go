package persistence

import (
	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type documentPermissionPers struct {
	db *gorm.DB
}

func NewDocumentPermissionPers(db *gorm.DB) *documentPermissionPers {
	return &documentPermissionPers{db: db}
}

func (p *documentPermissionPers) ListByDocumentId(documentId string) ([]domain.DocumentPermission, error) {
	var perms []domain.DocumentPermission
	err := p.db.Preload("User").
		Where("document_id = ? AND deleted_at IS NULL", documentId).
		Find(&perms).Error
	return perms, err
}

func (p *documentPermissionPers) GetByDocumentAndUser(documentId, userId string) (*domain.DocumentPermission, error) {
	var perm domain.DocumentPermission
	err := p.db.Where("document_id = ? AND user_id = ? AND deleted_at IS NULL", documentId, userId).
		First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (p *documentPermissionPers) Upsert(documentId, userId string, role domain.DocumentRole) error {
	var perm domain.DocumentPermission
	err := p.db.Where("document_id = ? AND user_id = ?", documentId, userId).First(&perm).Error
	if err != nil {
		// Create if not exists
		perm = domain.DocumentPermission{
			Id:         utils.UUIDv4(),
			DocumentId: documentId,
			UserId:     &userId,
			Role:       role,
		}
		return p.db.Create(&perm).Error
	}
	// Update existing
	perm.Role = role
	perm.DeletedAt = gorm.DeletedAt{} // Restore if soft-deleted
	return p.db.Save(&perm).Error
}

func (p *documentPermissionPers) Delete(documentId, userId string) error {
	return p.db.Where("document_id = ? AND user_id = ?", documentId, userId).
		Delete(&domain.DocumentPermission{}).Error
}
