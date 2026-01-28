package persistence

import (
	"fmt"

	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type documentPers struct {
	db *gorm.DB
}

func NewDocumentPers(db *gorm.DB) *documentPers {
	return &documentPers{db: db}
}

func (p *documentPers) GetDocumentWithPermissions(documentId, userId string) (*domain.Document, error) {
	var doc domain.Document
	err := p.db.
		Preload("Space", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Owner").
				Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId)
		}).
		Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId).
		Where("id = ?", documentId).
		First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (p *documentPers) GetDocumentByIdOrSlugWithUserPermissions(spaceId string, id *string, slug *string, userId string) (*domain.Document, error) {
	var doc domain.Document

	query := p.db.Debug().
		// Preload the space along with owner and its permissions for the user
		Preload("Space", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Owner").
				Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId)
		}).
		// Preload only the document's permissions for this user
		Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId).
		// Optionally preload the parent if needed
		Preload("Parent").
		Where("space_id = ?", spaceId)

	if id != nil {
		query = query.Where("id = ?", *id)
	} else if slug != nil {
		query = query.Where("slug = ?", *slug)
	} else {
		return nil, fmt.Errorf("either id or slug must be provided")
	}

	err := query.First(&doc).Error
	if err != nil {
		return nil, err
	}

	// Verify if the user has at least viewer permissions
	if !doc.HasPermission(userId, domain.PermissionRoleViewer) {
		return nil, fmt.Errorf("access denied: insufficient permissions")
	}

	return &doc, nil
}

func (p *documentPers) GetRootDocumentsFromSpaceWithUserPermissions(spaceId, userId string) ([]domain.Document, error) {
	var docs []domain.Document

	err := p.db.Debug().
		// Preload space with owner and permissions for the user
		Preload("Space", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Owner").
				Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId)
		}).
		// Preload document permissions for this user
		Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId).
		Where("space_id = ? AND parent_id IS NULL AND deleted_at IS NULL", spaceId).
		Order("position ASC, created_at ASC").
		Find(&docs).Error

	if err != nil {
		return nil, err
	}

	// Filter documents based on permissions
	var accessibleDocs []domain.Document
	for _, doc := range docs {
		if doc.HasPermission(userId, domain.PermissionRoleViewer) {
			accessibleDocs = append(accessibleDocs, doc)
		}
	}

	return accessibleDocs, nil
}

func (p *documentPers) GetChildDocumentsWithUserPermissions(parentId, userId string) ([]domain.Document, error) {
	var docs []domain.Document

	err := p.db.Debug().
		// Preload space with owner and permissions for the user
		Preload("Space", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Owner").
				Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId)
		}).
		// Preload document permissions for this user
		Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId).
		// Preload the parent to have the complete context
		Preload("Parent").
		Where("parent_id = ? AND deleted_at IS NULL", parentId).
		Order("position ASC, created_at ASC").
		Find(&docs).Error

	if err != nil {
		return nil, err
	}

	// Filter documents based on permissions
	var accessibleDocs []domain.Document
	for _, doc := range docs {
		if doc.HasPermission(userId, domain.PermissionRoleViewer) {
			accessibleDocs = append(accessibleDocs, doc)
		}
	}

	return accessibleDocs, nil
}

func (p *documentPers) Create(document *domain.Document, userId string) error {
	// If the document has a parent, check permissions on the parent
	if document.ParentId != nil {
		parentDoc, err := p.GetDocumentWithPermissions(*document.ParentId, userId)
		if err != nil {
			return fmt.Errorf("failed to get parent document: %w", err)
		}

		// Check if the user can edit the parent document
		if !parentDoc.HasPermission(userId, domain.PermissionRoleEditor) {
			return fmt.Errorf("access denied: insufficient permissions to create document in parent")
		}
	} else {
		// If it's a root document, check permissions on the space
		var space domain.Space
		err := p.db.Debug().
			Preload("Owner").
			Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId).
			Where("id = ?", document.SpaceId).
			First(&space).Error
		if err != nil {
			return fmt.Errorf("failed to get space: %w", err)
		}

		// Check if the user can edit in the space
		if !space.HasPermission(userId, domain.PermissionRoleEditor) {
			return fmt.Errorf("access denied: insufficient permissions to create document in space")
		}
	}

	// Assign position at the end
	maxPos, err := p.GetMaxPosition(document.SpaceId, document.ParentId)
	if err != nil {
		return fmt.Errorf("failed to get max position: %w", err)
	}
	document.Position = maxPos + 1

	// Perform the creation
	return p.db.Debug().Create(document).Error
}

