package permission

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type PermissionApplication struct {
	Config         config.Config
	Logger         zerolog.Logger
	PermissionPers domain.PermissionPers
	SpaceApp       ports.SpacePort
	DrawingApp     ports.DrawingPort
	DocumentApp    ports.DocumentPort
	DatabaseApp    ports.DatabasePort
}

func NewPermissionApplication(
	config config.Config,
	logger zerolog.Logger,
	permissionPers domain.PermissionPers,
) *PermissionApplication {
	return &PermissionApplication{
		Config:         config,
		Logger:         logger,
		PermissionPers: permissionPers,
	}
}
