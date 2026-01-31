package dto

import "github.com/labbs/nexo/domain"

type UpdateProfileInput struct {
	UserId      string
	Username    *string
	AvatarUrl   *string
	Preferences *domain.JSONB
}

type UpdateProfileOutput struct {
	User *domain.User
}

type ChangePasswordInput struct {
	UserId          string
	CurrentPassword string
	NewPassword     string
}
