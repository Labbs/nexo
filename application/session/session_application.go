package session

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type SessionApplication struct {
	Config      config.Config
	Logger      zerolog.Logger
	SessionPers domain.SessionPers
	UserApp     ports.UserPort
}

func NewSessionApplication(config config.Config, logger zerolog.Logger, sessionPers domain.SessionPers, userApp ports.UserPort) *SessionApplication {
	return &SessionApplication{
		Config:      config,
		Logger:      logger,
		SessionPers: sessionPers,
		UserApp:     userApp,
	}
}