func (p *documentPers) Update(document *domain.Document, userId string) error {
	// Get the existing document with permissions
	existingDoc, err := p.GetDocumentWithPermissions(document.Id, userId)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Check if the user can edit the document
	if !existingDoc.HasPermission(userId, domain.PermissionRoleEditor) {
		return fmt.Errorf("access denied: insufficient permissions to update document")
	}

	// Perform the update
	return p.db.Debug().Save(document).Error
}

func (p *documentPers) Delete(documentId, userId string) error {
	// Load document with permissions for the current user
	existingDoc, err := p.GetDocumentWithPermissions(documentId, userId)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Check if the user can edit/delete the document
	if !existingDoc.HasPermission(userId, domain.PermissionRoleEditor) {
		return fmt.Errorf("access denied: insufficient permissions to delete document")
	}

	// Prevent delete if there are active children
	var childrenCount int64
	if err := p.db.Model(&domain.Document{}).Where("parent_id = ? AND deleted_at IS NULL", documentId).Count(&childrenCount).Error; err != nil {
		return fmt.Errorf("failed to check child documents: %w", err)
	}
	if childrenCount > 0 {
		return fmt.Errorf("cannot delete document with existing child documents")
	}

	// Soft delete
	if err := p.db.Where("id = ?", documentId).Delete(&domain.Document{}).Error; err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

func (p *documentPers) Move(documentId string, newParentId *string, userId string) (*domain.Document, error) {
	// Load the document to move with permissions
	doc, err := p.GetDocumentWithPermissions(documentId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if !doc.HasPermission(userId, domain.PermissionRoleEditor) {
		return nil, fmt.Errorf("access denied: insufficient permissions to move document")
	}

	// Prevent self-parenting
	if newParentId != nil && *newParentId == documentId {
		return nil, fmt.Errorf("invalid move: cannot set document as its own parent")
	}

	// If moving under a parent, check permissions on parent and same space
	if newParentId != nil {
		parent, err := p.GetDocumentWithPermissions(*newParentId, userId)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent document: %w", err)
		}
		if !parent.HasPermission(userId, domain.PermissionRoleEditor) {
			return nil, fmt.Errorf("access denied: insufficient permissions on target parent")
		}
		if parent.SpaceId != doc.SpaceId {
			return nil, fmt.Errorf("invalid move: parent must be in the same space")
		}
		doc.ParentId = newParentId
	} else {
		// Move to root
		doc.ParentId = nil
	}

	// Assign position at the end of the new parent's children
	maxPos, err := p.GetMaxPosition(doc.SpaceId, doc.ParentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get max position: %w", err)
	}
	doc.Position = maxPos + 1

	if err := p.db.Save(doc).Error; err != nil {
		return nil, fmt.Errorf("failed to move document: %w", err)
	}

	return doc, nil
}

func (p *documentPers) GetDeletedDocuments(spaceId, userId string) ([]domain.Document, error) {
	var docs []domain.Document

	err := p.db.Unscoped().
		Preload("Space", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Owner").
				Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId)
		}).
		Where("space_id = ? AND deleted_at IS NOT NULL", spaceId).
		Order("deleted_at DESC").
		Find(&docs).Error

	if err != nil {
		return nil, err
	}

	// Filter based on space permissions (user must be at least editor to see trash)
	var accessibleDocs []domain.Document
	for _, doc := range docs {
		if doc.Space.HasPermission(userId, domain.PermissionRoleEditor) {
			accessibleDocs = append(accessibleDocs, doc)
		}
	}

	return accessibleDocs, nil
}

