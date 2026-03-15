package action

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type ActionApplication struct {
	Config        config.Config
	Logger        zerolog.Logger
	ActionPers    domain.ActionPers
	ActionRunPers domain.ActionRunPers
}

func NewActionApplication(config config.Config, logger zerolog.Logger, actionPers domain.ActionPers, actionRunPers domain.ActionRunPers) *ActionApplication {
	return &ActionApplication{
		Config:        config,
		Logger:        logger,
		ActionPers:    actionPers,
		ActionRunPers: actionRunPers,
	}
}
