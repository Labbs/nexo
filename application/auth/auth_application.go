package auth

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type AuthApplication struct {
	Config              config.Config
	Logger              zerolog.Logger
	UserApplication     ports.UserPort
	SessionApplication  ports.SessionPort
	SpaceApplication    ports.SpacePort
	DocumentApplication ports.DocumentPort
}

func NewAuthApplication(config config.Config, logger zerolog.Logger) *AuthApplication {
	return &AuthApplication{
		Config: config,
		Logger: logger,
	}
}
