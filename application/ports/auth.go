package ports

import (
	"github.com/labbs/nexo/application/auth/dto"
)

type AuthPort interface {
	Authenticate(input dto.AuthenticateInput) (*dto.AuthenticateOutput, error)
	Register(input dto.RegisterInput) error
	Logout(input dto.LogoutInput) error
}
