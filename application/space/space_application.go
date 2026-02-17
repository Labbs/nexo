package space

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type SpaceApplication struct {
	Config         config.Config
	Logger         zerolog.Logger
	SpacePres      domain.SpacePers
	DocumentPers   domain.DocumentPers
	PermissionPers domain.PermissionPers
}

func NewSpaceApplication(config config.Config, logger zerolog.Logger, spacePers domain.SpacePers, documentPers domain.DocumentPers, permissionPers domain.PermissionPers) *SpaceApplication {
	return &SpaceApplication{
		Config:         config,
		Logger:         logger,
		SpacePres:      spacePers,
		DocumentPers:   documentPers,
		PermissionPers: permissionPers,
	}
}
