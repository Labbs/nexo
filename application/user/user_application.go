package user

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type UserApp struct {
	Config   config.Config
	Logger   zerolog.Logger
	UserPres domain.UserPers
	GroupApp ports.GroupPort
}

func NewUserApp(
	config config.Config,
	logger zerolog.Logger,
	userPers domain.UserPers) *UserApp {
	return &UserApp{
		Config:   config,
		Logger:   logger,
		UserPres: userPers,
	}
}
