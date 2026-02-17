package database

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type DatabaseApplication struct {
	Config          config.Config
	Logger          zerolog.Logger
	DatabasePers    domain.DatabasePers
	DatabaseRowPers domain.DatabaseRowPers
	SpacePers       domain.SpacePers
	PermissionPers  domain.PermissionPers
}

func NewDatabaseApplication(config config.Config, logger zerolog.Logger, databasePers domain.DatabasePers, databaseRowPers domain.DatabaseRowPers, spacePers domain.SpacePers, permissionPers domain.PermissionPers) *DatabaseApplication {
	return &DatabaseApplication{
		Config:          config,
		Logger:          logger,
		DatabasePers:    databasePers,
		DatabaseRowPers: databaseRowPers,
		SpacePers:       spacePers,
		PermissionPers:  permissionPers,
	}
}
