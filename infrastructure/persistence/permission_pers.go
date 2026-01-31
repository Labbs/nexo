package persistence

import (
	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type permissionPers struct {
	db *gorm.DB
}

func NewPermissionPers(db *gorm.DB) *permissionPers {
	return &permissionPers{db: db}
}

// getResourceColumn returns the column name for the resource type
func getResourceColumn(resourceType domain.PermissionType) string {
	switch resourceType {
	case domain.PermissionTypeSpace:
		return "space_id"
	case domain.PermissionTypeDocument:
		return "document_id"
	case domain.PermissionTypeDatabase:
		return "database_id"
	case domain.PermissionTypeDrawing:
		return "drawing_id"
	default:
		return ""
	}
}

func (p *permissionPers) ListByResource(resourceType domain.PermissionType, resourceId string) ([]domain.Permission, error) {
	var perms []domain.Permission
	column := getResourceColumn(resourceType)
	if column == "" {
		return nil, nil
	}

	err := p.db.Preload("User").Preload("Group").
		Where("type = ? AND "+column+" = ? AND deleted_at IS NULL", resourceType, resourceId).
		Find(&perms).Error
	return perms, err
}

func (p *permissionPers) GetByResourceAndUser(resourceType domain.PermissionType, resourceId, userId string) (*domain.Permission, error) {
	var perm domain.Permission
	column := getResourceColumn(resourceType)
	if column == "" {
		return nil, nil
	}

	err := p.db.Where("type = ? AND "+column+" = ? AND user_id = ? AND deleted_at IS NULL", resourceType, resourceId, userId).
		First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (p *permissionPers) GetByResourceAndGroup(resourceType domain.PermissionType, resourceId, groupId string) (*domain.Permission, error) {
	var perm domain.Permission
	column := getResourceColumn(resourceType)
	if column == "" {
		return nil, nil
	}

	err := p.db.Where("type = ? AND "+column+" = ? AND group_id = ? AND deleted_at IS NULL", resourceType, resourceId, groupId).
		First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (p *permissionPers) UpsertUser(resourceType domain.PermissionType, resourceId, userId string, role domain.PermissionRole) error {
	column := getResourceColumn(resourceType)
	if column == "" {
		return nil
	}

	var perm domain.Permission
	err := p.db.Where("type = ? AND "+column+" = ? AND user_id = ?", resourceType, resourceId, userId).First(&perm).Error
	if err != nil {
		// Create if not exists
		perm = domain.Permission{
			Id:     utils.UUIDv4(),
			Type:   resourceType,
			UserId: &userId,
			Role:   role,
		}
		// Set the appropriate resource ID
		switch resourceType {
		case domain.PermissionTypeSpace:
			perm.SpaceId = &resourceId
		case domain.PermissionTypeDocument:
			perm.DocumentId = &resourceId
		case domain.PermissionTypeDatabase:
			perm.DatabaseId = &resourceId
		case domain.PermissionTypeDrawing:
			perm.DrawingId = &resourceId
		}
		return p.db.Create(&perm).Error
	}
	// Update existing
	perm.Role = role
	perm.DeletedAt = gorm.DeletedAt{} // Restore if soft-deleted
	return p.db.Save(&perm).Error
}

func (p *permissionPers) UpsertGroup(resourceType domain.PermissionType, resourceId, groupId string, role domain.PermissionRole) error {
	column := getResourceColumn(resourceType)
	if column == "" {
		return nil
	}

	var perm domain.Permission
	err := p.db.Where("type = ? AND "+column+" = ? AND group_id = ?", resourceType, resourceId, groupId).First(&perm).Error
	if err != nil {
		// Create if not exists
		perm = domain.Permission{
			Id:      utils.UUIDv4(),
			Type:    resourceType,
			GroupId: &groupId,
			Role:    role,
		}
		// Set the appropriate resource ID
		switch resourceType {
		case domain.PermissionTypeSpace:
			perm.SpaceId = &resourceId
		case domain.PermissionTypeDocument:
			perm.DocumentId = &resourceId
		case domain.PermissionTypeDatabase:
			perm.DatabaseId = &resourceId
		case domain.PermissionTypeDrawing:
			perm.DrawingId = &resourceId
		}
		return p.db.Create(&perm).Error
	}
	// Update existing
	perm.Role = role
	perm.DeletedAt = gorm.DeletedAt{} // Restore if soft-deleted
	return p.db.Save(&perm).Error
}

func (p *permissionPers) DeleteUser(resourceType domain.PermissionType, resourceId, userId string) error {
	column := getResourceColumn(resourceType)
	if column == "" {
		return nil
	}

	return p.db.Where("type = ? AND "+column+" = ? AND user_id = ?", resourceType, resourceId, userId).
		Delete(&domain.Permission{}).Error
}

func (p *permissionPers) DeleteGroup(resourceType domain.PermissionType, resourceId, groupId string) error {
	column := getResourceColumn(resourceType)
	if column == "" {
		return nil
	}

	return p.db.Where("type = ? AND "+column+" = ? AND group_id = ?", resourceType, resourceId, groupId).
		Delete(&domain.Permission{}).Error
}
