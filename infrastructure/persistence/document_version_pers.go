package persistence

import (
	"fmt"

	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type documentVersionPers struct {
	db *gorm.DB
}

func NewDocumentVersionPers(db *gorm.DB) *documentVersionPers {
	return &documentVersionPers{db: db}
}

func (p *documentVersionPers) Create(version *domain.DocumentVersion) error {
	return p.db.Debug().Create(version).Error
}

func (p *documentVersionPers) GetByDocumentId(documentId string, limit int, offset int) ([]domain.DocumentVersion, error) {
	var versions []domain.DocumentVersion

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	err := p.db.Debug().
		Preload("User").
		Where("document_id = ?", documentId).
		Order("version DESC").
		Limit(limit).
		Offset(offset).
		Find(&versions).Error

	if err != nil {
		return nil, err
	}

	return versions, nil
}

func (p *documentVersionPers) GetById(versionId string) (*domain.DocumentVersion, error) {
	var version domain.DocumentVersion
	err := p.db.Debug().
		Preload("User").
		Preload("Document").
		Where("id = ?", versionId).
		First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (p *documentVersionPers) GetLatestVersion(documentId string) (*domain.DocumentVersion, error) {
	var version domain.DocumentVersion
	err := p.db.Debug().
		Where("document_id = ?", documentId).
		Order("version DESC").
		First(&version).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &version, nil
}

func (p *documentVersionPers) GetVersionCount(documentId string) (int64, error) {
	var count int64
	err := p.db.Debug().Model(&domain.DocumentVersion{}).
		Where("document_id = ?", documentId).
		Count(&count).Error
	return count, err
}

func (p *documentVersionPers) DeleteOldVersions(documentId string, keepCount int) error {
	if keepCount <= 0 {
		return fmt.Errorf("keepCount must be positive")
	}

	// Get IDs of versions to keep
	var versionsToKeep []string
	err := p.db.Debug().Model(&domain.DocumentVersion{}).
		Select("id").
		Where("document_id = ?", documentId).
		Order("version DESC").
		Limit(keepCount).
		Pluck("id", &versionsToKeep).Error
	if err != nil {
		return err
	}

	if len(versionsToKeep) == 0 {
		return nil
	}

	// Delete versions not in the keep list
	return p.db.Debug().
		Where("document_id = ? AND id NOT IN ?", documentId, versionsToKeep).
		Delete(&domain.DocumentVersion{}).Error
}