func (p *documentPers) Restore(documentId, userId string) error {
	// Get the deleted document (unscoped to include soft-deleted)
	var doc domain.Document
	err := p.db.Unscoped().
		Preload("Space", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Owner").
				Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId)
		}).
		Where("id = ?", documentId).
		First(&doc).Error
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Check if user has editor permission on the space
	if !doc.Space.HasPermission(userId, domain.PermissionRoleEditor) {
		return fmt.Errorf("access denied: insufficient permissions to restore document")
	}

	// Check if document is actually deleted
	if !doc.DeletedAt.Valid {
		return fmt.Errorf("document is not deleted")
	}

	// If document had a parent, check if parent still exists
	if doc.ParentId != nil {
		var parent domain.Document
		err := p.db.Where("id = ? AND deleted_at IS NULL", *doc.ParentId).First(&parent).Error
		if err != nil {
			// Parent doesn't exist or is deleted, restore to root
			doc.ParentId = nil
		}
	}

	// Restore the document by clearing deleted_at
	return p.db.Unscoped().Model(&doc).Update("deleted_at", nil).Error
}

func (p *documentPers) SetPublic(documentId string, public bool, userId string) error {
	// Get the document with permissions
	doc, err := p.GetDocumentWithPermissions(documentId, userId)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Check if user has editor permission
	if !doc.HasPermission(userId, domain.PermissionRoleEditor) {
		return fmt.Errorf("access denied: insufficient permissions")
	}

	return p.db.Model(&domain.Document{}).Where("id = ?", documentId).Update("public", public).Error
}

func (p *documentPers) GetPublicDocument(spaceId string, id *string, slug *string) (*domain.Document, error) {
	var doc domain.Document

	query := p.db.
		Preload("Space").
		Preload("Parent").
		Where("space_id = ? AND public = ?", spaceId, true)

	if id != nil {
		query = query.Where("id = ?", *id)
	} else if slug != nil {
		query = query.Where("slug = ?", *slug)
	} else {
		return nil, fmt.Errorf("either id or slug must be provided")
	}

	err := query.First(&doc).Error
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (p *documentPers) Search(query string, userId string, spaceId *string, limit int) ([]domain.Document, error) {
	var docs []domain.Document

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	searchPattern := "%" + query + "%"

	dbQuery := p.db.
		Preload("Space", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Owner").
				Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId)
		}).
		Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId).
		Where("deleted_at IS NULL").
		Where("(name LIKE ? OR content LIKE ?)", searchPattern, searchPattern)

	if spaceId != nil {
		dbQuery = dbQuery.Where("space_id = ?", *spaceId)
	}

	err := dbQuery.
		Order("updated_at DESC").
		Limit(limit).
		Find(&docs).Error

	if err != nil {
		return nil, err
	}

	// Filter documents based on permissions
	var accessibleDocs []domain.Document
	for _, doc := range docs {
		if doc.HasPermission(userId, domain.PermissionRoleViewer) {
			accessibleDocs = append(accessibleDocs, doc)
		}
	}

	return accessibleDocs, nil
}

func (p *documentPers) Reorder(spaceId string, items []domain.ReorderItem, userId string) error {
	// Verify user has editor access on the space
	var space domain.Space
	err := p.db.Debug().
		Preload("Owner").
		Preload("Permissions", "user_id = ? AND deleted_at IS NULL", userId).
		Where("id = ?", spaceId).
		First(&space).Error
	if err != nil {
		return fmt.Errorf("failed to get space: %w", err)
	}
	if !space.HasPermission(userId, domain.PermissionRoleEditor) {
		return fmt.Errorf("access denied: insufficient permissions to reorder documents")
	}

	// Update positions in a transaction
	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&domain.Document{}).
				Where("id = ? AND space_id = ?", item.Id, spaceId).
				Update("position", item.Position).Error; err != nil {
				return fmt.Errorf("failed to update position for document %s: %w", item.Id, err)
			}
		}
		return nil
	})
}

func (p *documentPers) GetMaxPosition(spaceId string, parentId *string) (int, error) {
	var maxPos int
	query := p.db.Model(&domain.Document{}).
		Where("space_id = ? AND deleted_at IS NULL", spaceId)

	if parentId != nil {
		query = query.Where("parent_id = ?", *parentId)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	err := query.Select("COALESCE(MAX(position), -1)").Scan(&maxPos).Error
	if err != nil {
		return -1, err
	}
	return maxPos, nil
}
