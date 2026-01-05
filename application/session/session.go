package session

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type SessionApp struct {
	Config      config.Config
	Logger      zerolog.Logger
	SessionPers domain.SessionPers
	UserApp     ports.UserPort
}

func NewSessionApp(config config.Config, logger zerolog.Logger, sessionPers domain.SessionPers, userApp ports.UserPort) *SessionApp {
	return &SessionApp{
		Config:      config,
		Logger:      logger,
		SessionPers: sessionPers,
		UserApp:     userApp,
	}
}
