package auth

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type AuthApplication struct {
	Config      config.Config
	Logger      zerolog.Logger
	UserApp     ports.UserPort
	SessionApp  ports.SessionPort
	SpaceApp    ports.SpacePort
	DocumentApp ports.DocumentPort
}

func NewAuthApplication(config config.Config, logger zerolog.Logger, userApp ports.UserPort, sessionApp ports.SessionPort, spaceApp ports.SpacePort, documentApp ports.DocumentPort) *AuthApplication {
	return &AuthApplication{
		Config:      config,
		Logger:      logger,
		UserApp:     userApp,
		SessionApp:  sessionApp,
		SpaceApp:    spaceApp,
		DocumentApp: documentApp,
	}
}
