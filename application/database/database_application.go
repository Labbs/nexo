package database

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type DatabaseApplication struct {
	Config          config.Config
	Logger          zerolog.Logger
	DatabasePers    domain.DatabasePers
	DatabaseRowPers domain.DatabaseRowPers
	SpaceApp        ports.SpacePort
	PermissionApp   ports.PermissionPort
}

func NewDatabaseApplication(config config.Config, logger zerolog.Logger, databasePers domain.DatabasePers, databaseRowPers domain.DatabaseRowPers) *DatabaseApplication {
	return &DatabaseApplication{
		Config:          config,
		Logger:          logger,
		DatabasePers:    databasePers,
		DatabaseRowPers: databaseRowPers,
	}
}
