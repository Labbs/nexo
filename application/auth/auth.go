package auth

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type AuthApp struct {
	Config      config.Config
	Logger      zerolog.Logger
	UserApp     ports.UserPort
	SessionApp  ports.SessionPort
	SpaceApp    ports.SpacePort
	DocumentApp ports.DocumentPort
}

func NewAuthApp(config config.Config, logger zerolog.Logger, userApp ports.UserPort, sessionApp ports.SessionPort, spaceApp ports.SpacePort, documentApp ports.DocumentPort) *AuthApp {
	return &AuthApp{
		Config:      config,
		Logger:      logger,
		UserApp:     userApp,
		SessionApp:  sessionApp,
		SpaceApp:    spaceApp,
		DocumentApp: documentApp,
	}
}
