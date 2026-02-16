package ports

import (
	"github.com/labbs/nexo/application/user/dto"
)

type UserPort interface {
	Create(input dto.CreateUserInput) (*dto.CreateUserOutput, error)
	GetByEmail(input dto.GetByEmailInput) (*dto.GetByEmailOutput, error)
	GetByUserId(input dto.GetByUserIdInput) (*dto.GetByUserIdOutput, error)
	UpdateProfile(input dto.UpdateProfileInput) (*dto.UpdateProfileOutput, error)
	ChangePassword(input dto.ChangePasswordInput) error
	UpdateSpaceOrder(input dto.UpdateSpaceOrderInput) (*dto.UpdateSpaceOrderOutput, error)
}
