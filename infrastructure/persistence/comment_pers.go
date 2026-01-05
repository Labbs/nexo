package persistence

import (
	"fmt"

	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type commentPers struct {
	db *gorm.DB
}

func NewCommentPers(db *gorm.DB) *commentPers {
	return &commentPers{db: db}
}

func (p *commentPers) Create(comment *domain.Comment) error {
	return p.db.Create(comment).Error
}

func (p *commentPers) GetById(commentId string) (*domain.Comment, error) {
	var comment domain.Comment
	err := p.db.
		Preload("User").
		Preload("Parent").
		Where("id = ?", commentId).
		First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (p *commentPers) GetByDocumentId(documentId string) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := p.db.
		Preload("User").
		Where("document_id = ? AND deleted_at IS NULL", documentId).
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (p *commentPers) Update(comment *domain.Comment) error {
	return p.db.Save(comment).Error
}

func (p *commentPers) Delete(commentId string) error {
	if err := p.db.Where("id = ?", commentId).Delete(&domain.Comment{}).Error; err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

func (p *commentPers) Resolve(commentId string, resolved bool) error {
	return p.db.Model(&domain.Comment{}).Where("id = ?", commentId).Update("resolved", resolved).Error
}
