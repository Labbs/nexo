package space

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type SpaceApplication struct {
	Config         config.Config
	Logger         zerolog.Logger
	SpacePres      domain.SpacePers
	PermissionPers domain.PermissionPers // kept for admin operations only
	DocumentApp    ports.DocumentPort
	PermissionApp  ports.PermissionPort
}

func NewSpaceApplication(config config.Config, logger zerolog.Logger, spacePers domain.SpacePers, permissionPers domain.PermissionPers) *SpaceApplication {
	return &SpaceApplication{
		Config:         config,
		Logger:         logger,
		SpacePres:      spacePers,
		PermissionPers: permissionPers,
	}
}
