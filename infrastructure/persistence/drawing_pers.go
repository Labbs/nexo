package persistence

import (
	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type drawingPers struct {
	db *gorm.DB
}

func NewDrawingPers(db *gorm.DB) *drawingPers {
	return &drawingPers{db: db}
}

func (p *drawingPers) Create(drawing *domain.Drawing) error {
	return p.db.Create(drawing).Error
}

func (p *drawingPers) GetById(id string) (*domain.Drawing, error) {
	var drawing domain.Drawing
	err := p.db.Debug().
		Preload("Space").
		Preload("User").
		Where("id = ?", id).
		First(&drawing).Error
	if err != nil {
		return nil, err
	}
	return &drawing, nil
}

func (p *drawingPers) GetBySpaceId(spaceId string) ([]domain.Drawing, error) {
	var drawings []domain.Drawing
	err := p.db.Debug().
		Preload("User").
		Where("space_id = ?", spaceId).
		Order("created_at DESC").
		Find(&drawings).Error
	if err != nil {
		return nil, err
	}
	return drawings, nil
}

func (p *drawingPers) GetByDocumentId(documentId string) ([]domain.Drawing, error) {
	var drawings []domain.Drawing
	err := p.db.Debug().
		Preload("User").
		Where("document_id = ?", documentId).
		Order("created_at DESC").
		Find(&drawings).Error
	if err != nil {
		return nil, err
	}
	return drawings, nil
}

func (p *drawingPers) Update(drawing *domain.Drawing) error {
	return p.db.Debug().Save(drawing).Error
}

func (p *drawingPers) Delete(id string) error {
	return p.db.Debug().Where("id = ?", id).Delete(&domain.Drawing{}).Error
}
