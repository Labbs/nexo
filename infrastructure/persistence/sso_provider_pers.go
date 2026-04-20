package persistence

import (
	"errors"

	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"gorm.io/gorm"
)

type oauthProviderPers struct {
	db *gorm.DB
}

func NewOAuthProviderPers(db *gorm.DB) *oauthProviderPers {
	return &oauthProviderPers{db: db}
}

func (o *oauthProviderPers) FindByProviderAndSubject(provider, subject string) (domain.OAuthProvider, error) {
	var op domain.OAuthProvider
	err := o.db.Where("provider = ? AND provider_user_id = ?", provider, subject).First(&op).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return op, apperrors.ErrNotFound
	}
	return op, err
}

func (o *oauthProviderPers) FindByUserId(userId string) ([]domain.OAuthProvider, error) {
	var ops []domain.OAuthProvider
	err := o.db.Where("user_id = ?", userId).Find(&ops).Error
	return ops, err
}

func (o *oauthProviderPers) Create(op domain.OAuthProvider) (domain.OAuthProvider, error) {
	err := o.db.Create(&op).Error
	return op, err
}
