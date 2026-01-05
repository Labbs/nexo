package persistence

import (
	"github.com/google/uuid"
	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type groupPers struct {
	db *gorm.DB
}

func NewGroupPers(db *gorm.DB) *groupPers {
	return &groupPers{db: db}
}

func (g *groupPers) Create(group *domain.Group) error {
	if group.Id == "" {
		group.Id = uuid.New().String()
	}
	return g.db.Create(group).Error
}

func (g *groupPers) GetById(groupId string) (*domain.Group, error) {
	var group domain.Group
	err := g.db.Preload("Members").Preload("OwnerUser").Where("id = ?", groupId).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (g *groupPers) GetAll(limit, offset int) ([]domain.Group, int64, error) {
	var groups []domain.Group
	var total int64

	// Count total
	if err := g.db.Model(&domain.Group{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with members and owner
	err := g.db.Preload("Members").Preload("OwnerUser").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&groups).Error

	return groups, total, err
}

func (g *groupPers) Update(group *domain.Group) error {
	return g.db.Model(group).Updates(map[string]interface{}{
		"name":        group.Name,
		"description": group.Description,
		"role":        group.Role,
	}).Error
}

func (g *groupPers) Delete(groupId string) error {
	// First remove all members
	if err := g.db.Exec("DELETE FROM group_members WHERE group_id = ?", groupId).Error; err != nil {
		return err
	}
	// Then delete the group
	return g.db.Delete(&domain.Group{}, "id = ?", groupId).Error
}

func (g *groupPers) AddMember(groupId, userId string) error {
	return g.db.Exec(
		"INSERT INTO group_members (group_id, user_id, created_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT DO NOTHING",
		groupId, userId,
	).Error
}

func (g *groupPers) RemoveMember(groupId, userId string) error {
	return g.db.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ?", groupId, userId).Error
}

func (g *groupPers) GetMembers(groupId string) ([]domain.User, error) {
	var users []domain.User
	err := g.db.
		Joins("JOIN group_members ON group_members.user_id = \"user\".id").
		Where("group_members.group_id = ?", groupId).
		Find(&users).Error
	return users, err
}
