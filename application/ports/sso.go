package ports

import "github.com/labbs/nexo/application/auth/dto"

type SSOPort interface {
	GetRedirectURL() (*dto.SSORedirectOutput, error)
	HandleCallback(input dto.SSOCallbackInput) (*dto.SSOCallbackOutput, error)
}
