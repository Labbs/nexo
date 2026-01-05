package persistence

import (
	"time"

	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type actionPers struct {
	db *gorm.DB
}

func NewActionPers(db *gorm.DB) *actionPers {
	return &actionPers{db: db}
}

func (p *actionPers) Create(action *domain.Action) error {
	return p.db.Create(action).Error
}

func (p *actionPers) GetById(id string) (*domain.Action, error) {
	var action domain.Action
	err := p.db.
		Preload("User").
		Preload("Space").
		Where("id = ?", id).
		First(&action).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (p *actionPers) GetByUserId(userId string) ([]domain.Action, error) {
	var actions []domain.Action
	err := p.db.
		Preload("Space").
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Find(&actions).Error
	if err != nil {
		return nil, err
	}
	return actions, nil
}

func (p *actionPers) GetActiveByTrigger(triggerType domain.ActionTriggerType, spaceId *string, databaseId *string) ([]domain.Action, error) {
	var actions []domain.Action
	query := p.db.Where("active = ? AND trigger_type = ?", true, triggerType)

	if spaceId != nil {
		query = query.Where("space_id = ? OR space_id IS NULL", *spaceId)
	}

	if databaseId != nil {
		query = query.Where("database_id = ? OR database_id IS NULL", *databaseId)
	}

	err := query.Find(&actions).Error
	if err != nil {
		return nil, err
	}
	return actions, nil
}

func (p *actionPers) Update(action *domain.Action) error {
	return p.db.Save(action).Error
}

func (p *actionPers) Delete(id string) error {
	return p.db.Where("id = ?", id).Delete(&domain.Action{}).Error
}

func (p *actionPers) IncrementSuccess(id string) error {
	return p.db.Model(&domain.Action{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"success_count": gorm.Expr("success_count + 1"),
			"run_count":     gorm.Expr("run_count + 1"),
			"last_error":    "",
		}).Error
}

func (p *actionPers) RecordFailure(id string, errorMsg string) error {
	return p.db.Model(&domain.Action{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"failure_count": gorm.Expr("failure_count + 1"),
			"run_count":     gorm.Expr("run_count + 1"),
			"last_error":    errorMsg,
		}).Error
}

func (p *actionPers) UpdateLastRun(id string) error {
	now := time.Now()
	return p.db.Model(&domain.Action{}).
		Where("id = ?", id).
		Update("last_run_at", &now).Error
}

// ActionRun persistence
type actionRunPers struct {
	db *gorm.DB
}

func NewActionRunPers(db *gorm.DB) *actionRunPers {
	return &actionRunPers{db: db}
}

func (p *actionRunPers) Create(run *domain.ActionRun) error {
	return p.db.Create(run).Error
}

func (p *actionRunPers) GetByActionId(actionId string, limit int) ([]domain.ActionRun, error) {
	var runs []domain.ActionRun
	err := p.db.
		Where("action_id = ?", actionId).
		Order("created_at DESC").
		Limit(limit).
		Find(&runs).Error
	if err != nil {
		return nil, err
	}
	return runs, nil
}
