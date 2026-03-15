package group

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type GroupApplication struct {
	Config          config.Config
	Logger          zerolog.Logger
	GroupPers       domain.GroupPers
	UserApplication ports.UserPort
}

func NewGroupApplication(config config.Config, logger zerolog.Logger, groupPers domain.GroupPers) *GroupApplication {
	return &GroupApplication{
		Config:    config,
		Logger:    logger,
		GroupPers: groupPers,
	}
}
